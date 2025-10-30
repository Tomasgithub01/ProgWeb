test:
	hurl --test script.hurl 

dev: 
	docker compose -f docker-compose.dev.yml up --build

prod:
	docker compose --build 