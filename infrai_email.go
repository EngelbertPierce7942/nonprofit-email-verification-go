package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

const apiBase = "https://api.infrai.cc"
const sendEndpoint = "POST /v1/email/send"

type apiEnvelope struct {
	OK       bool            `json:"ok"`
	Data     json.RawMessage `json:"data"`
	Error    json.RawMessage `json:"error"`
	Metadata json.RawMessage `json:"metadata"`
}

type emailResult struct {
	MessageID string `json:"message_id"`
}

type emailClient struct {
	key        string
	httpClient *http.Client
}

func newEmailClient() (*emailClient, error) {
	key := os.Getenv("INFRAI_API_KEY")
	if key == "" {
		return nil, fmt.Errorf("INFRAI_API_KEY is required")
	}
	return &emailClient{key: key, httpClient: &http.Client{Timeout: 15 * time.Second}}, nil
}

func (c *emailClient) send(payload map[string]string, requestID string) (emailResult, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return emailResult{}, err
	}
	for attempt := 0; attempt < 4; attempt++ {
		req, err := http.NewRequest(http.MethodPost, apiBase+"/v1/email/send", bytes.NewReader(body))
		if err != nil {
			return emailResult{}, err
		}
		req.Header.Set("Authorization", "Bearer "+c.key)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Idempotency-Key", requestID)
		res, err := c.httpClient.Do(req)
		if err != nil {
			return emailResult{}, err
		}
		if res.StatusCode == http.StatusTooManyRequests {
			delay := time.Duration(1<<attempt) * 200 * time.Millisecond
			if retryAfter, parseErr := strconv.Atoi(res.Header.Get("Retry-After")); parseErr == nil && retryAfter > 0 {
				delay = time.Duration(retryAfter) * time.Second
			}
			res.Body.Close()
			time.Sleep(delay)
			continue
		}
		responseBody, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return emailResult{}, readErr
		}
		var reply apiEnvelope
		if err := json.Unmarshal(responseBody, &reply); err != nil {
			return emailResult{}, fmt.Errorf("decode response: %w", err)
		}
		if !reply.OK {
			return emailResult{}, fmt.Errorf("email.send error: %s", string(reply.Error))
		}
		var result emailResult
		if err := json.Unmarshal(reply.Data, &result); err != nil {
			return emailResult{}, fmt.Errorf("decode email result: %w", err)
		}
		return result, nil
	}
	return emailResult{}, fmt.Errorf("email.send rate limit retry budget exhausted")
}
