Host = 127.0.0.1
Port = 5432
User = postgres
Password = postgres
Database = secretsanta
SSLmode =  disable

connect_db = postgres "host=$(Host) port=$(Port) user=$(User) password=$(Password) database=$(Database) sslmode=$(SSLmode) "

migrate-create:
	goose -dir migrations create todo sql
migrate-up:
	goose -dir migrations $(connect_db) up
migrate-down:
	goose -dir migrations $(connect_db) down
migrate-status:
	goose $(connect_db) status