start:
	docker compose up -d

start-5:
	docker compose up -d --scale=5

down:
	docker compose down

rebuild:
	docker compose up --build -d

test:
	k6 run ./test/k6.js