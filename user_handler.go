package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/shubhrad1/rssagg/internal/database"
)

func (apiCfg *apiConfig) createUserHandler(w http.ResponseWriter, r *http.Request) {

	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondError(w, 400, fmt.Sprintf("Error parsing JSON: %s", err))
		return
	}
	usr, err := apiCfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdateAt:  time.Now().UTC(),
		Name:      params.Name,
	})
	if err != nil {
		respondError(w, 400, fmt.Sprintf("Couldn't create user: %s", err))
		return
	}

	respondJSON(w, 201, databaseUserToUser(usr))

}

func (apiCfg *apiConfig) getUserHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	respondJSON(w, 200, databaseUserToUser(user))

}

func (apiCfg *apiConfig) getPostsHandler(w http.ResponseWriter, r *http.Request, user database.User) {

	posts, err := apiCfg.DB.GetPostsForUsers(r.Context(), database.GetPostsForUsersParams{
		UserID: user.ID,
		Limit:  10,
	})
	if err != nil {
		respondError(w, 400, fmt.Sprintf("Couldn't get posts: %v", err))
		return
	}
	respondJSON(w, 200, databasePostsToPosts(posts))

}
