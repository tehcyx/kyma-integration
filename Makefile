default: build

KIS_VERSION=$(shell cat VERSION)
GIT_COMMIT=$(shell git rev-list -1 HEAD)

build: test cover
	go build -ldflags "-X github.com/tehcyx/kyma-github-integration/pkg/kis/cmd.Version=${KIS_VERSION} -X github.com/tehcyx/kyma-github-integration/pkg/lic/cmd.GitCommit=${GIT_COMMIT}" -i -o bin/kis ./cmd/kis

docker: test cover
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/tehcyx/kyma-github-integration/pkg/kis/cmd.Version=${LIC_VERSION} -X github.com/tehcyx/kyma-github-integration/pkg/lic/cmd.GitCommit=${GIT_COMMIT} -s" -a -installsuffix cgo -i -o bin/kisdocker ./cmd/kis

install: build
	go install

test:
	go test ./...

cover:
	go test ./... -cover

clean:
	rm -rf bin