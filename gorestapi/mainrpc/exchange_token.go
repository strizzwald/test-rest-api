package mainrpc

import (
	"context"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/snowzach/gorestapi/rutter"
	"github.com/snowzach/gorestapi/server"
)

// Exchanges authorization codes with access tokens from the Rutter API.
// @Tags Exchange Token
// @Summary Exchanges authorization code for access token.
// @Description Exchanges authorization code for access token.
// @Param sellerId path string true "Seller Id"
// @Param authCode query string true "Authorization Code"
// @Success 200
// @Failure 400 {object} server.ErrResponse "Invalid Argument"
// @Failure 500 {object} server.ErrResponse "Internal Server Error"
func (s *Server) ExchangeTokenSave() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authCode := chi.URLParam(r, "authCode")
		service, err := rutter.New()

		if err != nil {
			server.RenderErrInternal(w, err)
		}

		accessToken, err := service.GetAccesToken(context.Background(), authCode)

		if err != nil {
			server.RenderErrInvalidRequest(w, err)
		}

	}
}
