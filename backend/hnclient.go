package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Untuk memanggil/komunikasi Hacker News API

const baseURL = "https://hacker-news.firebaseio.com/v0"

type HNClient struct {
	httpClient *http.Client
}

func NewHNClient(timeout time.Duration) *HNClient {
	return &HNClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// helper: fetch raw JSON and decode into target
func (c *HNClient) fetchJSON(ctx context.Context, url string, target interface{}) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/json")
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 400 {
		b, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("status %d: %s", resp.StatusCode, string(b))
	}
	dec := json.NewDecoder(resp.Body)
	return dec.Decode(target)
}

// GetTopStoryIDs returns slice of ids (top 100 typically)
// mendapat list ID
func (c *HNClient) GetTopStoryIDs(ctx context.Context) ([]int, error) {
	url := baseURL + "/topstories.json"
	var ids []int
	if err := c.fetchJSON(ctx, url, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// GetItem fetches /v0/item/{id}.json
// mengambil detail tiap item
func (c *HNClient) GetItem(ctx context.Context, id int) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/item/%d.json", baseURL, id)
	var item map[string]interface{}
	if err := c.fetchJSON(ctx, url, &item); err != nil {
		return nil, err
	}
	return item, nil
}

// GetUser
func (c *HNClient) GetUser(ctx context.Context, user string) (map[string]interface{}, error) {
	url := fmt.Sprintf("%s/user/%s.json", baseURL, user)
	var u map[string]interface{}
	if err := c.fetchJSON(ctx, url, &u); err != nil {
		return nil, err
	}
	return u, nil
}
