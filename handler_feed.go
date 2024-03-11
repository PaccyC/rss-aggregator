package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
	"github.com/PaccyC/rss-aggregator/internal/database"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig)handlerCreateFeed(w http.ResponseWriter, r *http.Request,user database.User){
	
	type parameters struct{ 
		Name string `json:"name"`
		URL string `json:"url"`

	}

	decoder :=json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Error Partsing JSOn %v",err))
		return
	}
 feed, err :=apiCfg.DB.CreateFeed(r.Context(),database.CreateFeedParams{
		ID:uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name: params.Name,
		Url: params.URL,
		UserID: user.ID,
	})

	if err!= nil{
        respondWithError(w,400,fmt.Sprintf("Error Creating User: %s",err))
        return
    }
	respondWithJSON(w,201,databaseFeedToFeed(feed))
}


func (apiCfg *apiConfig) handlerGetFeeds(w http.ResponseWriter, r *http.Request) {
    // Assuming apiCfg.DB.GetFeeds returns []database.Feed and an error
    feeds, err := apiCfg.DB.GetFeeds(r.Context())
    if err != nil {
        respondWithError(w, 400, fmt.Sprintf("Couldn't get feeds: %s", err))
        return
    }
    respondWithJSON(w, 200, feeds)
}
