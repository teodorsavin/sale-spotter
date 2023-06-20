package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"teodorsavin/ah-bonus/util"
)

type APIClient struct {
	BaseURL    string
	HTTPClient *http.Client
}

func NewAPIClient(baseURL string, timeout time.Duration) *APIClient {
	return &APIClient{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *APIClient) Login() string {
	loginData := util.LoginData{
		ClientId: "appie",
	}
	marshalled, err := json.Marshal(loginData)
	if err != nil {
		log.Fatalf("impossible to marshal loginData: %s", err)
	}

	req, err := http.NewRequest("POST", c.BaseURL+"/mobile-auth/v1/auth/token/anonymous", bytes.NewReader(marshalled))
	if err != nil {
		log.Fatalf("impossible to build request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")

	res := c.DoRequest(req)
	accessToken := ReadResponseBodyLogin(res)
	return accessToken
}

func (c *APIClient) DoRequest(request *http.Request) *http.Response {
	res, err := c.HTTPClient.Do(request)
	if err != nil {
		log.Fatalf("impossible to send request: %s", err)
	}
	log.Printf("status Code: %d", res.StatusCode)

	return res
}
