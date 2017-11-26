all:
	go generate
	go build

cover:
	go test -cover
	go test -coverprofile=coverage.out
	go tool cover -html=coverage.out

