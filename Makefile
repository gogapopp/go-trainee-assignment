docker-compose-up:
	@docker-compose up --build

docker-compose-down:
	@docker-compose down

run: docker-compose-up

stop: docker-compose-down