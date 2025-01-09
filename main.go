package main

import (

    "github.com/cerecero/chirpy/internal"
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

type apiConfig struct {
	fileserverHits atomic.Int32
	dbQueries      *database.Queries
	platform       string
}
type Request struct {
	Body string `json:"body"`
}

type Response struct {
	Error       string `json:"error,omitempty"`
	CleanedBody string `json:"cleaned_body,omitempty"`
}
type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

var profaneWords = []string{"kerfuffle", "sharbert", "fornax"}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(Response{Error: msg})
	if err != nil {
		panic(err)
	}
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	enc := json.NewEncoder(w)
	enc.SetIndent("", "   ")
	if err := enc.Encode(payload); err != nil {
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

func (cfg *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	if cfg.platform != "dev" {
		respondWithError(w, http.StatusForbidden, "Access forbidden")
		return
	}
	err := cfg.dbQueries.DeleteUsers(r.Context())
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to reset database")
		return
	}

	respondWithJSON(w, http.StatusOK, map[string]string{"message": "All users deleted"})

}

func (cfg *apiConfig) handleChirp(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Body   string `json:"body"`
		UserID string `json:"user_id"`
	}
	var req requestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
	}
	if req.Body == "" || req.UserID == "" {
		respondWithError(w, http.StatusBadRequest, "Both body and user_id are required")
		return
	}
	if len(req.Body) > 140 {
		respondWithError(w, http.StatusBadRequest, "Chirp is too long")
		return
	}
	userID, err := uuid.Parse(req.UserID)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user_id fromat")
		return
	}
	cleanedBody := replaceProfanity(req.Body)
	chirpID := uuid.New()

	chirp, err := cfg.dbQueries.InsertChirp(r.Context(), database.InsertChirpParams{
		ID:     chirpID,
		Body:   cleanedBody,
		UserID: userID,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Failed to create chirp")
		return
	}

	respondWithJSON(w, http.StatusCreated, chirp)
}
func (cfg *apiConfig) handleCreateUser(w http.ResponseWriter, r *http.Request) {

	type requestBody struct {
		Email string `json:"email"`
		Password string `json:"password"`
	}
	var req requestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
	}

	hasshPass, err := auth.HashPassword(req.Password)
	if err != nil{
		respondWithError(w, http.StatusBadRequest, "error")
	}
	user, err := cfg.dbQueries.CreateUser(r.Context(), database.CreateUserParams{
		Email: req.Email,
		HashedPassword: sql.NullString{String: hasshPass, Valid: true},
	})
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "error")
	}

	usr := User{ID: user.ID, CreatedAt: user.CreatedAt, UpdatedAt: user.UpdatedAt, Email: user.Email}

	respondWithJSON(w, http.StatusCreated, usr)

}

func (cfg *apiConfig) handleQueryChirps(w http.ResponseWriter, r *http.Request) {

	queryChirps, err := cfg.dbQueries.QueryChirp(r.Context())
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Error")
	}
	type Chirp struct {
		ID        uuid.UUID
		CreatedAt time.Time
		UpdatedAt time.Time
		Body      string
		UserID    uuid.UUID
	}
	// chirp := Chirp{ID: queryChirps.ID, CreatedAt: queryChirps.CreatedAt, UpdatedAt: queryChirps.UpdatedAt, Body: queryChirps.Body, UserID: queryChirps.UserID}

	respondWithJSON(w, http.StatusOK, queryChirps)
}


func (cfg *apiConfig) handleLogin(w http.ResponseWriter, r *http.Request) {
	type requestBody struct {
		Password string `json:"password"`
		Email string `json:"email"`
	}

	var req requestBody
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request body")
	}
	usr, err := cfg.dbQueries.QueryUser(r.Context(), req.Email)
	if err != nil {
		respondWithError(w, http.StatusUnauthorized, "Invalid request")
	}
	usrHashPass

}
func main() {
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
		platform:  platform,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		_, err := w.Write([]byte("OK"))
		if err != nil {
			panic(err)
		}
	})
	fileserver := http.FileServer(http.Dir("."))
	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileserver)))

	mux.HandleFunc("GET /admin/metrics", apiCfg.handleMetrics)

	mux.HandleFunc("GET /api/chirps", apiCfg.handleQueryChirps)

	mux.HandleFunc("POST /admin/reset", apiCfg.handleReset)

	mux.HandleFunc("/api/chirps", apiCfg.handleChirp)

	mux.HandleFunc("POST /api/users", apiCfg.handleCreateUser)

	mux.HandleFunc("POST /api/login", apiCfg.handleLogin)

	server := &http.Server{
		Handler: mux,
		Addr:    ":8080",
	}
	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}

}
