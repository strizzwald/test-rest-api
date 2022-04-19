package postgres

import (
	"context"
	"fmt"

	"github.com/snowzach/gorestapi/rutter"
)

func (c *Client) SaveRutterAccessToken(ctx context.Context, accesstoken *rutter.RutterAccessToken, sellerId string) error {
	fmt.Println("TODO: Save Access Token")
	return nil
}
