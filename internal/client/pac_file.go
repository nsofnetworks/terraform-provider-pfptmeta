package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
)

const pacFilesEndpoint string = "v1/pac_files"

type PacFile struct {
	ID            string   `json:"id,omitempty"`
	Name          string   `json:"name,omitempty"`
	Description   string   `json:"description"`
	Enabled       bool     `json:"enabled"`
	ApplyToOrg    bool     `json:"apply_to_org"`
	Sources       []string `json:"sources"`
	ExemptSources []string `json:"exempt_sources"`
	HasContent    bool     `json:"has_content,omitempty"`
	Priority      int      `json:"priority"`
}

func NewPacFile(d *schema.ResourceData) *PacFile {
	res := &PacFile{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.ApplyToOrg = d.Get("apply_to_org").(bool)
	res.Enabled = d.Get("enabled").(bool)
	res.Sources = ConfigToStringSlice("sources", d)
	res.ExemptSources = ConfigToStringSlice("exempt_sources", d)
	res.Priority = d.Get("priority").(int)
	return res
}

func parsePacFile(resp *http.Response) (*PacFile, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	pg := &PacFile{}
	err = json.Unmarshal(body, pg)
	if err != nil {
		return nil, fmt.Errorf("could not parse pac file response: %v", err)
	}
	return pg, nil
}

func CreatePacFile(ctx context.Context, c *Client, pf *PacFile) (*PacFile, error) {
	pfUrl := fmt.Sprintf("%s/%s", c.BaseURL, pacFilesEndpoint)
	body, err := json.Marshal(pf)
	if err != nil {
		return nil, fmt.Errorf("could not convert pac file to json: %v", err)
	}
	resp, err := c.Post(ctx, pfUrl, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parsePacFile(resp)
}

func UpdatePacFile(ctx context.Context, c *Client, pfID string, pf *PacFile) (*PacFile, error) {
	pfUrl := fmt.Sprintf("%s/%s/%s", c.BaseURL, pacFilesEndpoint, pfID)
	body, err := json.Marshal(pf)
	if err != nil {
		return nil, fmt.Errorf("could not convert pac file to json: %v", err)
	}
	resp, err := c.Patch(ctx, pfUrl, bytes.NewReader(body))
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
	defer resp.Body.Close()
	decodedResp, err := ioutil.ReadAll(resp.Body)
	strResp := string(decodedResp)
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
