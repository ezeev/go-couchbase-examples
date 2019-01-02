build:
	#dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/hello fn/hello/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/world fn/world/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/signal fn/signal/main.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose
