include .env
export

migrate-up:
	goose -dir ./sql/schema postgres "$(DB_URL_LOCAL)" up

migrate-down:
	goose -dir ./sql/schema postgres "$(DB_URL_LOCAL)" down
