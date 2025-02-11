# go-trainee-assignment
## Протокол

Описан в файле [schema.yaml](schema.yaml) и [schema.json](schema.json).

## Как запустить (с докером)

создать .env из .env.example
```bash
cp .env.example .env
```

```bash
docker-compose up --build
```

## Остановить приложение

```bash
docker-compose-down
```

<!-- Постман коллекция: [documenter.getpostman.com](https://documenter.getpostman.com/view/2612412453/2sA123123DpC)  

## Generating code from a specification

Install [oapi-codegen](https://github.com/deepmap/oapi-codegen/) and generate:

```bash
oapi-codegen -package=handler -generate="chi-server,types,spec" api.yaml > internal/handler/api.gen.go
```

oapi-gen:
	@oapi-codegen -package=handler -generate="chi-server,types,spec" api.yaml > internal/handler/api.gen.go

https://github.com/avito-tech/tech-internship/tree/main/Tech%20Internships/Backend/Backend-trainee-assignment-winter-2025 -->