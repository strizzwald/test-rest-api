package rutter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/knadh/koanf"
	"go.uber.org/zap"
)

type RutterService struct {
	clientId     string
	clientSecret string
	apiUrl       string
	logger       *zap.SugaredLogger
}

type exchangeTokenRequest struct {
	ClientId    string `json:"client_id"`
	Secret      string `json:"secret"`
	PublicToken string `json:"public_token"`
}

type rutterError struct {
	ErrorType    string `json:"error_type"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
}

func New(config *koanf.Koanf) (*RutterService, error) {
	clientId := config.String("rutter.client_id")
	clientSecret := config.String("rutter.client_secret")
	apiUrl := config.String("rutter.api_url")

	if len(clientId) == 0 {
		return nil, fmt.Errorf("client_id is required")
	}

	if len(clientSecret) == 0 {
		return nil, fmt.Errorf("client_secret is required")
	}

	if len(apiUrl) == 0 {
		return nil, fmt.Errorf("api_url is required")
	}

	return &RutterService{clientId, clientSecret, apiUrl, zap.S().With("package", "rutter")}, nil
}

func (r *RutterService) GetAccesToken(ctx context.Context, authCode string) error {
	request := exchangeTokenRequest{
		ClientId:    r.clientId,
		Secret:      r.clientSecret,
		PublicToken: authCode,
	}

	serialized, _ := json.Marshal(request)

	body := bytes.NewBuffer(serialized)

	req, _ := http.NewRequestWithContext(ctx, "POST", r.apiUrl, body)
	client := http.DefaultClient

	resp, err := client.Do(req)

	if err != nil {
		var re rutterError
		body, _ := ioutil.ReadAll(resp.Body)

		json.Unmarshal(body, &re)
		r.logger.Error("Failed to fetch access token. %s", re.ErrorMessage)
		return fmt.Errorf("Failed to fetch access token. %s", re.ErrorMessage)
	}

	defer resp.Body.Close()

	return nil
}
