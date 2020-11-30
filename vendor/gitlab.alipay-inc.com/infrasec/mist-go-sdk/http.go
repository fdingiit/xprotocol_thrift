package mist

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var defaultTimeout = time.Duration(5) * time.Second

type mistIssueResponse struct {
	Code int `json:"code"`
	Data struct {
		JWTSvid string `json:"jwt-svid"`
	} `json:"data"`
	Err string `json:"err"`
}

type mistVerifyResponse struct {
	Code int `json:"code"`
	Data struct {
		Valid bool `json:"valid"`
	} `json:"data"`
	Err string `json:"err"`
}

func IssueJWTSVID(url string, exp int64) (jwtSvid string, err error) {
	if url == "" {
		return "", fmt.Errorf("url is empty")
	}

	req, err := http.NewRequest(
		"GET", url, nil,
	)
	if err != nil {
		return "", err
	}
	req.Header.Set("content-type", "application/json")

	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Parse the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	mistResp := &mistIssueResponse{}
	if err := json.Unmarshal(respBody, mistResp); err != nil {
		return "", err
	}
	if mistResp.Code != 0 {
		return "", fmt.Errorf(mistResp.Err)
	}
	return mistResp.Data.JWTSvid, nil
}

func VerifyJWTSVID(url string, svid string) (ok bool, err error) {
	if url == "" {
		return false, fmt.Errorf("url is empty")
	}

	req, err := http.NewRequest(
		"POST", url,
		strings.NewReader(fmt.Sprintf(`{"jwt-svid":"%s"}`, svid)),
	)
	if err != nil {
		return false, err
	}
	req.Header.Set("content-type", "application/json")

	client := &http.Client{Timeout: defaultTimeout}
	resp, err := client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	// Parse the response body
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return false, err
	}

	mistResp := &mistVerifyResponse{}
	if err := json.Unmarshal(respBody, mistResp); err != nil {
		return false, err
	}
	if mistResp.Code != 0 {
		return false, fmt.Errorf(mistResp.Err)
	}
	return mistResp.Data.Valid, nil
}
