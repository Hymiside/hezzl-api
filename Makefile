up:
	migrate -database "postgres://hymiside:@127.0.0.1:5432/hezzl?sslmode=disable" -path migrations/postgres up