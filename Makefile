# Подключение к бд
DB_DSN := "postgres://postgres:yourpassword@localhost:5432/postgres?sslmode=disable"

# общая команда для работы с миграциями 
MIGRATE := migrate -path ./migrations -database $(DB_DSN)

# создание новой миграции (создает 2 новых файла миграции - up и down)
migrate-new:
	migrate create -ext sql -dir ./migrations ${NAME}

# применение миграция

migrate:
	$(MIGRATE) up

# откат миграций
migrate-down:
	$(MIGRATE) down

# запуск прилы
run:
	go run cmd/app/main.go 

gen:
	oapi-codegen -config openapi/.openapi -include-tags tasks -package tasks openapi/openapi.yaml > ./internal/web/tasks/api.gen.go

lint:
	golangci-lint run -v --color=auto 

