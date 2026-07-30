package external

import "context"

// DingTalkClient is a stub for the DingTalk SDK.
// This provides a placeholder for future DingTalk integration.
type DingTalkClient struct {
	appKey    string
	appSecret string
}

func NewDingTalkClient(appKey, appSecret string) *DingTalkClient {
	return &DingTalkClient{
		appKey:    appKey,
		appSecret: appSecret,
	}
}

// GetUserIDByMobile is a stub for looking up a DingTalk user by mobile number.
func (c *DingTalkClient) GetUserIDByMobile(ctx context.Context, mobile string) (string, error) {
	return "", nil
}

// SendMessage is a stub for sending DingTalk messages.
func (c *DingTalkClient) SendMessage(ctx context.Context, userID, msg string) error {
	return nil
}
