package main

import (
	"fmt"
	"net/http"

	"github.com/shubhrad1/rssagg/internal/auth"
	"github.com/shubhrad1/rssagg/internal/database"
)

type authedHandler func(http.ResponseWriter, *http.Request, database.User)

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondError(w, 403, fmt.Sprintf("Auth error: %s", err))
			return
		}

		user, err := apiCfg.DB.GetUserByApiKey(r.Context(), apiKey)
		if err != nil {
			respondError(w, 400, fmt.Sprintf("Counln't get user: %v", err))
			return
		}
		handler(w, r, user)
	}

}
