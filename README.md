# Avito-Backend-trainee-assignment-autumn-2025
## Pull request assignment service

Техническое задание можно найти [здесь](docs/Task.md)

С API можно познакомиться [здесь](openapi.yml)

___
## **Запуск**
### Запуск тестов
`cp .env-test .env.test` \
`docker compose -f docker-compose.test.yml up -d`\
`go test ./tests/...`

### Запуск приложения
`cp .env-sample .env` \
`docker compose -f docker-compose.yml up -d`
___

## **Стек решения**
- Golang
- PostgreSQL
- docker

## **Библиотеки, использованные в задании**
- chi (роутеры)
- squirell (SQL-query builder)
- sqlx (БД)
- testify (Тесты)

___ 
## **Дополнительные задания**
- Добавлен эндпоинт статистики, сортирует пользователей по количеству PR, в которых пользователь назначен ревьюером, также есть статистика открытых и смерженных PR (доступен по эндпоинту `/stats/users`)
- Реализовано интеграционное тестирование, для запуска тестов требуется разворачивать другую БД через docker-compose.test.yml, чтобы не менять состояние базы данных
- Описана конфигурация линтера, а также сделан CI/CD для github с проверкой на линтер

## **Вопросы, мои решения**
- В openapi.yml в required нет флага needMoreReviewers для PR, в ТЗ есть - решил добавить
- В openapi.yml в properties есть createdAt, mergedAt, в ТЗ нет - решил добавить
