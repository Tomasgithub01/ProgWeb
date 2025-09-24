APP_NAME=my-app

DB_URL="postgres://admin:#Admin20250915@localhost:5434/tpespecialweb?sslmode=disable"

all: build

run:
	@air

generate:
	@echo ">= Generating SQLC code..."
	@sqlc generate 
#	@echo ">= Generating Templ code..."
#	@templ generate

#db/migrate:

build:
	@echo ">= Building application..."
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(APP_NAME) .
clean:
	@rm -f $(APP_NAME)