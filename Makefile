.PHONY: affiliate-system

build: affiliate-system
affiliate-system:
	go build -o ./affiliate-system ./cmd/main.go
clean:
	rm -i -f affiliate-system

run-test:
	go test -v ./internal/app/affiliate/caching/test
	go test -v ./internal/app/affiliate/util/test
	# go test -v ./internal/app/affiliate/watcher/test
	go test -v ./test

restart: stop clean build start
	@echo "Affiliate System restarted!"

build-service: clean build
	@echo "Restart service with cmd: 'systemctl restart affiliate-system'"
	systemctl restart affiliate-system

start: build
	@echo "Starting the Affiliate System..."
	@env DB_PASSWORD=${DB_PASSWORD} ./affiliate-system &
	@echo "Affiliate System running!"

stop:
	@echo "Stopping the Affiliate System..."
	@-pkill -SIGTERM -f "affiliate-system"
	@echo "Stopped Affiliate System"