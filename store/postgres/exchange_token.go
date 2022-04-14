package postgres

import (
	"context"

	"github.com/snowzach/gorestapi/rutter"
)

const (
	RutterConnectionSchema = ``
	RutterConnectionTable  = `connection`
	RutterConnectionFields = `
		COALESCE(connection.id, '') as "connection.id",
		COALESCE(connection.created, '') as "connection.created",
		COALESCE(connection.connection_id, '') as "connection.connection_id",
		COALESCE(connection.access_token, '') as "connection.access_token",
		COALESCE(connection.seller_id, '') as "connection.seller_id",
		COALESCE(connection.platform_name, '') as "connection.platform_name",
		COALESCE(connection.store_domain, '') as "connection.store_domain",
		COALESCE(connection.store_name, '') as "connection.store_name",
		COALESCE(connection.store_unique_id, '') as "connection.store_unique_id"
	`
)

func (c *Client) SaveRutterAccessToken(ctx context.Context, accesstoken *rutter.RutterAccessToken, sellerId string) error {

	fields, values, updates, args := composeUpsert([]field{
		{name: "created", insert: "NOW()", update: ""},
		{name: "updated", insert: "", update: "NOW()"},
		{name: "connection_id", insert: "$#", update: "$#", arg: accesstoken.ConnectionId},
		{name: "access_token", insert: "$#", update: "$#", arg: accesstoken.AccessToken},
		{name: "seller_id", insert: "$#", update: "$#", arg: sellerId},
		{name: "platform_name", insert: "$#", update: "$#", arg: accesstoken.Platform},
		{name: "store_name", insert: "$#", update: "$#", arg: accesstoken.StoreUniqueName},
		{name: "store_domain", insert: "$#", update: "$#", arg: accesstoken.StoreDomain},
		{name: "store_unique_id", insert: "$#", update: "$#", arg: accesstoken.StoreUniqueId},
	})

	err := c.db.GetContext(ctx, accesstoken, `
	WITH `+Rutt+` AS (
        INSERT INTO `+ThingSchema+WidgetTable+` (`+fields+`)
        VALUES(`+values+`) ON CONFLICT (id) DO UPDATE
        SET `+updates+` RETURNING *
	) `+WidgetFields+" FROM "+WidgetTable+WidgetJoins, args...)`)

	return nil
}
