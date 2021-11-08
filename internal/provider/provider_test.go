package provider

import (
	"context"
	"encoding/json"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
	"github.com/nsofnetworks/terraform-provider-pfptmeta/internal/client"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"os/user"
	"path/filepath"
	"testing"
)

var provider *schema.Provider

// providerFactories are used to instantiate a provider during acceptance testing.
// The factory function will be invoked for every Terraform CLI command executed
// to create a provider server to which the CLI can reattach.
var providerFactories = map[string]func() (*schema.Provider, error){
	"pfptmeta": func() (*schema.Provider, error) {
		if provider == nil {
			provider = New("dev")()
		}
		return provider, nil
	},
}

func TestProvider(t *testing.T) {
	if err := New("dev")().InternalValidate(); err != nil {
		t.Fatalf("err: %s", err)
	}
}

func testAccPreCheck(t *testing.T) {
	if os.Getenv("PFPTMETA_API_KEY") == "" {
		t.Fatalf("PFPTMETA_API_KEY env var must be set")
	}
	if os.Getenv("PFPTMETA_API_SECRET") == "" {
		t.Fatalf("PFPTMETA_API_SECRET env var must be set")
	}
	if os.Getenv("PFPTMETA_ORG_SHORTNAME") == "" {
		t.Fatalf("PFPTMETA_ORG_SHORTNAME env var must be set")
	}

}

func TestConfigure(t *testing.T) {
	server := configureAuthServer(t)
	setEnvVar(t, "PFPTMETA_BASE_URL", server.URL)
	defer os.Unsetenv("PFPTMETA_BASE_URL")
	setEnvVar(t, "PFPTMETA_API_KEY", "api-key-from-env-var")
	defer os.Unsetenv("PFPTMETA_API_KEY")
	setEnvVar(t, "PFPTMETA_API_SECRET", "api-secret-from-env-var")
	defer os.Unsetenv("PFPTMETA_API_SECRET")
	setEnvVar(t, "PFPTMETA_ORG_SHORTNAME", "org-from-env-var")
	defer os.Unsetenv("PFPTMETA_ORG_SHORTNAME")
	cases := map[string]struct {
		P      *schema.Provider
		Config map[string]interface{}
	}{
		"auth-in-config": {
			P: New("dev")(),
			Config: map[string]interface{}{
				"api_key":       "api-key-from-conf",
				"api_secret":    "api-secret-from-conf",
				"org_shortname": "org-from-conf",
			},
		},
		"auth-in-env-var": {
			P:      New("dev")(),
			Config: nil,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			c := terraform.NewResourceConfigRaw(tc.Config)
			diags := tc.P.Configure(context.Background(), c)
			if diags.HasError() {
				t.Errorf("%s failed: %+v", name, diags[0])
			}
		})
	}
}

func TestConfigureFromFile(t *testing.T) {
	server := configureAuthServer(t)
	setEnvVar(t, "PFPTMETA_BASE_URL", server.URL)
	defer os.Unsetenv("PFPTMETA_BASE_URL")
	usr, err := user.Current()
	if err != nil {
		t.Errorf("could not find current user's name: %v", err)
	}
	directoryPath := filepath.Join(usr.HomeDir, ".pfptmeta")
	_ = os.Mkdir(directoryPath, os.ModePerm)
	path := filepath.Join(usr.HomeDir, ".pfptmeta", "credentials.json")
	config := &client.Config{
		APIKey:       "api-key",
		APISecret:    "api-secret",
		OrgShortname: "org",
	}
	configBytes, _ := json.Marshal(config)
	err = os.WriteFile(path, configBytes, 0644)
	if err != nil {
		t.Fatalf("could not write config file: %v", err)
	}
	defer os.Remove(path)
	provider := New("dev")()
	c := terraform.NewResourceConfigRaw(nil)
	diags := provider.Configure(context.Background(), c)
	if diags.HasError() {
		t.Errorf("failed: %+v", diags[0])
	}
}

func configureAuthServer(t *testing.T) *httptest.Server {
	server := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.String() == "/v1/oauth/token" {
			defer req.Body.Close()
			credsReq := &client.Credentials{}
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				t.Errorf("Could not read token request: %v", err)
			}
			err = json.Unmarshal(body, credsReq)
			t.Logf("Got auth request with: %+v", credsReq)
			assert.Nil(t, err)
			assert.NotEqual(t, "", credsReq.Scope)
			assert.NotEqual(t, "", credsReq.ClientID)
			assert.NotEqual(t, "", credsReq.ClientSecret)
			assert.NotEqual(t, "", credsReq.GrantType)
			token := &client.Token{
				Token:     "token",
				Expiry:    300,
				TokenType: "access token",
			}
			response, _ := json.Marshal(token)
			rw.Write(response)
		}
	}))
	return server
}

func setEnvVar(t *testing.T, key, value string) {
	err := os.Setenv(key, value)
	if err != nil {
		t.Fatalf("Could not set %s env variable to %s: %v", key, value, err)
	}
}
