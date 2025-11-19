test:
	hurl --test script.hurl 

dev: 
	docker compose -f docker-compose.dev.yml up --build

prod:
	docker compose up --build 

down:
	docker compose -f docker-compose.dev.yml down -v

testdev:
	docker compose -f docker-compose.dev.yml up -d --build
	@sleep 30
	hurl --test script.hurl

generate: 
	go run github.com/a-h/templ/cmd/templ@latest generate
