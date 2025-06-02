package api

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"

	"github.com/jimbersoftware/pra_client/libraryinfo"
	"github.com/jimbersoftware/pra_client/logging"
)

const libraryInfoURL = "https://pra.testing.jimber.io/api/v1/pdes/library-info"

func (c *Client) SendLibraryInfo(libs []libraryinfo.PackageInfo) error {
	payload, err := json.Marshal(libs)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", libraryInfoURL, bytes.NewBuffer(payload))
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
		logging.Log(logging.ERROR, "Failed to send library information:", resp.Status)
		// Log the response body for debugging purposes
		body, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			logging.Log(logging.ERROR, "Error reading response body:", body)
			return readErr
		}

		return &httpError{Status: resp.Status, Code: resp.StatusCode}
	}
	logging.Log(logging.INFO, "Library information response body", resp.Status)
	if resp.StatusCode != http.StatusOK {
		logging.Log(logging.ERROR, "Failed to send library information:", resp.Status)
		return &httpError{Status: resp.Status, Code: resp.StatusCode}
	}
	return nil
}
