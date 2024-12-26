package main

import (
	"net/http"
	"sync/atomic"
)

type apiConfig struct{
	fileserverHits atomic.Int32
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func main(){

	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func( w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		
		_,err := w.Write([]byte("OK"))
		if err != nil {
			panic(err)
		}
	})
	fileserver := http.fileserver(http.dir("."))
	mux.handle("/app/", http.stripprefix("/app", fileserver))


	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
