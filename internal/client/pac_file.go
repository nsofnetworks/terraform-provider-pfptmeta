package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"net/http"
)

const pacFilesEndpoint string = "v1/pac_files"
const pacTypeManaged string = "managed"
const pacTypeBringYourOwn string = "bring_your_own"

type ManagedContent struct {
	Domains    []string `json:"domains,omitempty"`
	CloudApps  []string `json:"cloud_apps,omitempty"`
	IpNetworks []string `json:"ip_networks,omitempty"`
}

type PacFile struct {
	ID             string          `json:"id,omitempty"`
	Name           string          `json:"name,omitempty"`
	Description    string          `json:"description"`
	Enabled        bool            `json:"enabled"`
	ApplyToOrg     bool            `json:"apply_to_org"`
	Sources        []string        `json:"sources"`
	ExemptSources  []string        `json:"exempt_sources"`
	HasContent     bool            `json:"has_content,omitempty"`
	Priority       int             `json:"priority"`
	Type           string          `json:"type,omitempty"`
	ManagedContent *ManagedContent `json:"managed_content,omitempty"`
}

func NewManagedContent(d *schema.ResourceData) *ManagedContent {
	mc := d.Get("managed_content").([]interface{})
	fmt.Printf("[NADAV] NewManagedContent [0]: mc: %+v\n", mc)
	if len(mc) == 0 {
		return nil
	}
	fmt.Printf("[NADAV] NewManagedContent [1]\n")
	conf := mc[0].(map[string]interface{})
	fmt.Printf("[NADAV] NewManagedContent [2] conf: %+v\n", conf)
	res := &ManagedContent{}
	rng := conf["domains"].([]interface{})
	if len(rng) > 0 {
		fmt.Printf("[NADAV] NewManagedContent [2.1]\n")
		res.Domains = make([]string, len(rng))
		for i, val := range rng {
			res.Domains[i] = val.(string)
		}
	}
	fmt.Printf("[NADAV] NewManagedContent [3]\n")
	rng = conf["cloud_apps"].([]interface{})
	if len(rng) > 0 {
		fmt.Printf("[NADAV] NewManagedContent [3.1]\n")
		res.CloudApps = make([]string, len(rng))
		for i, val := range rng {
			res.CloudApps[i] = val.(string)
		}
	}
	fmt.Printf("[NADAV] NewManagedContent [4]\n")
	rng = conf["ip_networks"].([]interface{})
	if len(rng) > 0 {
		fmt.Printf("[NADAV] NewManagedContent [4.1]\n")
		res.IpNetworks = make([]string, len(rng))
		for i, val := range rng {
			res.IpNetworks[i] = val.(string)
		}
	}
	fmt.Printf("[NADAV] NewManagedContent [5] res: %+v\n\n", res)
	return res
}

func PacFileBase(d *schema.ResourceData) *PacFile {
	res := &PacFile{}
	res.Description = d.Get("description").(string)
	res.ApplyToOrg = d.Get("apply_to_org").(bool)
	res.Enabled = d.Get("enabled").(bool)
	res.Sources = ConfigToStringSlice("sources", d)
	res.ExemptSources = ConfigToStringSlice("exempt_sources", d)
	res.Priority = d.Get("priority").(int)
	return res
}

func NewPacFile(d *schema.ResourceData) *PacFile {
	res := PacFileBase(d)
	res.Name = d.Get("name").(string)
	res.Type = d.Get("type").(string)
	return res
}

func ModifiedPacFile(d *schema.ResourceData) *PacFile {
	res := PacFileBase(d)
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	return res
}

func parsePacFile(resp []byte) (*PacFile, error) {
	pg := &PacFile{}
	err := json.Unmarshal(resp, pg)
	if err != nil {
		return nil, fmt.Errorf("could not parse PAC file response: %v", err)
	}
	return pg, nil
}

func parseManagedContent(resp []byte) (*ManagedContent, error) {
	mc := &ManagedContent{}
	err := json.Unmarshal(resp, mc)
	if err != nil {
		return nil, fmt.Errorf("could not parse managed content response: %v", err)
	}
	return mc, nil
}

func CreatePacFile(ctx context.Context, c *Client, pf *PacFile) (*PacFile, error) {
	pfUrl := fmt.Sprintf("%s/%s", c.BaseURL, pacFilesEndpoint)
	body, err := json.Marshal(pf)
	if err != nil {
		return nil, fmt.Errorf("could not convert pac file to json: %v", err)
	}
	resp, err := c.Post(ctx, pfUrl, body)
	if err != nil {
		return nil, err
	}
	return parsePacFile(resp)
}

func UpdatePacFile(ctx context.Context, c *Client, pfID string, pf *PacFile) (*PacFile, error) {
	pfUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, pacFilesEndpoint, pfID)
	body, err := json.Marshal(pf)
	if err != nil {
		return nil, fmt.Errorf("could not convert PAC file to json: %v", err)
	}
	resp, err := c.Patch(ctx, pfUrl, body)
	if err != nil {
		return nil, err
	}
	return parsePacFile(resp)
}

func GetPacFile(ctx context.Context, c *Client, pfID string) (*PacFile, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, pacFilesEndpoint, pfID)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parsePacFile(resp)
}

func DeletePacFile(ctx context.Context, c *Client, pfID string) (*PacFile, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, pacFilesEndpoint, pfID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parsePacFile(resp)
}

func GetPacFileContent(ctx context.Context, c *Client, pfID string) (*string, error) {
	pfUrl := fmt.Sprintf("%s/%s/%s/content", c.BaseURL, pacFilesEndpoint, pfID)
	resp, err := c.Get(ctx, pfUrl, nil)
	if err != nil {
		return nil, err
	}
	strResp := string(resp)
	return &strResp, nil
}

func DeletePacFileContent(ctx context.Context, c *Client, pfID string) error {
	url := fmt.Sprintf("%s/%s/%s/content", c.BaseURL, pacFilesEndpoint, pfID)
	_, err := c.Delete(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}

func PutPacFileContent(ctx context.Context, c *Client, pfID, pfContent string) error {
	pfUrl := fmt.Sprintf("%s/%s/%s/content", c.BaseURL, pacFilesEndpoint, pfID)
	req, err := http.NewRequestWithContext(ctx, http.MethodPut, pfUrl, bytes.NewReader([]byte(pfContent)))
	req.Header.Set("Content-Type", "text/plain")
	if err != nil {
		return err
	}
	_, err = c.SendRequest(req)
	if err != nil {
		return err
	}
	return nil
}

func GetPacFileManagedContent(ctx context.Context, c *Client, pfID string) (*ManagedContent, error) {
	pfUrl := fmt.Sprintf("%s/%s/%s/content/managed", c.BaseURL, pacFilesEndpoint, pfID)
	resp, err := c.Get(ctx, pfUrl, nil)
	if err != nil {
		return nil, err
	}
	return parseManagedContent(resp)
}

func PatchPacFileManagedContent(ctx context.Context, c *Client, pfID string, mc *ManagedContent) error {
	pfUrl := fmt.Sprintf("%s/%s/%s/content/managed", c.BaseURL, pacFilesEndpoint, pfID)
	body, err := json.Marshal(mc)
	fmt.Printf("[NADAV] PatchPacFileManagedContent [0]: body: %s\n", body)
	if err != nil {
		fmt.Println("could not convert PAC file managed content to json")
		return err
	}

	_, err = c.Patch(ctx, pfUrl, body)
	if err != nil {
		return err
	}
	return nil
}
