package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"os"
	"sort"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GameState struct {
	Cards          []string `json:"cards"`
	DefuseCards    int      `json:"defuseCards"`
	RemainingCards int      `json:"remainingCards"`
	GameOver       bool     `json:"gameOver"`
	Won            bool     `json:"won"`
}

type LeaderboardEntry struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}

var redisClient *redis.Client

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})

	router := mux.NewRouter()

	router.HandleFunc("/api/login", loginHandler).Methods("POST")
	router.HandleFunc("/api/register", registerHandler).Methods("POST")
	router.HandleFunc("/api/getcards", getCardsHandler).Methods("GET")
	router.HandleFunc("/api/leaderboard", leaderboardHandler).Methods("GET")
	router.HandleFunc("/api/shuffle", shuffleHandler).Methods("GET")
	router.HandleFunc("/api/play", playHandler).Methods("POST")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	storedHash, err := redisClient.Get(r.Context(), "user:"+user.Username).Result()
	if err == redis.Nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := redisClient.Exists(r.Context(), "user:"+user.Username).Result()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	if exists == 1 {
		http.Error(w, "Username already exists", http.StatusConflict)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = redisClient.Set(r.Context(), "user:"+user.Username, string(hashedPassword), 0).Err()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
}

func getCardsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	gameStateJSON, err := redisClient.Get(r.Context(), "gamestate:"+username).Result()
	if err == redis.Nil {
		gameState := createNewGameState()
		gameStateBytes, err := json.Marshal(gameState)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		gameStateJSON = string(gameStateBytes)
		err = redisClient.Set(r.Context(), "gamestate:"+username, gameStateJSON, 0).Err()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var gameState GameState
	err = json.Unmarshal([]byte(gameStateJSON), &gameState)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if gameState.GameOver {
		gameState = createNewGameState()
		gameStateBytes, err := json.Marshal(gameState)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		gameStateJSON = string(gameStateBytes)
		err = redisClient.Set(r.Context(), "gamestate:"+username, gameStateJSON, 0).Err()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(gameStateJSON))
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	users, err := redisClient.Keys(r.Context(), "user:*").Result()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var leaderboard []LeaderboardEntry

	for _, user := range users {
		username := user[5:]
		score, err := redisClient.Get(r.Context(), "score:"+username).Int()
		if err == redis.Nil {
			score = 0
		} else if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		leaderboard = append(leaderboard, LeaderboardEntry{Username: username, Score: score})
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}

func shuffleHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	gameState := createNewGameState()
	gameStateBytes, err := json.Marshal(gameState)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	gameStateJSON := string(gameStateBytes)
	err = redisClient.Set(r.Context(), "gamestate:"+username, gameStateJSON, 0).Err()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(gameStateJSON))
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	if username == "" {
		http.Error(w, "Username is required", http.StatusBadRequest)
		return
	}

	var requestBody struct {
		Index int `json:"index"`
	}
	if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameStateJSON, err := redisClient.Get(r.Context(), "gamestate:"+username).Result()
	if err == redis.Nil {
		http.Error(w, "No active game found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var gameState GameState
	err = json.Unmarshal([]byte(gameStateJSON), &gameState)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if gameState.GameOver {
		http.Error(w, "Game is already over", http.StatusBadRequest)
		return
	}

	if len(gameState.Cards) == 0 {
		http.Error(w, "No cards left to draw", http.StatusBadRequest)
		return
	}

	if requestBody.Index < 0 || requestBody.Index >= len(gameState.Cards) {
		http.Error(w, "Invalid card index", http.StatusBadRequest)
		return
	}

	drawnCard := gameState.Cards[requestBody.Index]

	switch drawnCard {
	case "cat":
	case "exploding_kitten":
		if gameState.DefuseCards > 0 {
			gameState.DefuseCards--
		} else {
			gameState.GameOver = true
			gameState.Won = false
			resetGame(username)
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{"message": "Game over! You lost."})
			return
		}
	case "defuse":
		gameState.DefuseCards++
	case "shuffle":
		gameState = createNewGameState()
	}

	if len(gameState.Cards) == 0 && !gameState.GameOver {
		gameState.GameOver = true
		gameState.Won = true
		incrementScore(username)
	}

	updatedGameStateJSON, err := json.Marshal(gameState)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	err = redisClient.Set(r.Context(), "gamestate:"+username, string(updatedGameStateJSON), 0).Err()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(updatedGameStateJSON)
}

func resetGame(username string) {
	gameState := createNewGameState()
	gameStateBytes, _ := json.Marshal(gameState)
	redisClient.Set(context.Background(), "gamestate:"+username, gameStateBytes, 0)
}

func createNewGameState() GameState {
	cards := []string{"cat", "defuse", "shuffle", "exploding_kitten", "cat"}
	shuffleSlice(cards)
	return GameState{
		Cards:          cards,
		DefuseCards:    0,
		RemainingCards: len(cards),
		GameOver:       false,
		Won:            false,
	}
}

func shuffleSlice(slice []string) {
	for i := len(slice) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		slice[i], slice[j.Int64()] = slice[j.Int64()], slice[i]
	}
}

func incrementScore(username string) error {
	return redisClient.Incr(context.Background(), "score:"+username).Err()
}
