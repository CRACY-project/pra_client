package api

import (
	"bytes"
	"net/http"
	"time"

	"github.com/jimbersoftware/pra_client/logging"
)

const (
	heartbeatURL  = "https://pra.testing.jimber.io/api/v1/pdes/heartbeat"
	systemInfoURL = "https://pra.testing.jimber.io/api/v1/pdes/system-info"

	heartbeatInterval = time.Minute
)

type Client struct {
	token  string
	stopCh chan struct{}
}

// NewClient initializes the API client and starts the heartbeat goroutine.
func NewClient(token string) *Client {
	c := &Client{
		token:  token,
		stopCh: make(chan struct{}),
	}
	go c.startHeartbeat()
	return c
}

// startHeartbeat sends a heartbeat every minute until stopCh is closed.
func (c *Client) startHeartbeat() {
	c.sendHeartbeat()
	ticker := time.NewTicker(heartbeatInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := c.sendHeartbeat(); err != nil {
				logging.Log(logging.ERROR, "Failed to send heartbeat:", err)
			}
		case <-c.stopCh:
			return
		}
	}
}

// sendHeartbeat performs the HTTP POST for the heartbeat.
func (c *Client) sendHeartbeat() error {
	logging.Log(logging.INFO, "Sending heartbeat to", heartbeatURL)
	req, err := http.NewRequest("POST", heartbeatURL, bytes.NewBuffer(nil))
	if err != nil {
		return err
	}

	req.Header.Set("accept", "application/json")
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

// Stop terminates the heartbeat goroutine.
func (c *Client) Stop() {
	close(c.stopCh)
}

type httpError struct {
	Status string
	Code   int
}

func (e *httpError) Error() string {
	return "HTTP error: " + e.Status
}
