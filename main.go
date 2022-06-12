package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
)

func main (){
	port := "3000"

	if fromEnv := os.Getenv("PORT"); fromEnv != ""{
		port = fromEnv
	}

	log.Printf("Starting up on https://localhost:%v", port)

	r:= chi.NewRouter()

	r.Get("/", func(w http.ResponseWriter, r*http.Request){
		w.Header().Set("Content-Type", "text-plain")
		w.Write([]byte("Hello World!"))
	})

	r.Mount("/movies", moviesResource{}.Routes())

	log.Fatal(http.ListenAndServe(":"+ port, r))
}