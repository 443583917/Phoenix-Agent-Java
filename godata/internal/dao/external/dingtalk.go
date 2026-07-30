package external

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

type DingTalkClient struct {
	appKey    string
	appSecret string
	http      *http.Client
}

func NewDingTalkClient(appKey, appSecret string) *DingTalkClient {
	return &DingTalkClient{
		appKey:    appKey,
		appSecret: appSecret,
		http:      &http.Client{Timeout: 10 * time.Second},
	}
}

type dingtalkTokenResp struct {
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
	AccessToken string `json:"access_token"`
}

func (c *DingTalkClient) getAccessToken(ctx context.Context) (string, error) {
	url := fmt.Sprintf("https://oapi.dingtalk.com/gettoken?appkey=%s&appsecret=%s", c.appKey, c.appSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result dingtalkTokenResp
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("dingtalk gettoken: %s", result.ErrMsg)
	}
	return result.AccessToken, nil
}

func (c *DingTalkClient) GetUserIDByMobile(ctx context.Context, mobile string) (string, error) {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return "", err
	}

	body, _ := json.Marshal(map[string]string{"mobile": mobile})
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/v2/user/getbymobile?access_token=%s", token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		ErrCode int `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		Result  struct {
			UserID string `json:"userid"`
		} `json:"result"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("dingtalk getbymobile: %s", result.ErrMsg)
	}
	return result.Result.UserID, nil
}

func (c *DingTalkClient) SendMessage(ctx context.Context, userID, msg string) error {
	token, err := c.getAccessToken(ctx)
	if err != nil {
		return err
	}

	body, _ := json.Marshal(map[string]any{
		"agent_id": "",
		"userid_list": userID,
		"msg": map[string]any{
			"msgtype": "text",
			"text":    map[string]string{"content": msg},
		},
	})
	url := fmt.Sprintf("https://oapi.dingtalk.com/topapi/message/corpconversation/asyncsend_v2?access_token=%s", token)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.ErrCode != 0 {
		return fmt.Errorf("dingtalk send: %s", result.ErrMsg)
	}
	return nil
}
