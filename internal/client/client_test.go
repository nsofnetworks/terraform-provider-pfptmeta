package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"regexp"
	"testing"
	"time"
)

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
		HTTPClient: &http.Client{
			Transport: &http.Transport{
				MaxIdleConnsPerHost: maxIdleConnections,
			},
			Timeout: time.Duration(requestTimeout) * time.Second,
		},
		BaseURL:     server.URL,
		Credentials: &Credentials{},
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
		}
	}))
	return server
}
