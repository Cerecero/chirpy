package main

import "net/http"

type Server struct{

}

func main(){

	mux := http.NewServeMux()
	server := http.Server{
		Handler: mux,
		Addr: ":8080",
	}

	server.ListenAndServe()

	
	// mux.Handle("/", http.FileServer(http.Dir(".")))

}
