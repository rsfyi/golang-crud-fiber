# Getting started with Golang in docker development

- `docker compose run --service-ports web bash`
- Create a go module
  - `go mod init github.com/rsfyi/golang-docker-setup`
- Add Go fiber framework
  - `alias air='$(go env GOPATH)/bin/air'`
  - `go get github.com/gofiber/fiber/v2`
- Run Go project
  - `go run cmd/main.go -b 0.0.0.0`
  - b in above command is binding with localhost
- For hot module reload in golang we will make use of package called go air
  - `https://github.com/cosmtrek/air`
  - `go install github.com/cosmtrek/air@latest`
- We can make use of go database orm
  - `go get -u gorm.io/gorm`
  - `go get -u gorm.io/driver/sqlite`
