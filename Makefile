.PHONY: install
install: main.go weather.go
	go install ./...
