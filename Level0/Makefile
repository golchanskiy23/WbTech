include .env
.PHONY: docker-commit docker-compose compile run help

END=
DEVELOPMENT_STAGE=demo
VERSION=v1.0.0

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
docker-commit: docker-build docker-push

# +build - ?
docker-build:
	set ADDR=":3333" && docker-compose up --build

docker-push: docker-build
	docker push $(DOCKER_USERNAME)/$(DEVELOPMENT_STAGE)/demo-platform-handling:$(VERSION)

tools:
	@echo " > Install linter and mockgen dependencies here"

test:
	@echo " > Run tests, profiling, coverage checks here"

build: clean
	mkdir bin
	set ADDR="localhost:3333" && go build -o ./bin/${BINARY_NAME} ./cmd/app/main.go

## run: Запускает собранный бинарник
run:
	./bin/${BINARY_NAME}

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