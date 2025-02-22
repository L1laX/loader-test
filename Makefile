start:
	docker compose up -d

start5:
	docker compose up -d --scale app=5

down:
	docker compose down

rebuild:
	docker compose up --build -d

testk6:
	k6 run ./test/k6.js