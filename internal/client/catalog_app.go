package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	u "net/url"
)

const catalogAppsEndpoint string = "v1/catalog_apps"

type TenantAwarenessData struct {
	TenantCorpIdSupport bool `json:"tenant_corp_id_support"`
	TenantTypeSupport   bool `json:"tenant_type_support"`
}
type Attributes struct {
	TenantAwarenessData TenantAwarenessData `json:"tenant_awareness_data"`
}

type CatalogApp struct {
	ID         string     `json:"id"`
	Name       string     `json:"name"`
	Category   string     `json:"category"`
	Risk       int        `json:"risk"`
	Urls       []string   `json:"urls"`
	Vendor     string     `json:"vendor"`
	Verified   bool       `json:"verified"`
	Attributes Attributes `json:"attributes"`
}

func GetCatalogAppByName(ctx context.Context, c *Client, name, category string) (*CatalogApp, error) {
	var res []CatalogApp
	url := fmt.Sprintf("%s/%s", c.BaseURL, catalogAppsEndpoint)
	resp, err := c.Get(ctx, url, u.Values{"query": {name}})
	if err != nil {
		return nil, fmt.Errorf("could not get catalog app %s: %v", name, err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("could not read location response body")
	}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return nil, fmt.Errorf("could not parse catalog apps response: %v", err)
	}
	for _, catalogApp := range res {
		if catalogApp.Name == name && (category == "" || catalogApp.Category == category) {
			return &catalogApp, nil
		}
	}
	return nil, fmt.Errorf("could not find catalog app with the name %s", name)
}
