dcup:
	@docker-compose up --build

dcdown:
	@docker-compose down

run: docker-compose-up

stop: docker-compose-down