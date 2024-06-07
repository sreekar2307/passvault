package services

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"passVault/dtos"
	"passVault/interfaces"
	"passVault/resources"
	"strings"
	"time"
)

type captchaServiceImpl struct {
	client   *http.Client
	endpoint string
}

const (
	captchaEndpoint = "https://www.google.com/recaptcha/api/siteverify"
)

func NewCaptchaService(client *http.Client) interfaces.CaptchaService {
	return captchaServiceImpl{
		client:   client,
		endpoint: captchaEndpoint,
	}
}

func (c captchaServiceImpl) VerifyToken(ctx context.Context, token string) (bool, error) {
	var (
		secret     = resources.Config().GetString(dtos.ConfigKeys.Captcha.Secret)
		cancelFunc context.CancelFunc

		respBody struct {
			Success     bool     `json:"success"`
			ErrorCodes  []string `json:"error-codes"`
			HostName    string   `json:"hostname"`
			ChallengeTs string   `json:"challenge_ts"`
		}
	)

	ctx, cancelFunc = context.WithTimeout(ctx, 5*time.Second)
	defer cancelFunc()

	reqBody := strings.NewReader(fmt.Sprintf("secret=%s&response=%s", secret, token))

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.endpoint,
		reqBody,
	)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	resp, err := c.client.Do(req)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return false, err
	}

	return respBody.Success, nil
}
