
# Exploading Kittens

A luck based card game.

[![Watch the video](https://img.youtube.com/vi/iBNB24ZarnM/0.jpg)](https://youtu.be/iBNB24ZarnM)




## Installation

### Clone the Project
```bash
git clone https://github.com/AFZL210/exploding-kittens.git
```

### Setting Up Client/Frontend
1. Navigate to the `exploding-kittens/client` folder:
   ```bash
   cd exploding-kittens/client
   ```

2. Install the dependencies:
   ```bash
   npm install
   ```

3. Start the development server:
   ```bash
   npm run dev
   ```

### Setting Up Backend

**NOTE**: Make sure you have Go and Redis installed, and Redis is running either locally or on a remote server.

1. Navigate to the `exploding-kittens/server` folder:
   ```bash
   cd exploding-kittens/server
   ```

2. Install Go modules:
   ```bash
   go mod tidy
   ```

3. Update Redis configuration in `main.go`:

   - If Redis is running **locally**, update the configuration as follows:
     ```go
     redisClient := redis.NewClient(&redis.Options{
         Addr: "localhost:6379",
     })
     ```

   - If Redis is running on a **remote server**, configure it like this:
     ```go
     redisClient = redis.NewClient(&redis.Options{
         Addr:     "host:port",
         Username: "Username",
         Password: "Password",
         TLSConfig: &tls.Config{
             MinVersion: tls.VersionTLS12,
         },
     })
     ```

4. Run the server:
   ```bash
   go run main.go
   ```

No go to ```http://localhost:5173/``` to play the game!
## API Reference

#### Register user

```http
  GET /api/register
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `username`      | `string` | **Required**. username |
| `password`      | `string` | **Required**. password for username |



#### Login user

```http
  GET /api/login
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `username`      | `string` | **Required**. username |
| `password`      | `string` | **Required**. password for username |

### Get User's Cards

```http
  GET /api/getcards?username={$Username}
```

### Shuffle/Reset usersgame

```http
  GET /api/shuffle?username={$Username}
```

#### Play Move/Draw a card

```http
  GET /api/play?username={$Username}
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `index`      | `string` | **Required**. index of card to be drawn |

### Get Leaderboard

```http
  GET /api/leaderboard
```



### Get User's Rank

```http
  GET /api/user-rank?username={$Username}
```



## Features

- User game create an account and play game as per given rules
- Show leaderboard to users
- Pagination




## Bonus Features

- Automatically save the game for a user at every stage so the user can continue from where he left off last time.
- Real-time update of points on the leaderboard for all the users if they are playing simultaneously. 

