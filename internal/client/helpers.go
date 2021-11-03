package client

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type ResponseOnlyString string

// MarshalJSON When setting a new mapped domain or MappedHost to network element
// the name goes in the query params and not in the body.
// This method makes sure that when marshaling a MappedDomain it will not include the name
func (ResponseOnlyString) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

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

func resourceTypeSetToStringSlice(s *schema.Set) []string {
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
