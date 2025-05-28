include .env.local
include .env.container
.PHONY: docker-commit docker-compose compile linting run help

END=
DEVELOPMENT_STAGE=demo
APPVERSION=v1.0.0
mode=

ifndef DEVELOPMENT_STAGE
$(error DEVELOPMENT_STAGE is not set. Usage: make run DEVELOPMENT_STAGE=stage)
endif

ifneq ($(filter $(DEVELOPMENT_STAGE),demo test stage prod) , $(DEVELOPMENT_STAGE))
$(error Incorrect DEVELOPMENT_STAGE variable. Use on of suggestions: demo test stage prod)
endif

ifeq ($(OS),Windows_NT)
END := .exe
endif

BINARY_NAME = main$(END)

## compile: Собирает бинарник, запускает линтеры и тесты (tools + test + build)
compile: tools test build
## docker-commit: Собирает и пушит Docker-образ (docker-build + docker-push)
docker-commit: docker-push

docker-compose:
	docker compose --env-file .env.container up --build

docker-push: docker-compose
	docker push $(DOCKER_USERNAME)/$(DEVELOPMENT_STAGE)/demo-platform-handling:$(VERSION)

tools:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo " > Install linter and mockgen dependencies here"

linting:
	golangci-lint run ./...
#--config=.golangci.yaml


test:
	@echo " > Run tests, profiling, coverage checks here"

build: clean
	mkdir bin
	set ADDR="localhost:3333" && go build -o ./bin/${BINARY_NAME} ./cmd/app/main.go

## run: Запускает собранный бинарник или контейнер
run:
ifeq ($(mode),local)
	bin\${BINARY_NAME} -mode=local
else ifeq ($(mode),docker)
	docker compose --env-file .env.container up --build
else
	$(error Unknown mode "$(mode)". Use 'local' or 'docker')
endif


RM = rmdir /S /Q
ifneq ($(OS),Windows_NT)
	RM = rm -rf
endif

clean:
	mkdir bin 2>nul || rem
	go clean
	$(RM) bin

## help: Выводит список доступных команд Makefile , работает только в Unix-системах
help:
	@echo "Help commands list for make :"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'