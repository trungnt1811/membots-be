.PHONY: membots-be

build: membots-be
membots-be:
	go build -o ./membots-be ./cmd/main.go
clean:
	rm -i -f membots-be

run-test:
	go test -v ./internal/app/membots-be/caching/test
	go test -v ./internal/app/membots-be/util/test
	# go test -v ./internal/app/membots-be/watcher/test
	go test -v ./test

restart: stop clean build start
	@echo "membots-be restarted!"

build-service: clean build
	@echo "Restart service with cmd: 'systemctl restart membots-be'"
	systemctl restart membots-be

start: build
	@echo "Starting the membots-be..."
	@env DB_PASSWORD=${DB_PASSWORD} ./membots-be &
	@echo "membots-be running!"

stop:
	@echo "Stopping the membots-be..."
	@-pkill -SIGTERM -f "membots-be"
	@echo "Stopped membots-be"