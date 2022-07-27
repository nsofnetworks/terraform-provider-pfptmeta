package client

import (
	"context"
	"fmt"
	"strings"
)

func AssignAlias(ctx context.Context, c *Client, entityID, alias string) error {
	path := networkElementPathByPrefix(entityID)
	url := fmt.Sprintf("%s/%s/%s/aliases/%s", c.BaseURL, path, entityID, alias)
	_, err := c.Put(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}

func DeleteAlias(ctx context.Context, c *Client, entityID, alias string) error {
	path := networkElementPathByPrefix(entityID)
	url := fmt.Sprintf("%s/%s/%s/aliases/%s", c.BaseURL, path, entityID, alias)
	_, err := c.Delete(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}

func AliasExists(ctx context.Context, c *Client, neID, alias string) (bool, error) {
	prefix := strings.Split(neID, "-")[0]
	var aliases []string
	if prefix == "dev" {
		device, err := GetDevice(ctx, c, neID)
		if err != nil {
			return false, err
		}
		aliases = device.Aliases
	} else {
		ne, err := GetNetworkElement(ctx, c, neID)
		if err != nil {
			return false, err
		}
		aliases = ne.Aliases
	}
	if Contains(alias, aliases) {
		return true, nil
	}
	return false, nil
}
