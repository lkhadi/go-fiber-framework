GOLANG FIBER API
##STEP##
1. go mod tidy
2. go install -tags "mysql" github.com/golang-migrate/migrate/v4/cmd/migrate@latest
3. go build main.go -o web