package client

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsKeyExluded(t *testing.T) {
	assert.Equal(t, true, isKeyExluded("key", []string{"key", "other-key"}))
	assert.Equal(t, false, isKeyExluded("key", []string{"key1", "key2"}))
}

func TestResourceTypeSetToStringSlice(t *testing.T) {
	s := &schema.Set{F: schema.HashString}
	for _, i := range []string{"element1", "element2", "element3"} {
		s.Add(i)
	}
	assert.Contains(t, resourceTypeSetToStringSlice(s), "element1")
	assert.Contains(t, resourceTypeSetToStringSlice(s), "element2")
	assert.Contains(t, resourceTypeSetToStringSlice(s), "element3")
}

func TestConvertTagsListToMap(t *testing.T) {
	tags := []Tag{{"tag-name-1", "tag-value-1"}, {"tag-name-2", "tag-value-2"}}
	expected := map[string]string{
		"tag-name-1": "tag-value-1",
		"tag-name-2": "tag-value-2",
	}
	assert.Equal(t, expected, ConvertTagsListToMap(tags))
}
