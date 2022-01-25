package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"io/ioutil"
	"net/http"
	u "net/url"
)

const accessBridgeEndpoint = "v1/access_bridges"

type LdapProvisioningConfig struct {
	Certificate         string `json:"certificate"`
	FollowForeignGroups bool   `json:"follow_foreign_groups"`
	GroupsBaseDn        string `json:"groups_base_dn"`
	GroupsFilter        string `json:"groups_filter"`
	IdpId               string `json:"idp_id"`
	Password            string `json:"password"`
	Type                string `json:"type"`
	Url                 string `json:"url"`
	Username            string `json:"username"`
	UsersBaseDn         string `json:"users_base_dn"`
	UsersFilter         string `json:"users_filter"`
}

func newLdapProvisioningConfig(d *schema.ResourceData) *LdapProvisioningConfig {
	res := &LdapProvisioningConfig{}
	ldapConf := d.Get("ldap_provisioning_config").([]interface{})
	if len(ldapConf) == 0 {
		return nil
	}
	ldapConfMap := ldapConf[0].(map[string]interface{})
	res.Certificate = ldapConfMap["certificate"].(string)
	res.FollowForeignGroups = ldapConfMap["follow_foreign_groups"].(bool)
	res.GroupsBaseDn = ldapConfMap["groups_base_dn"].(string)
	res.GroupsFilter = ldapConfMap["groups_filter"].(string)
	res.IdpId = ldapConfMap["idp_id"].(string)
	res.Password = ldapConfMap["password"].(string)
	res.Type = ldapConfMap["type"].(string)
	res.Url = ldapConfMap["url"].(string)
	res.Username = ldapConfMap["username"].(string)
	res.UsersBaseDn = ldapConfMap["users_base_dn"].(string)
	res.UsersFilter = ldapConfMap["users_filter"].(string)
	return res
}

type ProofpointCasbConfig struct {
	Region   string `json:"region"`
	TenantId string `json:"tenant_id"`
}

func newProofpointCasbConfig(conf []interface{}) *ProofpointCasbConfig {
	res := &ProofpointCasbConfig{}
	if len(conf) == 0 {
		return nil
	}
	confMap := conf[0].(map[string]interface{})
	res.Region = confMap["region"].(string)
	res.TenantId = confMap["tenant_id"].(string)
	return res
}

type QradarHttpConfig struct {
	Certificate string `json:"certificate,omitempty"`
	Url         string `json:"url"`
}

func newQradarHttpConfig(conf []interface{}) *QradarHttpConfig {
	res := &QradarHttpConfig{}
	if len(conf) == 0 {
		return nil
	}
	confMap := conf[0].(map[string]interface{})
	res.Certificate = confMap["certificate"].(string)
	res.Url = confMap["url"].(string)
	return res
}

type S3Config struct {
	Bucket   string `json:"bucket"`
	Compress bool   `json:"compress"`
	Prefix   string `json:"prefix,omitempty"`
}

func newS3Config(conf []interface{}) *S3Config {
	res := &S3Config{}
	if len(conf) == 0 {
		return nil
	}
	confMap := conf[0].(map[string]interface{})
	res.Bucket = confMap["bucket"].(string)
	res.Compress = confMap["compress"].(bool)
	res.Prefix = confMap["prefix"].(string)
	return res
}

type SplunkHttpConfig struct {
	Certificate        string `json:"certificate,omitempty"`
	PubliclyAccessible bool   `json:"publicly_accessible"`
	Token              string `json:"token"`
	Url                string `json:"url"`
}

func newSplunkHttpConfig(conf []interface{}) *SplunkHttpConfig {
	res := &SplunkHttpConfig{}
	if len(conf) == 0 {
		return nil
	}
	confMap := conf[0].(map[string]interface{})
	res.Certificate = confMap["certificate"].(string)
	res.PubliclyAccessible = confMap["publicly_accessible"].(bool)
	res.Token = confMap["token"].(string)
	res.Url = confMap["url"].(string)
	return res
}

type SyslogConfig struct {
	Host  string `json:"host"`
	Port  int    `json:"port"`
	Proto string `json:"proto"`
}

func newSyslogConfig(conf []interface{}) *SyslogConfig {
	res := &SyslogConfig{}
	if len(conf) == 0 {
		return nil
	}
	confMap := conf[0].(map[string]interface{})
	res.Host = confMap["host"].(string)
	res.Port = confMap["port"].(int)
	res.Proto = confMap["proto"].(string)
	return res
}

type SiemConfig struct {
	Type                 string                `json:"type,omitempty"`
	ExportLogs           []string              `json:"export_logs"`
	ProofpointCasbConfig *ProofpointCasbConfig `json:"proofpoint_casb_config"`
	QradarHttpConfig     *QradarHttpConfig     `json:"qradar_http_config"`
	S3Config             *S3Config             `json:"s3_config"`
	SplunkHttpConfig     *SplunkHttpConfig     `json:"splunk_http_config"`
	SyslogConfig         *SyslogConfig         `json:"syslog_config"`
}

func newSiemConfig(d *schema.ResourceData) *SiemConfig {
	res := &SiemConfig{}
	siemConf := d.Get("siem_config").([]interface{})
	if len(siemConf) == 0 {
		return nil
	}
	siemConfMap := siemConf[0].(map[string]interface{})
	el := siemConfMap["export_logs"].([]interface{})
	res.ExportLogs = make([]string, len(el))
	for i, l := range el {
		res.ExportLogs[i] = l.(string)
	}
	res.ProofpointCasbConfig = newProofpointCasbConfig(siemConfMap["proofpoint_casb_config"].([]interface{}))
	res.QradarHttpConfig = newQradarHttpConfig(siemConfMap["qradar_http_config"].([]interface{}))
	res.S3Config = newS3Config(siemConfMap["s3_config"].([]interface{}))
	res.SplunkHttpConfig = newSplunkHttpConfig(siemConfMap["splunk_http_config"].([]interface{}))
	res.SyslogConfig = newSyslogConfig(siemConfMap["syslog_config"].([]interface{}))
	return res
}

type AccessBridge struct {
	ID                     string                  `json:"id,omitempty"`
	Name                   string                  `json:"name,omitempty"`
	Description            string                  `json:"description"`
	Enabled                bool                    `json:"enabled"`
	LdapProvisioningConfig *LdapProvisioningConfig `json:"ldap_provisioning_config,omitempty"`
	SiemConfig             *SiemConfig             `json:"siem_config,omitempty"`
	NotificationChannels   []string                `json:"notification_channels"`
	Status                 string                  `json:"status,omitempty"`
	StatusDescription      string                  `json:"status_description,omitempty"`
	Type                   string                  `json:"type,omitempty"`
}

func NewAccessBridge(d *schema.ResourceData) *AccessBridge {
	res := &AccessBridge{}
	if d.HasChange("name") {
		res.Name = d.Get("name").(string)
	}
	res.Description = d.Get("description").(string)
	res.Enabled = d.Get("enabled").(bool)
	res.NotificationChannels = ConfigToStringSlice("notification_channels", d)
	if _, exists := d.GetOk("ldap_provisioning_config"); exists {
		res.LdapProvisioningConfig = newLdapProvisioningConfig(d)
	}
	if _, exists := d.GetOk("siem_config"); exists {
		res.SiemConfig = newSiemConfig(d)
	}
	return res
}

func parseAccessBridge(resp *http.Response) (*AccessBridge, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	e := &AccessBridge{}
	err = json.Unmarshal(body, e)
	if err != nil {
		return nil, fmt.Errorf("could not parse access bridge response: %v", err)
	}
	return e, nil
}

func CreateAccessBridge(ctx context.Context, c *Client, e *AccessBridge) (*AccessBridge, error) {
	url := fmt.Sprintf("%s/%s", c.BaseURL, accessBridgeEndpoint)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert access bridge to json: %v", err)
	}
	resp, err := c.Post(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseAccessBridge(resp)
}

func GetAccessBridge(ctx context.Context, c *Client, eID string) (*AccessBridge, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, accessBridgeEndpoint, eID)
	resp, err := c.Get(ctx, url, u.Values{"expand": {"true"}})
	if err != nil {
		return nil, err
	}
	return parseAccessBridge(resp)
}

func UpdateAccessBridge(ctx context.Context, c *Client, eID string, e *AccessBridge) (*AccessBridge, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, accessBridgeEndpoint, eID)
	body, err := json.Marshal(e)
	if err != nil {
		return nil, fmt.Errorf("could not convert access bridge to json: %v", err)
	}
	resp, err := c.Patch(ctx, url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	return parseAccessBridge(resp)
}

func DeleteAccessBridge(ctx context.Context, c *Client, mID string) (*AccessBridge, error) {
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, accessBridgeEndpoint, mID)
	resp, err := c.Delete(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	return parseAccessBridge(resp)
}
