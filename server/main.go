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
	json.NewDecoder(r.Body).Decode(&user)

	storedHash, _ := redisClient.Get(r.Context(), "user:"+user.Username).Result()
	if bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password)) == nil {
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Login successful"})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	json.NewDecoder(r.Body).Decode(&user)

	if exists, _ := redisClient.Exists(r.Context(), "user:"+user.Username).Result(); exists == 0 {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		redisClient.Set(r.Context(), "user:"+user.Username, string(hashedPassword), 0)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
	} else {
		http.Error(w, "Username already exists", http.StatusConflict)
	}
}

func getCardsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	gameStateJSON, err := redisClient.Get(r.Context(), "gamestate:"+username).Result()
	if err != nil {
		gameState := createNewGameState()
		gameStateBytes, _ := json.Marshal(gameState)
		gameStateJSON = string(gameStateBytes)
		redisClient.Set(r.Context(), "gamestate:"+username, gameStateJSON, 0)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(gameStateJSON))
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	users, _ := redisClient.Keys(r.Context(), "user:*").Result()

	var leaderboard []LeaderboardEntry
	for _, user := range users {
		username := user[5:]
		score, _ := redisClient.Get(r.Context(), "score:"+username).Int()
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
	gameState := createNewGameState()
	gameStateBytes, _ := json.Marshal(gameState)
	redisClient.Set(r.Context(), "gamestate:"+username, string(gameStateBytes), 0)

	w.Header().Set("Content-Type", "application/json")
	w.Write(gameStateBytes)
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	var requestBody struct {
		Index int `json:"index"`
	}
	json.NewDecoder(r.Body).Decode(&requestBody)

	gameStateJSON, _ := redisClient.Get(r.Context(), "gamestate:"+username).Result()

	var gameState GameState
	json.Unmarshal([]byte(gameStateJSON), &gameState)

	drawnCard := gameState.Cards[requestBody.Index]

	switch drawnCard {
	case "exploding_kitten":
		if gameState.DefuseCards > 0 {
			gameState.DefuseCards--
		} else {
			gameState.GameOver = true
			gameState.Won = false
			resetGame(username)
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

	updatedGameStateJSON, _ := json.Marshal(gameState)
	redisClient.Set(r.Context(), "gamestate:"+username, string(updatedGameStateJSON), 0)

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

func incrementScore(username string) {
	redisClient.Incr(context.Background(), "score:"+username)
}
