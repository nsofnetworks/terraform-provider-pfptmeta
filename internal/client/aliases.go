package client

import (
	"context"
	"fmt"
	"strings"
)

func AssignNetworkElementAlias(ctx context.Context, c *Client, neID, alias string) error {
	path := networkElementPathByPrefix(neID)
	url := fmt.Sprintf("%s/%s/%s/aliases/%s", c.BaseURL, path, neID, alias)
	_, err := c.Put(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNetworkElementAlias(ctx context.Context, c *Client, neID, alias string) error {
	path := networkElementPathByPrefix(neID)
	url := fmt.Sprintf("%s/%s/%s/aliases/%s", c.BaseURL, path, neID, alias)
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
