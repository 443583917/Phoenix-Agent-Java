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

type FeishuClient struct {
	appID     string
	appSecret string
	http      *http.Client
}

func NewFeishuClient(appID, appSecret string) *FeishuClient {
	return &FeishuClient{
		appID:     appID,
		appSecret: appSecret,
		http:      &http.Client{Timeout: 10 * time.Second},
	}
}

type FeishuToken struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

type FeishuUser struct {
	UserID string `json:"user_id"`
	Name   string `json:"name"`
	Mobile string `json:"mobile"`
	Email  string `json:"email"`
}

func (c *FeishuClient) getTenantAccessToken(ctx context.Context) (string, error) {
	body, _ := json.Marshal(map[string]string{
		"app_id":     c.appID,
		"app_secret": c.appSecret,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result struct {
		Code           int    `json:"code"`
		Msg            string `json:"msg"`
		TenantAccessToken string `json:"tenant_access_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	if result.Code != 0 {
		return "", fmt.Errorf("feishu tenant token: %s", result.Msg)
	}
	return result.TenantAccessToken, nil
}

func (c *FeishuClient) GetUserIDByMobile(ctx context.Context, mobile string) (string, error) {
	token, err := c.getTenantAccessToken(ctx)
	if err != nil {
		return "", err
	}

	body, _ := json.Marshal(map[string]string{"mobile": mobile})
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://open.feishu.cn/open-apis/contact/v3/users/batch_get_id?user_id_type=open_id", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, _ := io.ReadAll(resp.Body)
	var result struct {
		Code int `json:"code"`
		Msg  string `json:"msg"`
		Data struct {
			UserList []struct {
				UserID string `json:"user_id"`
			} `json:"user_list"`
		} `json:"data"`
	}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", err
	}
	if result.Code != 0 {
		return "", fmt.Errorf("feishu getbymobile: %s", result.Msg)
	}
	if len(result.Data.UserList) == 0 {
		return "", fmt.Errorf("feishu: user not found for mobile")
	}
	return result.Data.UserList[0].UserID, nil
}

func (c *FeishuClient) SendMessage(ctx context.Context, userID, msg string) error {
	token, err := c.getTenantAccessToken(ctx)
	if err != nil {
		return err
	}

	body, _ := json.Marshal(map[string]any{
		"receive_id": userID,
		"msg_type":   "text",
		"content":    `{"text":"` + msg + `"}`,
	})

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://open.feishu.cn/open-apis/im/v1/messages?receive_id_type=open_id", bytes.NewReader(body))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := c.http.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var result struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}
	if result.Code != 0 {
		return fmt.Errorf("feishu send: %s", result.Msg)
	}
	return nil
}
