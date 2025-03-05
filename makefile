PORT?=8080

lint:
	golangci-lint run .

fmt:
	go fmt .
	golangci-lint run --fix .

dev:
	go run github.com/joho/godotenv/cmd/godotenv@v1.5.1 \
		go run github.com/air-verse/air@v1.52.2 \
			--build.cmd "go build -o .tmp/bin/main ." \
			--build.bin ".tmp/bin/main" \
			--build.exclude_dir "node_modules" \
			--build.include_ext "go" \
			--build.stop_on_error "false" \
			--misc.clean_on_exit true \
			-- -dev -port $(PORT)


build:
	go build -o ./.dist/app .

run: build
	./.dist/app

clean:
	# Remove generated directories
	if [[ -d ".dist" ]]; then rm -r ./.dist; fi
	if [[ -d "tmp" ]]; then rm -r ./tmp; fi
	if [[ -d "bin" ]]; then rm -r ./bin; fi
