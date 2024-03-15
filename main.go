package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/PaccyC/rss-aggregator/internal/database"
	"github.com/joho/godotenv"

	_ "github.com/lib/pq"
)

type apiConfig struct {
	DB *database.Queries
}


func main(){
 

 feed,err :=	urlToFeed("https://wagslane.dev/index.xml")
 if err!= nil {
    log.Fatal(err)
 }
   fmt.Println(feed)


	godotenv.Load(".env")
	fmt.Println("Hello world!")
	portString := os.Getenv("PORT")

	if portString == ""{
		log.Fatal("PORT is not found in the environment")
	}

	dbUrl := os.Getenv("DB_URL")
	log.Println(dbUrl)
	if dbUrl == "" {
        log.Fatal("DB_URL is not found in the environment")
    }

  conn, err :=sql.Open("postgres", dbUrl)
  if err != nil {
	log.Fatal("Can't connect to database",err)
  }

    db := database.New(conn)
   apiCfg := apiConfig{
	DB:db,
   }

   go startScraping(db,10,time.Minute)


 router:= chi.NewRouter()
 router.Use(cors.Handler(cors.Options{
	AllowedOrigins: []string{"https://*","http://*"},
	AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	AllowedHeaders: []string{"*"},
	ExposedHeaders: []string{"Link"},
	AllowCredentials: false,
	MaxAge: 300,
 }))

 v1Router := chi.NewRouter()

 v1Router.Get("/healthz",handlerReadiness)
 v1Router.Get("/err",handlerError)
 router.Mount("/v1",v1Router)

 v1Router.Post("/users",apiCfg.handlerCreateUser)
 v1Router.Get("/users",apiCfg.middlewareAuth(apiCfg.handlerGetUser))
 v1Router.Post("/feeds",apiCfg.middlewareAuth(apiCfg.handlerCreateFeed))
 v1Router.Get("/feeds",apiCfg.handlerGetFeeds)
v1Router.Get("/posts",apiCfg.middlewareAuth(apiCfg.handlerGetPostsForUser))

 v1Router.Post("/feed-follows",apiCfg.middlewareAuth(apiCfg.handlerCreateFeedFollow))
 v1Router.Get("/feed-follows",apiCfg.middlewareAuth(apiCfg.handlerGetFeedFollows))
 v1Router.Delete("/feed-follows/{feedFollowID}",apiCfg.middlewareAuth(apiCfg.handlerDeleteFeedFollows))

 srv:= &http.Server{
	Handler: router,
	Addr: ":"+portString,
 }
 fmt.Printf("Server listening on port %v\n",portString)
 srv.ListenAndServe()
 if err != nil{
	log.Fatal(err)
 }
}