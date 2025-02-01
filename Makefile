
.PHONY: build
build:
	ls -GFlash
	go build -o ./tmp/main cmd/*.go

.PHONY: run
run: build
	./tmp/main