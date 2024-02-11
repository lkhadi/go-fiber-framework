GOLANG FIBER API
##STEP##
1. go mod tidy
2. go build main.go -o web


##DATABASE MIGRATION##
1. install package : go install -tags "mysql,mongodb" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
2. create migration file : migrate create -ext sql -dir database/migrations create_table_users
3. run migration : migrate -database "mysql://username:password@tcp(localhost:3306)/database_name" -path database/migrations up


##NOTES##
jwt.go.bak for generate jwt key