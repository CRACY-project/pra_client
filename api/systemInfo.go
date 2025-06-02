package api

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jimbersoftware/pra_client/sysinfo"
)

func (c *Client) SendSystemInfo(info sysinfo.SystemInfo) error {
	payload, err := json.Marshal(info)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", systemInfoURL, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-token", c.token)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		return &httpError{Status: resp.Status, Code: resp.StatusCode}
	}

	return nil
}
