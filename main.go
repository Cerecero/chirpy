package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cerecero/chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)
type apiConfig struct{
	fileserverHits atomic.Int32
	dbQueries *database.Queries
	platform string
}
type Request struct {
	Body string `json:"body"`
}

type Response struct {
	Error string `json:"error,omitempty"`
	CleanedBody string `json:"cleaned_body,omitempty"`
}
type User struct {
	ID uuid.UUID 		`json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email string 		`json:"email"`
}
var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(Response{Error: msg})
	if err != nil{
		panic(err)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(payload)
	if err != nil {
		panic(err)
	}
}
func replaceProfanity(input string) string {
	words := strings.Split(input, " ")
	for i, word := range words {
		cleaned := strings.ToLower(word)
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

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Access forbidden")
		return
	}
 	err = cfg.dbQueries.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reset database")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "All users deleted"})

}

func (cfg *apiConfig) handleValidateReq (w http.ResponseWriter, r *http.Request){
	if r.Method != http.MethodPost{
		respondWithError(w, http.StatusMethodNotAllowed, "Method not allowed")
		return
	}

	req := Request{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "Invalid JSON")
		return
	}

	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}	
	cleaned := replaceProfanity(req.Body)

	respondWithJSON(w, http.StatusOK, Response{CleanedBody: cleaned})
}
func (cfg *apiConfig) handleCreateUser (w http.ResponseWriter, r *http.Request){

	// req := Request{}
	type requestBody struct {
		Email string `json:"email"`
	}
	var req requestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
	}
	user, err := cfg.dbQueries.CreateUser(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error")
	}

	usr := User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}
	
	respondWithJSON(w, http.StatusCreated, usr)

}

func main(){
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	dbURL := os.Getenv("DB_URL")
	platform := os.Getenv("PLATFORM")

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		panic(err)
	}
	dbQueries := database.New(db)

	apiCfg := &apiConfig{
		dbQueries: dbQueries,
		platform: platform,
	}

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

	mux.HandleFunc("POST /api/validate_chirp", apiCfg.handleValidateReq)

	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)

	server := &http.Server{
		Handler: mux,
		Addr: ":8080",
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
