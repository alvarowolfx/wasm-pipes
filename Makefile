install:
	go install ./

test:
	go test ./...

cover:
	go test -coverprofile=c.out ./... && go tool cover -html=c.out