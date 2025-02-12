dcup:
	@docker-compose up --build

dcdown:
	@docker-compose down

run: dcup

stop: dcdown