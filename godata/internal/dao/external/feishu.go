package external

import "context"

// FeishuClient is a stub for the Feishu (Lark) SDK.
// This provides a placeholder for future Feishu integration.
type FeishuClient struct {
	appID     string
	appSecret string
}

func NewFeishuClient(appID, appSecret string) *FeishuClient {
	return &FeishuClient{
		appID:     appID,
		appSecret: appSecret,
	}
}

// GetUserIDByMobile is a stub for looking up a Feishu user by mobile number.
func (c *FeishuClient) GetUserIDByMobile(ctx context.Context, mobile string) (string, error) {
	return "", nil
}

// SendMessage is a stub for sending Feishu messages.
func (c *FeishuClient) SendMessage(ctx context.Context, userID, msg string) error {
	return nil
}
