package client

import (
	"context"
	"fmt"
)

func AssignNetworkElementAlias(ctx context.Context, c *Client, neID, alias string) error {
	url := fmt.Sprintf("%s/%s/%s/aliases/%s", c.BaseURL, networkElementsEndpoint, neID, alias)
	_, err := c.Put(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNetworkElementAlias(ctx context.Context, c *Client, neID, alias string) error {
	url := fmt.Sprintf("%s/%s/%s/aliases/%s", c.BaseURL, networkElementsEndpoint, neID, alias)
	_, err := c.Delete(ctx, url, nil)
	if err != nil {
		return err
	}
	return nil
}

func AliasExists(ctx context.Context, c *Client, neID, alias string) (bool, error) {
	ne, err := GetNetworkElement(ctx, c, neID)
	if err != nil {
		return false, err
	}
	if Contains(alias, ne.Aliases) {
		return true, nil
	}
	return false, nil
}
