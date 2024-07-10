package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/shubhrad1/rssagg/internal/database"
)

func (apiCfg *apiConfig) createFeedFollowHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	type parameters struct {
		FeedId uuid.UUID `json:"feed_id"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}
	feedFollow, err := apiCfg.DB.CreateFeedFollow(r.Context(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdateAt:  time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    params.FeedId,
	})
	if err != nil {
		respondError(w, 400, fmt.Sprintf("Couldn't create feed follow: %s", err))
		return
	}

	respondJSON(w, 201, databaseFeedFollowtoFeedFollow(feedFollow))

}

func (apiCfg *apiConfig) getFeedFollowsHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollows, err := apiCfg.DB.GetFeedFollows(r.Context(), user.ID)
	if err != nil {
		respondError(w, 400, fmt.Sprintf("Couldn't get feed follows: %s", err))
		return
	}

	respondJSON(w, 201, databaseFeedFollowstoFeedFollowsAll(feedFollows))

}

func (apiCfg *apiConfig) deleteFeedFollowsHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	feedFollowID_str := chi.URLParam(r, "feedFollowID")
	feedFollowID, err := uuid.Parse(feedFollowID_str)
	if err != nil {
		respondError(w, 400, fmt.Sprintf("[ERROR] Couldn't parse feed follow ID: %v", err))
		return
	}
	err = apiCfg.DB.DeleteFeedFollow(r.Context(), database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: user.ID,
	})
	if err != nil {
		respondError(w, 400, fmt.Sprintf("[ERROR] Couldn't delete feed follow: %v", err))
		return
	}
	respondJSON(w, 200, struct{}{})
}
