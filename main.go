package main

import (
	sqlc "ProgWeb/db/sqlc"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

func main() {
	// En tu main.go, establece la conexión con tu base de datos
	connStr := "user=admin password=#Admin20250915 dbname=tpespecialweb sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("failed to connect to DB: %v", err)
	}
	defer db.Close()
	//Crea una instancia del repositorio generado por sqlc para poder usarlo en tus handlers.
	queries := sqlc.New(db)
	//ctx := context.Background()

	http.HandleFunc("/", handleMain)

	port := ":8080"
	fmt.Printf("Servidor escuchando en http://localhost%s\n", port)

	err = http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Printf("Error al iniciar el servidor: %s\n", err)
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
