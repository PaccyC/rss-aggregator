package main

import (
	"net/http"
	"fmt"

	"github.com/PaccyC/rss-aggregator/internal/database"
	"github.com/PaccyC/rss-aggregator/internal/auth"
)

type authedHandler func (http.ResponseWriter,*http.Request,database.User) 

func (apiCfg *apiConfig) middlewareAuth(handler authedHandler)http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request) {
		apiKey,err := auth.GetApiKey(r.Header)
		if err !=nil{
			respondWithError(w,403,fmt.Sprintf("Auth Eror: %v",err))
			return
		}	
		user,err := apiCfg.DB.GetUserByAPIKey(r.Context(),apiKey)
		
		if err != nil{
			respondWithError(w,400,fmt.Sprintf("Counldn't get user: %v",err))
			return
		}
     handler(w,r,user)
	}
}