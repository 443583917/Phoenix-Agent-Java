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

type WeComClient struct {
	corpID     string
	corpSecret string
	http       *http.Client
}

func NewWeComClient(corpID, corpSecret string) *WeComClient {
	return &WeComClient{
		corpID:     corpID,
		corpSecret: corpSecret,
		http:       &http.Client{Timeout: 10 * time.Second},
	}
}

func (c *WeComClient) GetAccessToken(ctx context.Context) (string, error) {
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s", c.corpID, c.corpSecret)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return "", err
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("wecom gettoken: %s", result.ErrMsg)
	}
	return result.AccessToken, nil
}

func (c *WeComClient) GetUserIDByMobile(ctx context.Context, mobile string) (string, error) {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return "", err
	}

	body, _ := json.Marshal(map[string]string{"mobile": mobile})
	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/user/getbymobile?access_token=%s", token)
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
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
		UserID  string `json:"userid"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if result.ErrCode != 0 {
		return "", fmt.Errorf("wecom getbymobile: %s", result.ErrMsg)
	}
	return result.UserID, nil
}

func (c *WeComClient) SendMessage(ctx context.Context, userID, msg string) error {
	token, err := c.GetAccessToken(ctx)
	if err != nil {
		return err
	}

	body, _ := json.Marshal(map[string]any{
		"touser": userID,
		"msgtype": "text",
		"agentid": 0,
		"text": map[string]string{
			"content": msg,
		},
	})

	url := fmt.Sprintf("https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s", token)
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
		return fmt.Errorf("wecom send: %s", result.ErrMsg)
	}
	return nil
}
