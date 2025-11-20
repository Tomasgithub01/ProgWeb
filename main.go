package main

import (
	sqlc "ProgWeb/db/sqlc"
	"ProgWeb/logic"
	"ProgWeb/views"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/a-h/templ"
	"github.com/google/uuid"
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

	// Handler healthcheck
	http.HandleFunc("/health", handleHealthcheck)

	// Handler del dashboard
	http.HandleFunc("/dashboard", requireLogin(handleDashboard))
	http.HandleFunc("/logout", logoutHandler)

	// Servir archivos estáticos (CSS, JS, imágenes)
	fs := http.FileServer(http.Dir("static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//Logeo
	http.HandleFunc("/login", loginHandler)

	// Handlers de API
	http.HandleFunc("/games", gamesHandler)
	http.HandleFunc("/games/", gameHandler)
	http.HandleFunc("/users", usersHandler)
	http.HandleFunc("/users/", userHandler)
	http.HandleFunc("/plays", playsHandler)
	http.HandleFunc("/plays/", playHandler)
	http.HandleFunc("/steam/search", SearchSteamGames)
	http.HandleFunc("/steam/search/", SearchSteamGames)

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
	//http.ServeFile(w, r, "static/login.html")

	if currentUser(r) != nil {
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	templ.Handler(views.LayoutLogin()).ServeHTTP(w, r)
}

func handleDashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	user := currentUser(r)
	games, err := queries.ListGamesByUserID(ctx, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plays, err := queries.ListPlaysByUserID(ctx, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	searchedGames := []sqlc.Game{}
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	templ.Handler(views.LayoutIndex(games, searchedGames, user, plays)).ServeHTTP(w, r)
}

// Handler /games
func gamesHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parse the form data
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Error parsing form data", http.StatusBadRequest)
			return
		}

		newGame := sqlc.Game{
			Name:        r.FormValue("name"),
			Description: r.FormValue("description"),
			Image:       r.FormValue("image"),
			Link:        r.FormValue("link"),
			Custom:      r.FormValue("custom"),
		}

		if err := logic.ValidateGame(newGame); err != nil {
			http.Error(w, "Falla de validar el juego", http.StatusBadRequest)
			return
		}

		var createdGame sqlc.Game

		existing, err := queries.GetGameByName(ctx, newGame.Name)
		if err == nil && existing.ID != 0 && createdGame.Custom == "0" {
			createdGame = existing
		} else {
			createdGame, err = queries.CreateGame(ctx, sqlc.CreateGameParams{
				Name:        newGame.Name,
				Description: newGame.Description,
				Image:       newGame.Image,
				Link:        newGame.Link,
				Custom:      newGame.Custom,
			})
			if err != nil {
				http.Error(w, "Falla al insertar en la base", http.StatusBadRequest)
				return
			}
		}

		user := currentUser(r)
		_, err = queries.CreateUserPlaysGame(ctx, sqlc.CreateUserPlaysGameParams{
			IDGame: createdGame.ID,
			IDUser: user.ID,
		})

		// Redirect to reload page or show success
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
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
	log.Printf("Game handler: path=%s method=%s\n", r.URL.Path, r.Method)
	parts := strings.Split(r.URL.Path, "/")
	//Simular que la peticion DELETE llega como POST (al no tener HTMX todavia usamos un formulario que solo admite GET y POST)
	if r.Method == http.MethodPost && r.FormValue("_method") == "DELETE" {
		r.Method = http.MethodDelete
	}
	if r.Method == http.MethodPost && r.FormValue("_method") == "PUT" {
		r.Method = http.MethodPut
	}
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
	//var updatedGame sqlc.Game
	/* if err := json.NewDecoder(r.Body).Decode(&updatedGame); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	} */
	if err := r.ParseForm(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := queries.UpdateGame(ctx, sqlc.UpdateGameParams{
		ID:          id,
		Name:        r.FormValue("name"),
		Description: r.FormValue("description"),
		Image:       r.FormValue("image"),
		Link:        r.FormValue("link"),
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func deleteGame(w http.ResponseWriter, r *http.Request, id int32) {
	if err := queries.Delete(ctx, id); err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	//w.WriteHeader(http.StatusNoContent)
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)

}

// ================== Users ==================

func usersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		//var newUser sqlc.User
		/* if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		} */
		if err := r.ParseForm(); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		newUser := sqlc.User{
			Name:     r.FormValue("username"),
			Password: r.FormValue("password"),
		}
		if err := logic.ValidateUser(newUser); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		existing, err := queries.GetUserByName(ctx, newUser.Name)
		if err == nil && existing.ID != 0 {
			http.Error(w, "El usuario ya existe", http.StatusConflict)
			return
		}
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "El usuario ya existe", http.StatusConflict)
			return
		}

		_, err = queries.CreateUser(ctx, sqlc.CreateUserParams{
			Name:     newUser.Name,
			Password: newUser.Password,
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//w.Header().Set("Content-Type", "application/json")
		//w.WriteHeader(http.StatusCreated)
		//fmt.Fprintf(w, "Usuario creado con éxito")
		//json.NewEncoder(w).Encode(createdUser)
		http.Redirect(w, r, "/", http.StatusSeeOther)
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
		createdPlay, err := queries.CreateUserPlaysGame(ctx, sqlc.CreateUserPlaysGameParams{
			IDGame: newPlay.IDGame,
			IDUser: newPlay.IDUser,
		})
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

// Handler /plays/{gameID}/{userID}
func playHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost && r.FormValue("_method") == "DELETE" {
		r.Method = http.MethodDelete
	}

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
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

// para buscar
func SearchSteamGames(w http.ResponseWriter, r *http.Request) {
	log.Printf("Steam handler: path=%s method=%s query=%s\n", r.URL.Path, r.Method, r.URL.RawQuery)
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Missing query", http.StatusBadRequest)
		return
	}

	// Llamar a la API de Steam
	steamURL := fmt.Sprintf("https://store.steampowered.com/api/storesearch/?term=%s&cc=us", url.QueryEscape(query))
	log.Printf("Consultando Steam: %s\n", steamURL)
	resp, err := http.Get(steamURL)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	var data map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	searchedGames := []sqlc.Game{}
	user := currentUser(r)
	games, err := queries.ListGamesByUserID(ctx, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	plays, err := queries.ListPlaysByUserID(ctx, user.ID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	for _, item := range data["items"].([]interface{}) {
		game := sqlc.Game{
			ID:     int32(item.(map[string]interface{})["id"].(float64)),
			Name:   item.(map[string]interface{})["name"].(string),
			Image:  item.(map[string]interface{})["tiny_image"].(string),
			Link:   fmt.Sprintf("https://store.steampowered.com/app/%d", int32(item.(map[string]interface{})["id"].(float64))),
			Custom: "0",
		}
		log.Printf("Juego %s\n", game.Name)
		searchedGames = append(searchedGames, game)

	}
	log.Printf("Total juegos encontrados: %d\n", len(searchedGames))

	templ.Handler(views.LayoutIndex(games, searchedGames, user, plays)).ServeHTTP(w, r)
}

// Aca vemos si podemos hacer lo de inicio de sesión.
// store de auth state en memoria (por ahora simple)
var sessions = map[string]int32{}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Bad Request", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")

		user, err := queries.GetUserByName(ctx, username)
		if err != nil || user.Password != password {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		// generar session id (token corto que representa el auth state)
		sessionID := uuid.NewString()
		sessions[sessionID] = user.ID

		cookie := http.Cookie{
			Name:     "session",
			Value:    sessionID,
			Path:     "/",
			HttpOnly: true, // recomendado en teoría
		}

		http.SetCookie(w, &cookie)
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	http.NotFound(w, r)
}

func currentUser(r *http.Request) *sqlc.User {
	cookie, err := r.Cookie("session")
	if err != nil {
		return nil
	}

	uid, ok := sessions[cookie.Value]
	if !ok {
		return nil
	}

	user, err := queries.GetUser(ctx, uid)
	if err != nil {
		return nil
	}

	return &user
}

func requireLogin(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if currentUser(r) == nil {
			http.Redirect(w, r, "/", http.StatusSeeOther)
			return
		}
		handler(w, r)
	}
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	cookie, _ := r.Cookie("session")
	delete(sessions, cookie.Value)

	expired := http.Cookie{
		Name:   "session",
		Value:  "",
		MaxAge: -1,
		Path:   "/",
	}
	http.SetCookie(w, &expired)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// healthcheck
func handleHealthcheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
