package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"github.com/joho/godotenv"
)

func main(){
	godotenv.Load(".env")
	fmt.Println("Hello world!")
	portString := os.Getenv("PORT")

	if portString == ""{
		log.Fatal("PORT is not found in the environment")
	}
 fmt.Println("port: ",portString)

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

 srv:= &http.Server{
	Handler: router,
	Addr: ":"+portString,
 }
 fmt.Printf("Server listening on port %v\n",portString)
 err := srv.ListenAndServe()
 if err != nil{
	log.Fatal(err)
 }
}