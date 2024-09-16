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
	"strconv"

	"github.com/go-redis/redis/v8"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Card struct {
	CardType  string `json:"cardType"`
	IsFlipped bool   `json:"isFlipped"`
}

type GameState struct {
	Cards          []Card `json:"cards"`
	DefuseCards    int    `json:"defuseCards"`
	RemainingCards int    `json:"remainingCards"`
	IsWon          bool   `json:"isWon"`
	IsLost         bool   `json:"isLost"`
}

type LeaderboardEntry struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}

var redisClient *redis.Client

func main() {
	redisClient = redis.NewClient(&redis.Options{
		Addr: "rediss://default:AVNS_jVtEm0xoUKh04kmumfk@caching-2daea98a-ak341668-1ede.c.aivencloud.com:17537",
	})

	router := mux.NewRouter()

	router.HandleFunc("/api/login", loginHandler).Methods("POST")
	router.HandleFunc("/api/register", registerHandler).Methods("POST")
	router.HandleFunc("/api/getcards", getCardsHandler).Methods("GET")
	router.HandleFunc("/api/leaderboard", leaderboardHandler).Methods("GET")
	router.HandleFunc("/api/shuffle", shuffleHandler).Methods("GET")
	router.HandleFunc("/api/play", playHandler).Methods("POST")
	router.HandleFunc("/api/user-rank", userRankHandler).Methods("GET")

	corsAllowedOrigins := handlers.AllowedOrigins([]string{"*"})
	corsAllowedMethods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
	corsAllowedHeaders := handlers.AllowedHeaders([]string{"Content-Type", "Authorization"})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Fatal(http.ListenAndServe(":"+port, handlers.CORS(corsAllowedOrigins, corsAllowedMethods, corsAllowedHeaders)(router)))
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

	var gameState GameState
	json.Unmarshal([]byte(gameStateJSON), &gameState)

	response := struct {
		Cards          []Card `json:"cards"`
		IsWon          bool   `json:"isWon"`
		IsLost         bool   `json:"isLost"`
		DefuseCount    int    `json:"defuseCount"`
		RemainingCards int    `json:"remainingCards"`
	}{
		Cards:          gameState.Cards,
		IsWon:          gameState.IsWon,
		IsLost:         gameState.IsLost,
		DefuseCount:    gameState.DefuseCards,
		RemainingCards: gameState.RemainingCards,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func leaderboardHandler(w http.ResponseWriter, r *http.Request) {
	users, _ := redisClient.Keys(r.Context(), "user:*").Result()

	var leaderboard []LeaderboardEntry
	for _, user := range users {
		username := user[5:]

		scoreStr, err := redisClient.Get(r.Context(), "score:"+username).Result()
		if err == redis.Nil {
			continue
		}
		score, _ := strconv.Atoi(scoreStr)

		leaderboard = append(leaderboard, LeaderboardEntry{Username: username, Score: score})
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(leaderboard)
}

func userRankHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	users, _ := redisClient.Keys(r.Context(), "user:*").Result()

	var leaderboard []LeaderboardEntry
	for _, user := range users {
		userKey := user[5:]
		scoreStr, err := redisClient.Get(r.Context(), "score:"+userKey).Result()
		if err == redis.Nil {
			continue
		}
		score, _ := strconv.Atoi(scoreStr)
		leaderboard = append(leaderboard, LeaderboardEntry{Username: userKey, Score: score})
	}

	sort.Slice(leaderboard, func(i, j int) bool {
		return leaderboard[i].Score > leaderboard[j].Score
	})

	for rank, entry := range leaderboard {
		if entry.Username == username {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]interface{}{
				"username": username,
				"rank":     rank + 1,
				"score":    entry.Score,
			})
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func shuffleHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	gameState := createNewGameState()
	gameStateBytes, _ := json.Marshal(gameState)
	redisClient.Set(r.Context(), "gamestate:"+username, string(gameStateBytes), 0)

	response := struct {
		Cards          []Card `json:"cards"`
		IsWon          bool   `json:"isWon"`
		IsLost         bool   `json:"isLost"`
		DefuseCount    int    `json:"defuseCount"`
		RemainingCards int    `json:"remainingCards"`
	}{
		Cards:          gameState.Cards,
		IsWon:          gameState.IsWon,
		IsLost:         gameState.IsLost,
		DefuseCount:    gameState.DefuseCards,
		RemainingCards: gameState.RemainingCards,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
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

	if gameState.IsWon || gameState.IsLost {
		http.Error(w, "Game is already over", http.StatusBadRequest)
		return
	}

	if requestBody.Index < 0 || requestBody.Index >= len(gameState.Cards) {
		http.Error(w, "Invalid card index", http.StatusBadRequest)
		return
	}

	if gameState.Cards[requestBody.Index].IsFlipped {
		http.Error(w, "Card is already flipped", http.StatusBadRequest)
		return
	}

	gameState.Cards[requestBody.Index].IsFlipped = true
	drawnCard := gameState.Cards[requestBody.Index]

	switch drawnCard.CardType {
	case "Cat":
		gameState.RemainingCards--
	case "Defuse":
		gameState.DefuseCards++
		gameState.RemainingCards--
	case "Shuffle":
		if gameState.RemainingCards == 1 {
			gameState.IsWon = true
		} else {
			gameState.IsLost = true
		}
	case "Bomb":
		if gameState.DefuseCards > 0 {
			gameState.DefuseCards--
			gameState.RemainingCards--
		} else {
			gameState.IsLost = true
		}
	}

	if !gameState.IsWon && !gameState.IsLost && gameState.RemainingCards == 0 {
		gameState.IsWon = true
	}

	response := struct {
		Cards          []Card `json:"cards"`
		IsWon          bool   `json:"isWon"`
		IsLost         bool   `json:"isLost"`
		DefuseCount    int    `json:"defuseCount"`
		RemainingCards int    `json:"remainingCards"`
	}{
		Cards:          gameState.Cards,
		IsWon:          gameState.IsWon,
		IsLost:         gameState.IsLost,
		DefuseCount:    gameState.DefuseCards,
		RemainingCards: gameState.RemainingCards,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)

	if gameState.IsWon {
		incrementScore(username)
	}

	if gameState.IsWon || gameState.IsLost {
		resetGame(username)
	} else {
		updatedGameStateJSON, _ := json.Marshal(gameState)
		redisClient.Set(r.Context(), "gamestate:"+username, string(updatedGameStateJSON), 0)
	}
}

func resetGame(username string) {
	gameState := createNewGameState()
	gameStateBytes, _ := json.Marshal(gameState)
	redisClient.Set(context.Background(), "gamestate:"+username, gameStateBytes, 0)
}

func createNewGameState() GameState {
	cardTypes := []string{"Cat", "Defuse", "Shuffle", "Bomb", "Cat"}
	cards := make([]Card, len(cardTypes))
	for i, cardType := range cardTypes {
		cards[i] = Card{CardType: cardType, IsFlipped: false}
	}
	shuffleCards(cards)
	return GameState{
		Cards:          cards,
		DefuseCards:    0,
		RemainingCards: len(cards),
		IsWon:          false,
		IsLost:         false,
	}
}

func shuffleCards(cards []Card) {
	for i := len(cards) - 1; i > 0; i-- {
		j, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		cards[i], cards[j.Int64()] = cards[j.Int64()], cards[i]
	}
}

func incrementScore(username string) {
	redisClient.Incr(context.Background(), "score:"+username)
}
