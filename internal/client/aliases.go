package client

import (
	"fmt"
)

func AssignNetworkElementAlias(c *Client, neID, alias string) error {
	url := fmt.Sprintf("%s/%s/%s/aliases/%s", c.BaseURL, networkElementsEndpoint, neID, alias)
	_, err := c.Put(url, nil)
	if err != nil {
		return err
	}
	return nil
}

func DeleteNetworkElementAlias(c *Client, neID, alias string) error {
	url := fmt.Sprintf("%s/%s/%s/aliases/%s", c.BaseURL, networkElementsEndpoint, neID, alias)
	_, err := c.Delete(url, nil)
	if err != nil {
		return err
	}
	return nil
}

func AliasExists(c *Client, neID, alias string) (bool, error) {
	ne, err := GetNetworkElement(c, neID)
	if err != nil {
		return false, err
	}
	if Contains(alias, ne.Aliases) {
		return true, nil
	}
	return false, nil
}
