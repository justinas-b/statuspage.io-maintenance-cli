package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

func NewClient(apiKey string) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://api.statuspage.io",
		HTTPClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (c *Client) GetPage(pageId string) (Page, error) {
	URL := fmt.Sprintf("%s/v1/pages/%s", c.BaseURL, pageId)

	resp := c.doRequest("GET", URL, nil)
	if resp.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("unexpected response status %q", resp.Status)
	}
	defer resp.Body.Close()

	obj, err := parseResponse[Page](resp.Body)
	if err != nil {
		return Page{}, fmt.Errorf("invalid API response %s: %w", resp.Body, err)
	}
	obj.client = c
	return obj, nil
}

func (c *Client) GetPages() ([]*Page, error) {
	URL := fmt.Sprintf("%s/v1/pages", c.BaseURL)

	resp := c.doRequest("GET", URL, nil)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response status %q", resp.Status)
	}
	defer resp.Body.Close()

	pages, err := parseResponse[[]Page](resp.Body)
	if err != nil {
		return nil, fmt.Errorf("invalid API response %s: %w", resp.Body, err)
	}

	for idx := range pages {
		pages[idx].client = c
	}
	c.Pages = pages
	return PointersOf(c.Pages).([]*Page), nil
}

func (c *Client) SetMaintenance(pageID string, name string, body string, startTime time.Time, endTime time.Time) error {
	URL := fmt.Sprintf("%s/v1/pages/%s/incidents", c.BaseURL, pageID)
	incident := Incident{
		Name:                             name,
		Status:                           "scheduled",
		ImpactOverride:                   "maintenance",
		ScheduledAutoInProgress:          true,
		ScheduledAutoCompleted:           true,
		AutoTransitionToMaintenanceState: true, // Check if actually needed
		AutoTransitionToOperationalState: true, // Check if actually needed
		Body:                             body,
		ScheduledFor:                     startTime.Format(time.RFC3339Nano), // "2013-05-07T03:00:00.007Z"
		ScheduledUntil:                   endTime.Format(time.RFC3339Nano),   // "2013-05-07T03:00:00.007Z"
	}
	incidentBody := map[string]Incident{"incident": incident}
	payload, _ := json.Marshal(incidentBody)
	resp := c.doRequest("POST", URL, bytes.NewBuffer(payload))

	if resp.StatusCode != http.StatusCreated {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}

		var r map[string][]string
		err = json.Unmarshal(b, &r)
		if err != nil {
			panic(err)
		}

		return fmt.Errorf(r["error"][0])

		//b, err := httputil.DumpResponse(resp, false)
		//if err != nil {
		//	log.Fatalln(err)
		//}
		//return fmt.Errorf("unexpected response status %q", b)
	}
	defer resp.Body.Close()

	return nil
}

func (p *Page) SetMaintenance(name string, body string, startTime time.Time, endTime time.Time) error {
	return p.client.SetMaintenance(
		p.ID,
		name,
		body,
		startTime,
		endTime,
	)
}

func (p *Page) String() string {
	return p.Name
}
