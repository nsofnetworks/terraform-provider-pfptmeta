package client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
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
	baseUrlEnvVar             string = "PFPTMETA_BASE_URL"
	baseURL                   string = "https://api.metanetworks.com"
	oauthURL                  string = "/v1/oauth/token"
	maxIdleConnections        int    = 10
	requestTimeout            int    = 15
	configPath                string = ".pfptmeta/credentials.json"
	grantType                 string = "client_credentials"
	eventuallyConsistentSleep        = time.Millisecond * 300
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
	HTTP              *retryablehttp.Client
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

// RetryPolicy is a callback for Client.CheckRetry, which
// will status codes 409, 502, 504.
func RetryPolicy(ctx context.Context, resp *http.Response, err error) (bool, error) {
	// do not retry on context.Canceled or context.DeadlineExceeded
	if ctx.Err() != nil {
		return false, ctx.Err()
	}
	if resp == nil {
		return false, err
	}
	switch resp.StatusCode {
	case http.StatusConflict, http.StatusBadGateway, http.StatusGatewayTimeout, http.StatusTooManyRequests:
		return true, nil
	}
	return false, nil
}

func errorHandler(resp *http.Response, err error, _ int) (*http.Response, error) {
	return resp, err
}

func NewClient(ctx context.Context, d *schema.ResourceData, userAgent string) (*Client, error) {
	client := &Client{
		HTTP:      retryablehttp.NewClient(),
		UserAgent: userAgent,
	}
	client.HTTP.HTTPClient = &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: maxIdleConnections},
		Timeout:   time.Duration(requestTimeout) * time.Second,
	}
	client.HTTP.CheckRetry = RetryPolicy
	client.HTTP.ErrorHandler = errorHandler
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
	err = client.tokenRequest(ctx)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (c *Client) tokenRequest(ctx context.Context) error {
	jsonData, err := json.Marshal(c.Credentials)
	if err != nil {
		return fmt.Errorf("could not convert credentials to json: %v", err)
	}
	u := fmt.Sprintf("%s%s", c.BaseURL, oauthURL)
	r, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(jsonData))
	if err != nil {
		return fmt.Errorf("error while trying to create access token request: %v", err)
	}
	r.Header.Set("Content-Type", "application/json")
	retryableRequest, err := retryablehttp.FromRequest(r)
	if err != nil {
		return fmt.Errorf("error while trying to create access token retryable request: %v", err)
	}
	resp, err := c.HTTP.Do(retryableRequest)
	if err != nil {
		return fmt.Errorf("error while trying to create access token: %v", err)
	}
	if resp.StatusCode != 200 {
		return parseHttpError(resp)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("could not read authentication response: %v", err)
	}
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
		err := c.tokenRequest(r.Context())
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
	retryableRequest, err := retryablehttp.FromRequest(r)
	if err != nil {
		return nil, err
	}
	resp, err := c.HTTP.Do(retryableRequest)
	if err != nil {
		if resp == nil {
			return nil, fmt.Errorf("failed to execute %s request to %s: %v", r.Method, r.URL, err)
		} else {
			return nil, parseHttpError(resp)
		}
	}
	if resp.StatusCode < http.StatusOK || resp.StatusCode > http.StatusIMUsed {
		return nil, parseHttpError(resp)
	}
	// Sometimes after writing it can take up to 100 milliseconds for the resource to actually be consistent in all documentdb server instances.
	// To make sure the next read will be consistent we will sleep after writing finishes.
	switch r.Method {
	case http.MethodPost, http.MethodPut, http.MethodPatch, http.MethodDelete:
		time.Sleep(eventuallyConsistentSleep)
	}
	return resp, nil
}

func (c *Client) Get(ctx context.Context, url string, queryParams url.Values) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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

func (c *Client) Delete(ctx context.Context, url string, queryParams url.Values) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
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

func (c *Client) Post(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	var resp *http.Response
	var err error
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return nil, err
	}
	resp, err = c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Patch(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
	var resp *http.Response
	var err error
	req, err := http.NewRequest(http.MethodPatch, url, body)
	if err != nil {
		return nil, err
	}
	resp, err = c.SendRequest(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) Put(ctx context.Context, url string, body io.Reader) (*http.Response, error) {
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

func (c *Client) GetResource(ctx context.Context, resourceUrl, ID string) (*http.Response, error) {
	u := fmt.Sprintf("%s/%s/%s", c.BaseURL, resourceUrl, ID)
	resp, err := c.Get(ctx, u, nil)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
