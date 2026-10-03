package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

const (
	DefaultHost = "localhost"
	DefaultPort = 10350
)

var DefaultAddr = fmt.Sprintf("%s:%d", DefaultHost, DefaultPort)

type Client struct {
	addr string
	http *http.Client
}

func NewClient(addr string) *Client {
	if addr == "" {
		addr = DefaultAddr
	}
	return &Client{
		addr: addr,
		http: &http.Client{Timeout: 5 * time.Second},
	}
}

func (c *Client) FetchView() (*View, error) {
	resp, err := c.http.Get(fmt.Sprintf("http://%s/api/view", c.addr))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tilt API returned %d", resp.StatusCode)
	}

	var vr ViewResponse
	if err := json.NewDecoder(resp.Body).Decode(&vr); err != nil {
		return nil, err
	}
	return &vr.View, nil
}
