GO_FLAGS = -ldflags "-X 'github.com/leighmacdonald/verimapcom/consts.BuildVersion=`git describe --long`'"

all: build

vet:
	@go vet . ./...

fmt:
	@go fmt . ./...

gen:
	@protoc --go_out=plugins=grpc:. --go_opt=paths=source_relative pb/rpc.proto
	@protoc pb/rpc.proto    --js_out=import_style=commonjs:frontend/src/app \
		--grpc-web_out=import_style=commonjs,mode=grpcwebtext:frontend/src/app

yarn_install:
	@cd frontend && yarn install && cd ..

frontend:
	@cd frontend && yarn run build && cd ..

watch:
	@cd frontend && yarn run watch && cd ..

demo:
	@go run main.go demo

client:
	@go run main.go

serve:
	@go run main.go serve

deps: yarn_install
	@go get github.com/golangci/golangci-lint/cmd/golangci-lint@v1.27.0

build: clean gen fmt lint vet
	@go build $(GO_FLAGS)

run:
	@go run $(GO_FLAGS) -race main.go

install: deps
	@go install $(GO_FLAGS) ./...

test:
	@go test $(GO_FLAGS) -race -cover . ./...

testcover:
	@go test -race -coverprofile c.out $(GO_FLAGS) ./...

lint:
	@golangci-lint run

bench:
	@go test -run=NONE -bench=. $(GO_FLAGS) ./...

clean:
	@go clean $(GO_FLAGS) -i
	@rm -rf ./dist