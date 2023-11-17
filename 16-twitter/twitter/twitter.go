package twitter

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	DefaultBaseURL    = "https://api.twitter.com"
	DefaultAPIVersion = "1.1"
)

type Client struct {
	httpClient http.Client

	consumerKey    string
	consumerSecret string
	accessToken    string

	BaseURL    string
	APIVersion string
}

func NewClient(consumerKey, consumerSecret string) *Client {
	return &Client{
		consumerKey:    consumerKey,
		consumerSecret: consumerSecret,
		BaseURL:        DefaultBaseURL,
		APIVersion:     DefaultAPIVersion,
	}
}

func (c *Client) authenticate() error {
	endpointURL := fmt.Sprintf("%s/oauth2/token", c.BaseURL)

	reqBody := url.Values{}
	reqBody.Add("grant_type", "client_credentials")
	req, err := http.NewRequest(
		http.MethodPost,
		endpointURL,
		strings.NewReader(reqBody.Encode()))
	if err != nil {
		return err
	}
	req.SetBasicAuth(c.consumerKey, c.consumerSecret)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if code := res.StatusCode; code != http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("got: %d %s", code, b)
	}

	var resBody struct {
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return err
	}

	c.accessToken = resBody.AccessToken

	return nil
}

type Status struct {
	User struct {
		ScreenName string `json:"screen_name"`
	} `json:"user"`
}

func (c *Client) StatusRetweets(id string) ([]Status, error) {
	if c.accessToken == "" {
		if err := c.authenticate(); err != nil {
			return nil, err
		}
	}

	endpointURL := fmt.Sprintf("%s/%s/statuses/retweets/%s.json", c.BaseURL, c.APIVersion, id)
	// endpointURL := path.Join(c.BaseURL, "statuses/retweets", fmt.Sprintf("%s.json", id))

	req, err := http.NewRequest(http.MethodGet, endpointURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Authorization", fmt.Sprint("Bearer ", c.accessToken))

	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if code := res.StatusCode; code != http.StatusOK {
		b, err := ioutil.ReadAll(res.Body)
		if err != nil {
			return nil, err
		}
		return nil, fmt.Errorf("got: %d %s", code, b)
	}

	var resBody []Status
	if err := json.NewDecoder(res.Body).Decode(&resBody); err != nil {
		return nil, err
	}
	return resBody, nil
}
