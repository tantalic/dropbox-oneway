package dropbox // import "tantalic.com/dropbox"

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
)

const (
	rpcBaseURL     = "https://api.dropboxapi.com/2/"
	contentBaseURL = "https://content.dropboxapi.com/2/"
)

type APIError struct {
	Summary     string                 `json:"error_summary"`
	UserMessage string                 `json:"user_message"`
	Details     map[string]interface{} `json:"error"`
}

func (e APIError) Error() string {
	return e.Summary
}

type Client struct {
	AuthorizationToken string
	BaseURL            string
	HttpClient         http.Client
}

func (c *Client) rpc(method string, args interface{}, result interface{}) error {
	URL := rpcBaseURL + method

	body, err := json.Marshal(args)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", URL, bytes.NewReader(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	if c.AuthorizationToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthorizationToken)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		var apiError APIError
		json.Unmarshal(bytes, &apiError)
		return apiError
	}

	err = json.Unmarshal(bytes, result)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) content(method string, args interface{}) (io.ReadCloser, error) {

	URL := contentBaseURL + method
	req, err := http.NewRequest("POST", URL, nil)
	if err != nil {
		return nil, err
	}

	arg, err := json.Marshal(args)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Dropbox-API-Arg", string(arg))

	if c.AuthorizationToken != "" {
		req.Header.Set("Authorization", "Bearer "+c.AuthorizationToken)
	}

	resp, err := c.HttpClient.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		defer resp.Body.Close()
		bytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		var apiError APIError
		json.Unmarshal(bytes, &apiError)
		return nil, apiError
	}

	return resp.Body, nil
}
