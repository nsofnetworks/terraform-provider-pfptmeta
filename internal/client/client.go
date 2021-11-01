package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"time"
)

const (
	baseUrlEnvVar      string = "PFPTMETA_BASE_URL"
	baseURL            string = "https://api.metanetworks.com"
	oauthURL           string = "/v1/oauth/token"
	maxIdleConnections int    = 10
	requestTimeout     int    = 15
	configPath         string = ".pfptmeta/credentials.json"
	grantType          string = "client_credentials"
)

type Config struct {
	APIKey       string `json:"api_key"`
	APISecret    string `json:"api_secret"`
	OrgShortname string `json:"org_shortname"`
}

type Token struct {
	Token     string `json:"access_token"`
	Expiry    int64  `json:"expires_in"`
	TokenType string `json:"token_type"`
}

type Credentials struct {
	GrantType    string `json:"grant_type"`
	Scope        string `json:"scope"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type ErrorResponse struct {
	URL    string
	Method string
	Detail string `json:"detail"`
	Status int    `json:"status"`
	Title  string `json:"title"`
	Type   string `json:"type"`
}

func (err *ErrorResponse) Error() string {
	return fmt.Sprintf("%s request to %s failed with status code %d: %s - %s",
		err.Method, err.URL, err.Status, err.Title, err.Detail)
}

func newCredentials(d *schema.ResourceData) (*Credentials, error) {
	credentials := &Credentials{GrantType: grantType}
	apiKey, haveAPIKey := d.GetOk("api_key")
	apiSecret, haveAPISecret := d.GetOk("api_secret")
	org, haveOrg := d.GetOk("org_shortname")
	// If one is set
	if haveAPIKey || haveAPISecret || haveOrg {
		// but not all are set
		if !(haveAPIKey && haveAPISecret && haveOrg) {
			return nil, errors.New("please provide an api_key, api_secret and org shortname")
		}
		credentials.ClientID = apiKey.(string)
		credentials.ClientSecret = apiSecret.(string)
		credentials.Scope = fmt.Sprintf("org:%s", org)
		return credentials, nil
	}
	return parseCredentialsFile(credentials)
}

func parseCredentialsFile(credentials *Credentials) (*Credentials, error) {
	usr, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("could not find current user's name: %v", err)
	}
	path := filepath.Join(usr.HomeDir, configPath)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open credentials file %s: %v", path, err)
	}
	configBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("could not read credentials file %s: %v", path, err)
	}
	var config Config
	err = json.Unmarshal(configBytes, &config)
	if err != nil {
		return nil, fmt.Errorf("could not parse credentials file %s: %v", path, err)
	}
	credentials.ClientID = config.APIKey
	credentials.ClientSecret = config.APISecret
	credentials.Scope = fmt.Sprintf("org:%s", config.OrgShortname)
	return credentials, nil
}

type Client struct {
	Credentials       *Credentials
	Token             *Token
	HTTPClient        *http.Client
	BaseURL           string
	TokenCreationTime int64
	UserAgent         string
}

func parseHttpError(resp *http.Response) error {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not parse authentication response: %v", err)
	}
	errorResponse := &ErrorResponse{}
	err = json.Unmarshal(body, errorResponse)
	if err != nil {
		return fmt.Errorf("could not parse error response: %v", err)
	}
	errorResponse.URL = resp.Request.URL.String()
	errorResponse.Method = resp.Request.Method
	return errorResponse
}

func NewClient(d *schema.ResourceData, userAgent string) (*Client, error) {
	client := &Client{
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: maxIdleConnections,
			},
			Timeout: time.Duration(requestTimeout) * time.Second,
		},
		UserAgent: userAgent,
	}
	if url := os.Getenv(baseUrlEnvVar); url != "" {
		client.BaseURL = url
	} else {
		client.BaseURL = baseURL
	}

	credentials, err := newCredentials(d)
	if err != nil {
		return nil, fmt.Errorf("could not find credentials: %v", err)
	}
	client.Credentials = credentials
	err = client.tokenRequest()
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) tokenRequest() error {
	jsonData, err := json.Marshal(c.Credentials)
	if err != nil {
		return fmt.Errorf("could not convert credentials to json: %v", err)
	}
	url := fmt.Sprintf("%s%s", c.BaseURL, oauthURL)
	resp, err := c.HTTPClient.Post(url, "application/json", bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error while trying to create access token: %v", err)
	}
	if err != nil {
		return fmt.Errorf("could not read authentication response: %v", err)
	}
	if resp.StatusCode != 200 {
		return parseHttpError(resp)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	token := &Token{}
	err = json.Unmarshal(body, token)
	if err != nil {
		return fmt.Errorf("could not parse authentication response: %v", err)
	}
	c.Token = token
	c.TokenCreationTime = time.Now().Unix()
	return nil
}

func (c *Client) SendRequest(r *http.Request) (*http.Response, error) {
	now := time.Now().Unix()
	if c.Token == nil || c.TokenCreationTime+c.Token.Expiry-now < 30 {
		err := c.tokenRequest()
		if err != nil {
			return nil, err
		}
	}
	r.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token.Token))
	r.Header.Add("User-Agent", c.UserAgent)
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		r.Header.Set("Content-Type", "application/json")
	case http.MethodPatch:
		r.Header.Set("Content-Type", "application/merge-patch+json")
	}
	resp, err := c.HTTPClient.Do(r)
	if err != nil {
		return nil, fmt.Errorf("failed to execute %s request to %s: %v", r.Method, r.URL, err)
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		return nil, parseHttpError(resp)
	}
	return resp, nil
}

func (c *Client) Get(url string, queryParams url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = queryParams.Encode()
	resp, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Delete(url string, queryParams url.Values) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	if err != nil {
		return nil, err
	}
	req.URL.RawQuery = queryParams.Encode()
	resp, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Post(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Patch(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Put(url string, body io.Reader) (*http.Response, error) {
	req, err := http.NewRequest(http.MethodPut, url, body)
	if err != nil {
		return nil, err
	}
	resp, err := c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
