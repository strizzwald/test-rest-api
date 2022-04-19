package rutter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/snowzach/gorestapi/conf"
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

type RutterAccessToken struct {
	AccessToken     string `json:"access_token"`
	ConnectionId    string `json:"connection_id"`
	RequestId       string `json:"request_id"`
	IsReady         bool   `json:"is_ready"`
	StoreUniqueName string `json:"store_unique_name"`
	StoreUniqueId   string `json:"store_unique_id"`
	StoreDomain     string `json:"store_domain"`
	Platform        string `json:"platform"`
}

func New() (*RutterService, error) {

	clientId := conf.C.String("rutter.client_id")
	clientSecret := conf.C.String("rutter.client_secret")
	apiUrl := conf.C.String("rutter.api_url")

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

func (r *RutterService) GetAccesToken(ctx context.Context, authCode string) (*RutterAccessToken, error) {
	request := exchangeTokenRequest{
		ClientId:    r.clientId,
		Secret:      r.clientSecret,
		PublicToken: authCode,
	}

	serialized, _ := json.Marshal(request)

	body := bytes.NewBuffer(serialized)

	req, _ := http.NewRequestWithContext(ctx, "POST", r.apiUrl+"/item/public_token/exchange", body)
	req.Header.Add("Content-Type", "application/json")
	client := http.DefaultClient

	resp, err := client.Do(req)
	b, _ := ioutil.ReadAll(resp.Body)

	if resp.StatusCode != 200 {
		var re rutterError

		json.Unmarshal(b, &re)
		r.logger.Error("Failed to fetch access token. %s", re.ErrorMessage)
		return nil, fmt.Errorf("Failed to fetch access token. %s", re.ErrorMessage)
	}

	if err != nil {
		r.logger.Error("Failed to fetch access token. %s", err.Error())
		return nil, fmt.Errorf("Failed to fetch access token. %s", err.Error())
	}

	defer resp.Body.Close()

	var accessToken RutterAccessToken
	json.Unmarshal(b, &accessToken)

	return &accessToken, nil
}
