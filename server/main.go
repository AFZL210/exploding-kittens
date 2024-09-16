package main

import (
	"context"
	"crypto/rand"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"sort"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt/v4"
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
		Addr:     "caching-2daea98a-ak341668-1ede.c.aivencloud.com:17537",
		Username: "default",
		Password: "AVNS_jVtEm0xoUKh04kmumfk",
		TLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
		},
	})

	pong, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		fmt.Println("Error connecting to Redis:", err)
		return
	} else {
		fmt.Println("Connected to Redis:", pong)
	}

	router := mux.NewRouter()

	router.HandleFunc("/api/login", loginHandler).Methods("POST")
	router.HandleFunc("/api/register", registerHandler).Methods("POST")
	router.Handle("/api/getcards", jwtMiddleware(http.HandlerFunc(getCardsHandler))).Methods("GET")
	router.Handle("/api/shuffle", jwtMiddleware(http.HandlerFunc(shuffleHandler))).Methods("GET")
	router.Handle("/api/play", jwtMiddleware(http.HandlerFunc(playHandler))).Methods("POST")
	router.HandleFunc("/api/leaderboard", leaderboardHandler).Methods("GET")
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

var jwtSecret = []byte("your-secret-key")

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func generateJWT(username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &Claims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

func jwtMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var tokenStr string

		authHeader := r.Header.Get("Authorization")
		if authHeader != "" && len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenStr = authHeader[7:]
		} else {
			cookie, err := r.Cookie("token")
			if err != nil {
				http.Error(w, "Unauthorized", http.StatusUnauthorized)
				return
			}
			tokenStr = cookie.Value
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtSecret, nil
		})
		if err != nil || !token.Valid {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "username", claims.Username)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	storedHash, err := redisClient.Get(context.Background(), "user:"+user.Username).Result()
	if err == redis.Nil {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(storedHash), []byte(user.Password)); err == nil {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"username": user.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(),
		})

		tokenString, err := token.SignedString(jwtSecret)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"token": tokenString})
	} else {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
	}
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	exists, err := redisClient.Exists(context.Background(), "user:"+user.Username).Result()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	if exists == 0 {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		err = redisClient.Set(context.Background(), "user:"+user.Username, string(hashedPassword), 0).Err()
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"message": "User registered successfully"})
	} else {
		http.Error(w, "Username already exists", http.StatusConflict)
	}
}

func getCardsHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	gameStateJSON, err := redisClient.Get(context.Background(), "gamestate:"+username).Result()
	if err == redis.Nil {
		gameState := createNewGameState()
		gameStateBytes, _ := json.Marshal(gameState)
		gameStateJSON = string(gameStateBytes)
		err = redisClient.Set(context.Background(), "gamestate:"+username, gameStateJSON, 0).Err()
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
	users, err := redisClient.Keys(context.Background(), "user:*").Result()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var leaderboard []LeaderboardEntry
	for _, user := range users {
		username := user[5:]

		scoreStr, err := redisClient.Get(context.Background(), "score:"+username).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
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
	users, err := redisClient.Keys(context.Background(), "user:*").Result()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var leaderboard []LeaderboardEntry
	for _, user := range users {
		userKey := user[5:]
		scoreStr, err := redisClient.Get(context.Background(), "score:"+userKey).Result()
		if err == redis.Nil {
			continue
		} else if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			return
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
	err := redisClient.Set(context.Background(), "gamestate:"+username, string(gameStateBytes), 0).Err()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
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
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")

	var requestBody struct {
		Index int `json:"index"`
	}
	err := json.NewDecoder(r.Body).Decode(&requestBody)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	gameStateJSON, err := redisClient.Get(context.Background(), "gamestate:"+username).Result()
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

	var gameState GameState
	err = json.Unmarshal([]byte(gameStateJSON), &gameState)
	if err != nil {
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}

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
		err = redisClient.Set(context.Background(), "gamestate:"+username, string(updatedGameStateJSON), 0).Err()
		if err != nil {
			log.Println("Error updating game state:", err)
		}
	}
}

func resetGame(username string) {
	gameState := createNewGameState()
	gameStateBytes, _ := json.Marshal(gameState)
	err := redisClient.Set(context.Background(), "gamestate:"+username, gameStateBytes, 0).Err()
	if err != nil {
		log.Println("Error resetting game:", err)
	}
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
	err := redisClient.Incr(context.Background(), "score:"+username).Err()
	if err != nil {
		log.Println("Error incrementing score:", err)
	}
}
