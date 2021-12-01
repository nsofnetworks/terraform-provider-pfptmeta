package client

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type Tag struct {
	Name  string `json:"name"`
	Value string `json:"value"`
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
