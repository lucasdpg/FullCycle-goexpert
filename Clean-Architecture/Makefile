createmigration:
	migrate create -ext=sql -dir=sql/migrations -seq init

migrate:
	sleep 15 && migrate -path=sql/migrations -database "mysql://root:root@tcp(mysql:3306)/orders" -verbose up

migratedown:
	sleep 15 && migrate -path=sql/migrations -database "mysql://root:root@tcp(mysql:3306)/orders" -verbose down

.PHONY: migrate migratedown createmigration