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

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
	}
}

func gamesHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method == http.MethodPost {
		var newGame sqlc.Game
		err := json.NewDecoder(r.Body).Decode(&newGame)

		if err != nil {
			http.Error(w, "Falla del JSON", http.StatusBadRequest)
			return
		}

		err = logic.ValidateGame(newGame)
		if err != nil {
			http.Error(w, "Falla de validar el juego", http.StatusBadRequest)
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
			http.Error(w, "Falla al insertar en la base", http.StatusBadRequest)
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
	err := json.NewDecoder(r.Body).Decode(&updatedGame)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// err = logic.ValidateGame(updatedGame)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusBadRequest)
	// 	return
	// }

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

	//.WriteHeader(http.StatusNoContent)
}

func deleteGame(w http.ResponseWriter, r *http.Request, id int32) {
	err := queries.Delete(ctx, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func handleMain(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, "index.html")
}
