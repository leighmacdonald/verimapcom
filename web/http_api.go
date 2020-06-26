package web

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

type exampleKML struct {
	ID               int         `json:"id"`
	Name             string      `json:"name"`
	Hash             string      `json:"hash"`
	Sha256           string      `json:"sha256"`
	Ext              string      `json:"ext"`
	Mime             string      `json:"mime"`
	Size             string      `json:"size"`
	URL              string      `json:"url"`
	Provider         string      `json:"provider"`
	ProviderMetadata interface{} `json:"provider_metadata"`
	CreatedAt        time.Time   `json:"created_at"`
	UpdatedAt        time.Time   `json:"updated_at"`
}

type examplePage struct {
	ID           int               `json:"id"`
	Name         string            `json:"name"`
	Description  string            `json:"description"`
	Latitude     float64           `json:"latitude"`
	Longitude    float64           `json:"longitude"`
	Zoom         int               `json:"zoom"`
	Public       bool              `json:"public"`
	Stats        string            `json:"stats"`
	StatsMap     map[string]string `json:"-"`
	CreatedAt    time.Time         `json:"created_at"`
	UpdatedAt    time.Time         `json:"updated_at"`
	Layer        string            `json:"layer"`
	ZoomMin      int               `json:"zoom_min"`
	ZoomMax      int               `json:"zoom_max"`
	VectorLayers interface{}       `json:"vector_layers"`
	Kml          exampleKML        `json:"kml,omitempty"`
}

type Showcase struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Body      string    `json:"body"`
	LinkText  string    `json:"link_text"`
	LinkURL   string    `json:"link_url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Order     int       `json:"order"`
	Image     struct {
		ID               int         `json:"id"`
		Name             string      `json:"name"`
		Hash             string      `json:"hash"`
		Sha256           string      `json:"sha256"`
		Ext              string      `json:"ext"`
		Mime             string      `json:"mime"`
		Size             string      `json:"size"`
		URL              string      `json:"url"`
		Provider         string      `json:"provider"`
		ProviderMetadata interface{} `json:"provider_metadata"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
	} `json:"image"`
}

type FireWatch struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	Published time.Time `json:"published"`
	Body      string    `json:"body"`
	User      struct {
		ID        int       `json:"id"`
		Username  string    `json:"username"`
		Email     string    `json:"email"`
		Provider  string    `json:"provider"`
		Confirmed bool      `json:"confirmed"`
		Blocked   bool      `json:"blocked"`
		Role      int       `json:"role"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	} `json:"user"`
	CreatedAt time.Time   `json:"created_at"`
	UpdatedAt time.Time   `json:"updated_at"`
	LinkText  string      `json:"link_text"`
	LinkURL   interface{} `json:"link_url"`
	Thumbnail struct {
		ID               int         `json:"id"`
		Name             string      `json:"name"`
		Hash             string      `json:"hash"`
		Sha256           string      `json:"sha256"`
		Ext              string      `json:"ext"`
		Mime             string      `json:"mime"`
		Size             string      `json:"size"`
		URL              string      `json:"url"`
		Provider         string      `json:"provider"`
		ProviderMetadata interface{} `json:"provider_metadata"`
		CreatedAt        time.Time   `json:"created_at"`
		UpdatedAt        time.Time   `json:"updated_at"`
	} `json:"thumbnail"`
	Gallery []interface{} `json:"gallery"`
}

func cmsURL(path string) string {
	h := viper.GetString("cms_host_internal")
	if h == "" {
		h = "https://cms.verimap.com"
	}
	return fmt.Sprintf("%s%s", h, path)
}
func apiGetFireWatches(ctx context.Context, client *http.Client, count int) ([]FireWatch, error) {
	var resp []FireWatch
	u := fmt.Sprintf(cmsURL("/firewatches?_sort=published:desc&_limit=%d"), count)
	if err := get(ctx, client, u, &resp); err != nil {
		return nil, errors.Wrapf(err, "Failed to make get request")
	}
	return resp, nil
}

func apiGetShowcases(ctx context.Context, client *http.Client) ([]Showcase, error) {
	var resp []Showcase
	if err := get(ctx, client, cmsURL("/showcases?_sort=order"), &resp); err != nil {
		return nil, errors.Wrapf(err, "Failed to make get request")
	}
	return resp, nil
}

func apiGetExample(ctx context.Context, client *http.Client, ID int) (examplePage, error) {
	var resp []examplePage
	url := fmt.Sprintf(cmsURL("/examples?public=true&id=%d"), ID)
	if err := get(ctx, client, url, &resp); err != nil {
		return examplePage{}, errors.Wrapf(err, "Failed to make get request")
	}
	for i, page := range resp {
		m := make(map[string]string)
		for _, row := range strings.Split(page.Stats, "\n") {
			cols := strings.SplitN(row, "|", 2)
			if len(cols) == 2 {
				m[cols[0]] = cols[1]
			}
		}
		page.StatsMap = m
		if i == 0 {
			return page, nil
		}
	}
	return examplePage{}, errors.New("Unknown result")
}

func apiGetExamples(ctx context.Context, client *http.Client) ([]examplePage, error) {
	var resp []examplePage
	if err := get(ctx, client, cmsURL("/examples?public=true"), &resp); err != nil {
		return nil, errors.Wrapf(err, "Failed to make get request")
	}
	for i, page := range resp {
		m := make(map[string]string)
		for _, row := range strings.Split(page.Stats, "\n") {
			cols := strings.SplitN(row, "|", 2)
			if len(cols) == 2 {
				m[cols[0]] = cols[1]
			}
		}
		resp[i].StatsMap = m
	}
	return resp, nil
}

func get(ctx context.Context, client *http.Client, url string, recv interface{}) error {
	c, cancel := context.WithTimeout(ctx, time.Second*30)
	defer cancel()
	req, err := http.NewRequestWithContext(c, "GET", url, nil)
	if err != nil {
		return errors.Wrapf(err, "Failed to create request")
	}
	resp, err := client.Do(req)
	if err != nil {
		return errors.Wrapf(err, "Failed to perform request")
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrapf(err, "Failed to read response body")
	}
	if err := json.Unmarshal(b, recv); err != nil {
		return errors.Wrapf(err, "Failed to decode json response")
	}
	return nil
}
