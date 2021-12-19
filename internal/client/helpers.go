package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

func NewTags(d *schema.ResourceData) []Tag {
	rawTags := d.Get("tags").(map[string]interface{})
	tags := make([]Tag, len(rawTags))
	index := 0
	for key, value := range rawTags {
		t := Tag{
			Name:  key,
			Value: value.(string),
		}
		tags[index] = t
		index += 1
	}
	return tags
}

func AssignTagsToResource(ctx context.Context, c *Client, rID, rName string, tags []Tag) error {
	body, err := json.Marshal(tags)
	if err != nil {
		return fmt.Errorf("could not convert tags to json: %v", err)
	}
	url := fmt.Sprintf("%s/v1/%s/%s/tags", c.BaseURL, rName, rID)
	_, err = c.Put(ctx, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	return nil
}

func MapResponseToResource(r interface{}, d *schema.ResourceData, excludedKeys []string) error {
	var interfaceResponse map[string]interface{}
	marshaledResponse, err := json.Marshal(r)
	if err != nil {
		return err
	}
	err = json.Unmarshal(marshaledResponse, &interfaceResponse)
	if err != nil {
		return err
	}
	for key, val := range interfaceResponse {
		if isKeyExluded(key, excludedKeys) {
			continue
		}
		err = d.Set(key, val)
		if err != nil {
			return err
		}
	}
	return nil
}

func isKeyExluded(key string, excludedKeys []string) bool {
	for _, excludedKey := range excludedKeys {
		if key == excludedKey {
			return true
		}
	}
	return false
}

func ResourceTypeSetToStringSlice(s *schema.Set) []string {
	valuesList := s.List()
	values := make([]string, len(valuesList))
	for i := 0; i < len(valuesList); i++ {
		values[i] = fmt.Sprint(valuesList[i])
	}
	return values
}

func ConfigToStringSlice(key string, d *schema.ResourceData) []string {
	data := d.Get(key).([]interface{})
	res := make([]string, len(data))
	for i, val := range data {
		res[i] = val.(string)
	}
	return res
}

func ConvertTagsListToMap(tags []Tag) map[string]string {
	res := make(map[string]string)
	for _, tag := range tags {
		res[tag.Name] = tag.Value
	}
	return res
}

func Contains(v string, a []string) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}
	return false
}
