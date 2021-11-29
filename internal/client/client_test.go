package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"
)

var retryCounter int = 0

func TestParseHttpError(t *testing.T) {
	errorResponse := &ErrorResponse{
		"",
		"",
		"error details",
		403,
		"title",
		"type",
	}
	endpoint, _ := url.Parse("https://example.com")
	bytesRes, _ := json.Marshal(errorResponse)
	res := &http.Response{
		Body: io.NopCloser(bytes.NewReader(bytesRes)),
		Request: &http.Request{
			Method: http.MethodPost,
			URL:    endpoint,
		},
	}
	formattedError := parseHttpError(res)
	assert.EqualError(t, formattedError, "POST request to https://example.com failed with status code 403: title - error details")
}

func TestDoRequest(t *testing.T) {
	server := configureServer(t)
	client := &Client{
		HTTP:        retryablehttp.NewClient(),
		BaseURL:     server.URL,
		Credentials: &Credentials{},
	}
	client.HTTP.HTTPClient = &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: maxIdleConnections},
		Timeout:   time.Duration(requestTimeout) * time.Second,
	}
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/v1/test", nil)

	t.Run("check-token-loaded-on-first-request", func(t *testing.T) {
		// Asserting token struct is empty before request
		assert.Nil(t, client.Token)
		_, err := client.SendRequest(req)
		assert.Nil(t, err)
		// Asserting token struct is not empty after request
		assert.NotNil(t, client.Token)
		//Asserting token was created for the first time
		assert.Equal(t, "token-1", client.Token.Token)
	})
	t.Run("check-token-reused-before-expiration", func(t *testing.T) {
		_, err := client.SendRequest(req)
		assert.Nil(t, err)
		//Asserting token was created for the second time
		assert.Equal(t, "token-1", client.Token.Token)
	})
	t.Run("check-token-refreshed-on-expiration", func(t *testing.T) {
		client.TokenCreationTime = time.Now().Unix() - 60*5
		_, err := client.SendRequest(req)
		assert.Nil(t, err)
		//Asserting token was created for the second time
		assert.Equal(t, "token-2", client.Token.Token)
	})
}

func TestRetry(t *testing.T) {
	server := configureServer(t)
	client := &Client{
		HTTP:        retryablehttp.NewClient(),
		BaseURL:     server.URL,
		Credentials: &Credentials{},
	}
	client.HTTP.HTTPClient = &http.Client{
		Transport: &http.Transport{MaxIdleConnsPerHost: maxIdleConnections},
		Timeout:   time.Duration(requestTimeout) * time.Second,
	}
	client.HTTP.CheckRetry = RetryPolicy
	client.HTTP.Backoff = func(min, max time.Duration, attemptNum int, resp *http.Response) time.Duration {
		return time.Duration(1)
	}
	req, _ := http.NewRequest(http.MethodGet, server.URL+"/v1/patch_200", nil)
	resp, err := client.SendRequest(req)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Nil(t, err)
	assert.Equal(t, 1, retryCounter)
	req, _ = http.NewRequest(http.MethodGet, server.URL+"/v1/patch_400", nil)
	resp, err = client.SendRequest(req)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, 2, retryCounter)
	req, _ = http.NewRequest(http.MethodGet, server.URL+"/v1/patch_409", nil)
	resp, err = client.SendRequest(req)
	assert.Nil(t, resp)
	assert.NotNil(t, err)
	assert.Equal(t, 7, retryCounter)
}

func configureServer(t *testing.T) *httptest.Server {
	tokenCounter := 1
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		accessToken := fmt.Sprintf("token-%d", tokenCounter)
		switch req.URL.String() {
		case "/v1/oauth/token":
			tokenCounter++
			defer req.Body.Close()
			token := &Token{
				Token:     accessToken,
				Expiry:    300,
				TokenType: "access token",
			}
			response, _ := json.Marshal(token)
			rw.Write(response)
		case "/v1/test":
			assert.Regexp(t, regexp.MustCompile("Bearer token-[\\d]"), req.Header.Get("Authorization"))
			rw.Write([]byte("ok"))
		case "/v1/patch_409":
			errorResponse := &ErrorResponse{
				Detail: "error details",
				Status: 409,
				Title:  "title",
				Type:   "type",
			}
			bytesRes, _ := json.Marshal(errorResponse)
			retryCounter++
			rw.WriteHeader(http.StatusConflict)
			rw.Write(bytesRes)
		case "/v1/patch_200":
			retryCounter++
			rw.Write([]byte("ok"))
		case "/v1/patch_400":
			errorResponse := &ErrorResponse{
				Detail: "error details",
				Status: 400,
				Title:  "title",
				Type:   "type",
			}
			bytesRes, _ := json.Marshal(errorResponse)
			retryCounter++
			rw.WriteHeader(http.StatusBadRequest)
			rw.Write(bytesRes)
		}
	}))
	return server
}
