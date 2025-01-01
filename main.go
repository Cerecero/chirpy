package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync/atomic"
)

type apiConfig struct{
	fileserverHits atomic.Int32
}
type Request struct {
	Body string `json:"body"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	CleanedBody string `json:"cleaned_body,omitempty"`
}
var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(Response{Error: msg})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}
func replaceProfanity(input string) string {
	words := strings.Split(input, " ")
	for i, word := range words {
		cleaned := strings.ToLower(word) // Convert word to lowercase for case-insensitive matching
		for _, profane := range profaneWords {
			if cleaned == profane {
				words[i] = "****"
				break
			}
		}
	}
	return strings.Join(words, " ")
}

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler{
	return http.HandlerFunc(func( w http.ResponseWriter, r *http.Request){
		cfg.fileserverHits.Add(1)
		next.ServeHTTP(w, r)
	})
}

func (cfg *apiConfig) handleMetrics(w http.ResponseWriter, r *http.Request) {
	hits := cfg.fileserverHits.Load()
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, err := w.Write([]byte(fmt.Sprintf(`
	<html>
		<body>
			<h1>Welcome, Chirpy Admin</h1>
			<p>Chirpy has been visited %d times!</p>
		</body>
	</html>`, hits)))
	if err != nil {
		panic(err)
	}
}


func (cfg *apiConfig) handleReset(w http.ResponseWriter, r *http.Request) {
	cfg.fileserverHits.Store(0)
	w.Header().Set("Content-type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_,err := w.Write([]byte("Hits back to 0"))
	if err != nil {
		panic(err)
	}

}

func (cfg *apiConfig) handleValidateRe (w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "Invalid Json"})
		return
	}

	if len(req.Body) > 140 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(Response{Error: "Chirp is too long"})
		return
	}	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(Response{Valid: true})
}

func main(){

	apiCfg := &apiConfig{}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", func( w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		
		_,err := w.Write([]byte("OK"))
		if err != nil {
			panic(err)
		}
	})
	fileserver := http.FileServer(http.Dir("."))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileserver)))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)

	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handleValidateRe)
	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
