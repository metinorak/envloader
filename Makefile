generate-mocks:
	mockgen -destination=mocks/env_reader_mock.go -package mocks github.com/metinorak/goenv EnvReader 

test:
	go test -v ./...