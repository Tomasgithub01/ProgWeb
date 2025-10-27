APP_NAME=my-app

DB_URL="postgres://admin:#Admin20250915@localhost:5434/tpespecialweb?sslmode=disable"

# Significa que si se ejecuta 'make' por defecto es 'make build'
all: build

# 'make run' ejecuta el air, antes levanta la base de datos. Luego air ejecuta 'make build' cuando ve un cambio.
run: db
	@echo ">= Esperando que la base de datos esté lista..."
	sleep 10
	@echo ">= Iniciando aplicación con Air..."
	@air

generate:
	@echo ">= Generating SQLC code..."
	@sqlc generate 
#	@echo ">= Generating Templ code..."
#	@templ generate

#db/migrate:

# Cada vez que air vea un cambio en un archivo ejecuta 'make generate' y 'make build'
build: generate
	@echo ">= Building application..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) .

clean:
	@rm -f $(APP_NAME)

# Levantar solo la base de datos
db: 
	docker compose up db

stop-db: 
	@docker compose stop db 