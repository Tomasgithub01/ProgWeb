package main

import (
	"ProgWeb/logic"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	sqlc "ProgWeb/db/sqlc"

	_ "github.com/lib/pq"
)

var queries *sqlc.Queries
var ctx context.Context

func main() {
	// En tu main.go, establece la conexión con tu base de datos
	connStr := "user=admin password=#Admin20250915 dbname=tpespecialweb sslmode=disable"
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

	}

}

func handleMain(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" || r.Method != http.MethodGet {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	http.ServeFile(w, r, "index.html")
}
