test:
	hurl --test script.hurl 

dev: 
	docker compose -f docker-compose.dev.yml up --build

prod:
	docker compose up --build 

down:
	docker compose -f docker-compose.dev.yml down

testdev:
	docker compose -f docker-compose.dev.yml up -d --build
	docker compose -f docker-compose.dev.yml wait db
	hurl --test script.hurl
