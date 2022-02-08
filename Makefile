# Makefile 基础教程
# https://www.kancloud.cn/kancloud/make-command/45596
# https://seisman.github.io/how-to-write-makefile/overview.html

.PHONY:build
build:
	go build -v .

.PHONY:run
run:
	go run main.go

.PHONY:test
test:
	go test ./...
