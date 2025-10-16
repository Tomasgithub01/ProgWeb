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
	// En tu main.go, establece la conexión con tu base de datos
	connStr := "host=db user=admin password=#Admin20250915 dbname=tpespecialweb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()
	//Crea una instancia del repositorio generado por sqlc para poder usarlo en tus handlers.
	queries = sqlc.New(db)
	ctx = context.Background()

	http.HandleFunc("/", handleMain)
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

// Handler para /games - maneja operaciones relacionadas con juegos
func gamesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		var newGame sqlc.Game
		err := json.NewDecoder(r.Body).Decode(&newGame)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = logic.ValidateGame(newGame)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdGame, err := queries.CreateGame(ctx,
			sqlc.CreateGameParams{
				Name:        newGame.Name,
				Description: newGame.Description,
				Image:       newGame.Image,
				Link:        newGame.Link,
			})

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdGame)
	} else if r.Method == http.MethodGet {
		games, err := queries.ListGames(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(games)
	}

}

// Handler para /games/{id} - maneja operaciones específicas de un juego
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

// Obtener un juego por ID
func getGame(w http.ResponseWriter, r *http.Request, id int32) {
	game, err := queries.GetGame(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(game)
}

// Actualizar un juego por ID
func updateGame(w http.ResponseWriter, r *http.Request, id int32) {
	var updatedGame sqlc.Game
	err := json.NewDecoder(r.Body).Decode(&updatedGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = logic.ValidateGame(updatedGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = queries.UpdateGame(ctx, sqlc.UpdateGameParams{
		ID:          id,
		Name:        updatedGame.Name,
		Description: updatedGame.Description,
		Image:       updatedGame.Image,
		Link:        updatedGame.Link,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Eliminar un juego por ID
func deleteGame(w http.ResponseWriter, r *http.Request, id int32) {
	err := queries.Delete(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Handler para /users - maneja operaciones relacionadas con usuarios
func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newUser sqlc.User
		err := json.NewDecoder(r.Body).Decode(&newUser)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		err = logic.ValidateUser(newUser)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		createdUser, err := queries.CreateUser(ctx,
			sqlc.CreateUserParams{
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
	} else if r.Method == http.MethodGet {
		users, err := queries.ListUsers(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}

}

// Handler para /users/{id} - maneja operaciones específicas de un usuario
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
	err := json.NewDecoder(r.Body).Decode(&updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = logic.ValidateUser(updatedUser)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = queries.UpdateUser(ctx, sqlc.UpdateUserParams{
		ID:       id,
		Name:     updatedUser.Name,
		Password: updatedUser.Password,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func deleteUser(w http.ResponseWriter, r *http.Request, id int32) {
	err := queries.DeleteUser(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Handler para /plays - maneja operaciones relacionadas con la relación de juegos y usuarios
func playsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var newPlay sqlc.Play
		err := json.NewDecoder(r.Body).Decode(&newPlay)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		_, err = queries.GetGame(ctx, newPlay.IDGame)
		if err != nil {
			http.Error(w, "Game not found", http.StatusBadRequest)
			return
		}
		_, err = queries.GetUser(ctx, newPlay.IDUser)
		if err != nil {
			http.Error(w, "User not found", http.StatusBadRequest)
			return
		}

		createdPlays, err := queries.CreateUserPlaysGame(ctx,
			sqlc.CreateUserPlaysGameParams(newPlay))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(createdPlays)
	} else if r.Method == http.MethodGet {
		plays, err := queries.ListUserPlaysGames(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(plays)
	}

}

// Handler para /plays/{game_id}/{user_id} - maneja operaciones específicas de una relación
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

// Obtener una relación juego-usuario por IDs
func getPlays(w http.ResponseWriter, r *http.Request, gameID int32, userID int32) {
	play, err := queries.GetUserPlaysGame(ctx, sqlc.GetUserPlaysGameParams{
		IDGame: int32(gameID),
		IDUser: int32(userID),
	})
	if err != nil {
		http.Error(w, "Game is not played by user", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(play)
}

// Actualizar una relación juego-usuario
func updatePlays(w http.ResponseWriter, r *http.Request) {
	var updatedPlays sqlc.Play
	err := json.NewDecoder(r.Body).Decode(&updatedPlays)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = queries.GetGame(ctx, updatedPlays.IDGame)
	if err != nil {
		http.Error(w, "Game not found", http.StatusBadRequest)
		return
	}
	_, err = queries.GetUser(ctx, updatedPlays.IDUser)
	if err != nil {
		http.Error(w, "User not found", http.StatusBadRequest)
		return
	}

	err = queries.UpdateUserPlaysGame(ctx, sqlc.UpdateUserPlaysGameParams(updatedPlays))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Eliminar una relación juego-usuario
func deletePlays(w http.ResponseWriter, r *http.Request, gameID int32, userID int32) {
	err := queries.DeleteUserPlaysGame(ctx, sqlc.DeleteUserPlaysGameParams{
		IDGame: int32(gameID),
		IDUser: int32(userID),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Handler para el main - sirve el archivo HTML principal
func handleMain(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, "index.html")
}
