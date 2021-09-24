GOOS=linux

build:
	@go build -o bin/webhook main.go
	@echo "[OK] Project built!"

deploy:
	@sls deploy --verbose
	@echo "[OK] Application was deployed!"
