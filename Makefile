default: build

KIS_VERSION=$(shell cat VERSION)
GIT_COMMIT=$(shell git rev-list -1 HEAD)

run: bin/templates
	go run cmd/kis/main.go

build: bin/templates test cover
	go build -ldflags "-X github.com/tehcyx/kyma-integration/pkg/cmd.Version=${KIS_VERSION} -X github.com/tehcyx/kyma-integration/pkg/cmd.GitCommit=${GIT_COMMIT}" -i -o bin/kis ./cmd/kis

docker: bin/templates test cover
	CGO_ENABLED=0 GOOS=linux go build -ldflags "-X github.com/tehcyx/kyma-integration/pkg/cmd.Version=${KIS_VERSION} -X github.com/tehcyx/kyma-integration/pkg/cmd.GitCommit=${GIT_COMMIT} -s" -a -installsuffix cgo -i -o bin/kisdocker ./cmd/kis

dockertest: bin/templates
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -i -o bin/appdocker ./cmd/test
	docker build -f build/package/Dockerfile.test -t test-app:label .
	kind load docker-image test-app:label --name "config-map"

setupkind:
	kind create cluster --name "config-map" --config kind-config.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.27.0/deploy/static/mandatory.yaml
	kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/nginx-0.27.0/deploy/static/provider/baremetal/service-nodeport.yaml
	kubectl patch deployments -n ingress-nginx nginx-ingress-controller -p '{"spec":{"template":{"spec":{"containers":[{"name":"nginx-ingress-controller","ports":[{"containerPort":80,"hostPort":80},{"containerPort":443,"hostPort":443}]}],"nodeSelector":{"ingress-ready":"true"},"tolerations":[{"key":"node-role.kubernetes.io/master","operator":"Equal","effect":"NoSchedule"}]}}}}'

install: build
	go install

bin/templates:
	mkdir -p internal/tmpl
	go run hack/packtemplates.go
	go fmt github.com/tehcyx/kyma-integration/internal/tmpl

test:
	go test ./...

cover:
	go test ./... -cover

clean:
	rm -rf bin