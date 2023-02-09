package client

import (
	"context"
	"encoding/json"
	"fmt"
)

const (
	locationsEndpoint string = "v1/locations"
)

type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Name      string  `json:"name"`
	State     string  `json:"state"`
	Status    string  `json:"status"`
}

func GetLocation(ctx context.Context, c *Client, lName string) (*Location, error) {
	location := &Location{}
	url := fmt.Sprintf("%s/%s/%s", c.BaseURL, locationsEndpoint, lName)
	resp, err := c.Get(ctx, url, nil)
	if err != nil {
		return nil, fmt.Errorf("could not get location %s: %v", lName, err)
	}
	err = json.Unmarshal(resp, location)
	if err != nil {
		return nil, fmt.Errorf("could not parse location response: %v", err)
	}
	return location, nil
}
