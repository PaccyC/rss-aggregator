package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/PaccyC/rss-aggregator/internal/database"
	"github.com/go-chi/chi"
	"github.com/google/uuid"
)

func (apiCfg *apiConfig)handlerCreateFeedFollow(w http.ResponseWriter, r *http.Request,user database.User){
	
	type parameters struct{ 
	FeedID uuid.UUID   `json:"feed_id"`

	}

	decoder :=json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil{
		respondWithError(w,400,fmt.Sprintf("Error Partsing JSOn %v",err))
		return
	}
 feedFollow, err :=apiCfg.DB.CreateFeedFollows(r.Context(),database.CreateFeedFollowsParams{
		ID:uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),

		UserID: user.ID,
		FeedID: params.FeedID,
	})

	if err!= nil{
        respondWithError(w,400,fmt.Sprintf("Error Creating feed follows: %s",err))
        return
    }
	respondWithJSON(w,201,databaseFeedFollowToFeedFollow(feedFollow))
}
// Getting feeds


func (apiCfg *apiConfig)handlerGetFeedFollows(w http.ResponseWriter, r *http.Request,user database.User){
	

 feedFollows, err :=apiCfg.DB.GetFeedFollows(r.Context(),user.ID)

	if err!= nil{
        respondWithError(w,400,fmt.Sprintf("Error Creating feed follows: %s",err))
        return
    }
	respondWithJSON(w,201,databaseFeedFollowsToFeedFollows(feedFollows))
}

func (apiCfg *apiConfig)handlerDeleteFeedFollows(w http.ResponseWriter, r *http.Request,user database.User){
	
feedFollowStr := chi.URLParam(r,"feedFollowID")

feedFollowID, err := uuid.Parse(feedFollowStr)
if err != nil {
	respondWithError(w,400,"Failed to parse feedFollowID ")
}

 err = apiCfg.DB.DeleteFeedFollows(r.Context(),database.DeleteFeedFollowsParams{
	ID: feedFollowID,
	UserID: user.ID,
})
if err != nil {
	respondWithError(w,400, "Could not delete feed follow")
}

respondWithJSON(w,200,struct{}{})

   }
   