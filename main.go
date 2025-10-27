package main

import (
	"ProgWeb/logic"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	sqlc "ProgWeb/db/sqlc"

	_ "github.com/lib/pq"
)

var queries *sqlc.Queries
var ctx context.Context

func main() {
	// Conexión a la base de datos
	connStr := "host=db user=admin password=Admin20250915 dbname=tpespecialweb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()

	queries = sqlc.New(db)
	ctx = context.Background()

	// Handler para la página principal
	http.HandleFunc("/", handleMain)

	// Servir archivos estáticos (CSS, JS, imágenes)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Handlers de API
	http.HandleFunc("/games", gamesHandler)
	http.HandleFunc("/games/", gameHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)
	http.HandleFunc("/plays", playsHandler)
	http.HandleFunc("/plays/", playHandler)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)
	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}

// ================== Handlers ==================

// Handler principal - sirve index.html
func handleMain(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	http.ServeFile(w, r, "static/index.html")
}

// Handler /games
func gamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newGame sqlc.Game
		if err := json.NewDecoder(r.Body).Decode(&newGame); err != nil {
			http.Error(w, "Falla del JSON", http.StatusBadRequest)
			return
		}

		if err := logic.ValidateGame(newGame); err != nil {
			http.Error(w, "Falla de validar el juego", http.StatusBadRequest)
			return
		}

		createdGame, err := queries.CreateGame(ctx, sqlc.CreateGameParams{
			Name:        newGame.Name,
			Description: newGame.Description,
			Image:       newGame.Image,
			Link:        newGame.Link,
		})
		if err != nil {
			http.Error(w, "Falla al insertar en la base", http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdGame)
		return
	}

	if r.Method == http.MethodGet {
		games, err := queries.ListGames(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	}
}

// Handler /games/{id}
func gameHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}
	id := int32(idInt)

	switch r.Method {
	case http.MethodGet:
		getGame(w, r, id)
	case http.MethodPut:
		updateGame(w, r, id)
	case http.MethodDelete:
		deleteGame(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getGame(w http.ResponseWriter, r *http.Request, id int32) {
	game, err := queries.GetGame(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

func updateGame(w http.ResponseWriter, r *http.Request, id int32) {
	var updatedGame sqlc.Game
	if err := json.NewDecoder(r.Body).Decode(&updatedGame); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := queries.UpdateGame(ctx, sqlc.UpdateGameParams{
		ID:          id,
		Name:        updatedGame.Name,
		Description: updatedGame.Description,
		Image:       updatedGame.Image,
		Link:        updatedGame.Link,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func deleteGame(w http.ResponseWriter, r *http.Request, id int32) {
	if err := queries.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ================== Users ==================

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newUser sqlc.User
		if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := logic.ValidateUser(newUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		createdUser, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
			Name:     newUser.Name,
			Password: newUser.Password,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdUser)
		return
	}
	if r.Method == http.MethodGet {
		users, err := queries.ListUsers(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}

func userHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 3 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	idInt, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	id := int32(idInt)

	switch r.Method {
	case http.MethodGet:
		getUser(w, r, id)
	case http.MethodPut:
		updateUser(w, r, id)
	case http.MethodDelete:
		deleteUser(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getUser(w http.ResponseWriter, r *http.Request, id int32) {
	user, err := queries.GetUser(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

func updateUser(w http.ResponseWriter, r *http.Request, id int32) {
	var updatedUser sqlc.User
	if err := json.NewDecoder(r.Body).Decode(&updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := logic.ValidateUser(updatedUser); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err := queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:       id,
		Name:     updatedUser.Name,
		Password: updatedUser.Password,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func deleteUser(w http.ResponseWriter, r *http.Request, id int32) {
	if err := queries.DeleteUser(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ================== Plays ==================

func playsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newPlay sqlc.Play
		if err := json.NewDecoder(r.Body).Decode(&newPlay); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if _, err := queries.GetGame(ctx, newPlay.IDGame); err != nil {
			http.Error(w, "Game not found", http.StatusBadRequest)
			return
		}
		if _, err := queries.GetUser(ctx, newPlay.IDUser); err != nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}
		createdPlay, err := queries.CreateUserPlaysGame(ctx, sqlc.CreateUserPlaysGameParams(newPlay))
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdPlay)
		return
	}
	if r.Method == http.MethodGet {
		plays, err := queries.ListUserPlaysGames(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plays)
	}
}

func playHandler(w http.ResponseWriter, r *http.Request) {
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) != 4 {
		http.Error(w, "Invalid URL", http.StatusBadRequest)
		return
	}
	gameID, err := strconv.Atoi(parts[2])
	if err != nil {
		http.Error(w, "Invalid game ID", http.StatusBadRequest)
		return
	}
	userID, err := strconv.Atoi(parts[3])
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case http.MethodGet:
		getPlays(w, r, int32(gameID), int32(userID))
	case http.MethodPut:
		updatePlays(w, r)
	case http.MethodDelete:
		deletePlays(w, r, int32(gameID), int32(userID))
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getPlays(w http.ResponseWriter, r *http.Request, gameID, userID int32) {
	play, err := queries.GetUserPlaysGame(ctx, sqlc.GetUserPlaysGameParams{
		IDGame: gameID,
		IDUser: userID,
	})
	if err != nil {
		http.Error(w, "Game is not played by user", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(play)
}

func updatePlays(w http.ResponseWriter, r *http.Request) {
	var updatedPlays sqlc.Play
	if err := json.NewDecoder(r.Body).Decode(&updatedPlays); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if _, err := queries.GetGame(ctx, updatedPlays.IDGame); err != nil {
		http.Error(w, "Game not found", http.StatusBadRequest)
		return
	}
	if _, err := queries.GetUser(ctx, updatedPlays.IDUser); err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}
	if err := queries.UpdateUserPlaysGame(ctx, sqlc.UpdateUserPlaysGameParams(updatedPlays)); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func deletePlays(w http.ResponseWriter, r *http.Request, gameID, userID int32) {
	if err := queries.DeleteUserPlaysGame(ctx, sqlc.DeleteUserPlaysGameParams{
		IDGame: gameID,
		IDUser: userID,
	}); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
