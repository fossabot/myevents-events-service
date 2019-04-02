package rest

import (
	"encoding/hex"
	"encoding/json"
	"github.com/danielpacak/myevents-events-service/domain"
	"net/http"
	"net/url"
)

type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

func NewClient(baseURL *url.URL) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: http.DefaultClient,
	}
}

func (c *Client) ListEvents() ([]domain.Event, error) {
	rel := &url.URL{Path: "/events/"}
	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var users []domain.Event
	err = json.NewDecoder(resp.Body).Decode(&users)

	return users, nil
}

func (c *Client) GetById(id []byte) (*domain.Event, error) {
	rel := &url.URL{Path: "/events/id/" + hex.EncodeToString(id)}
	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var event domain.Event
	err = json.NewDecoder(resp.Body).Decode(&event)

	return &event, nil
}

func (c *Client) GetByName(name string) (*domain.Event, error) {
	rel := &url.URL{Path: "/events/name/" + name}
	u := c.baseURL.ResolveReference(rel)

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var event domain.Event
	err = json.NewDecoder(resp.Body).Decode(&event)

	return &event, nil
}
