package external

import "context"

// WeComClient is a stub for the WeCom (WeChat Work) SDK.
// This provides a placeholder for future WeCom integration.
type WeComClient struct {
	corpID     string
	corpSecret string
}

func NewWeComClient(corpID, corpSecret string) *WeComClient {
	return &WeComClient{
		corpID:     corpID,
		corpSecret: corpSecret,
	}
}

// GetUserIDByMobile is a stub for looking up a WeCom user by mobile number.
func (c *WeComClient) GetUserIDByMobile(ctx context.Context, mobile string) (string, error) {
	return "", nil
}

// SendMessage is a stub for sending WeCom messages.
func (c *WeComClient) SendMessage(ctx context.Context, userID, msg string) error {
	return nil
}
