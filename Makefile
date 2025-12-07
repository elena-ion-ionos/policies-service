.PHONY: docker-build

GOLANGCI_LINT_VERSION := v1.64.5
GO ?= go
GO_RUN_TOOLS ?= $(GO) run -modfile ./tools/go.mod
GO_TEST = $(GO_RUN_TOOLS) gotest.tools/gotestsum --format pkgname
VERSION ?= 0.0.1
LDFLAGS ?= -ldflags "-X main.version=$(VERSION)"
TAG ?= policies-service:latest
ARGS ?= ""
GITHUB_USER ?= ""# make sure you dontt commit this with your real username
GITHUB_TOKEN ?= ""# make sure you dontt commit this with your real token
GITHUB_PRIVATE_PATH ?= "ionos-cloud"

.PHONY: run
run:
	@-$(GO) run $(LDFLAGS) ./cmd/service $(ARGS)

.PHONY: build
build: git-config
	$(GO) build $(LDFLAGS) ./cmd/service

tag-generate:
	TIME=`date +%Y%m%d%H%M%S`; \
	BUILD_ID='workshop'; \

.PHONY: generate
generate:
	$(GO) generate ./...

.PHONY: api
api: internal/api/api.gen.go

internal/api:
	mkdir -p $@

.PHONY: internal/api/api.gen.go
internal/api/api.gen.go: openapi/openapi.yaml openapi/config.yaml internal/api
	$(GO_RUN_TOOLS) github.com/deepmap/oapi-codegen/cmd/oapi-codegen \
		-config ./openapi/config.yaml $< > $@

docker-build: ## Build docker image
	docker build -t $(TAG) --secret id=github_token .

docker-push: docker-build## Push image
	docker push $(TAG)

.PHONY: clean
clean:
	rm -rf build

.PHONY: test
## make test will run unit tests with short flag (without unit tests that require a db connection)
test: git-config
	$(GO) mod tidy
	mkdir -p build/reports
	$(GO_TEST) --junitfile build/reports/unit-test.xml -- -short -tags unit -p 1 -race ./... -count=1 -cover -coverprofile build/reports/unit-test-coverage.out

.PHONY: test-all
## make test-all will run all unit tests (including unit tests that require a db connection) on CI environment
test-all:
	$(GO) mod tidy
	mkdir -p build/reports
	$(GO_TEST) --junitfile build/reports/unit-test.xml -- -tags unit -p 1 -race ./... -count=1 -cover -coverprofile build/reports/unit-test-coverage.out

.PHONY: test-local-all
## make test-local-all will run all unit tests (including unit tests that require a db connection) on local environment
test-local-all: git-config db-create
	$(GO) mod tidy
	mkdir -p build/reports
	$(GO_TEST) --junitfile build/reports/unit-test.xml -- -tags unit -p 1 -race ./... -count=1 -cover -coverprofile build/reports/unit-test-coverage.out


.PHONY: vet
vet: ## Run go vet against code.
	$(GO) vet -tags unit ./...
	
.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GO_RUN_TOOLS) mvdan.cc/gofumpt -w .

git-config:
	git config --global \
	url."https://$(GITHUB_USER):$(GITHUB_TOKEN)@github.com/$(GITHUB_PRIVATE_PATH)".insteadOf \
	"https://github.com/$(GITHUB_PRIVATE_PATH)"

.PHONY: lint
lint: 
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run --build-tags unit

.PHONY: lint-with-fix
lint-with-fix:
	$(GO) run github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run --fix --build-tags unit

.PHONY: create-config
create-config:
	@echo "kind: Cluster" > config.yaml
	@echo "apiVersion: kind.x-k8s.io/v1alpha4" >> config.yaml
	@echo "nodes:" >> config.yaml
	@echo "- role: control-plane" >> config.yaml
	@echo "  kubeadmConfigPatches:" >> config.yaml
	@echo "  - |" >> config.yaml
	@echo "    kind: InitConfiguration" >> config.yaml
	@echo "    nodeRegistration:" >> config.yaml
	@echo "      kubeletExtraArgs:" >> config.yaml
	@echo "        node-labels: \"ingress-ready=true\"" >> config.yaml
	@echo "  extraPortMappings:" >> config.yaml
	@echo "  - containerPort: 80" >> config.yaml
	@echo "    hostPort: 80" >> config.yaml
	@echo "    protocol: TCP" >> config.yaml
	@echo "  - containerPort: 443" >> config.yaml
	@echo "    hostPort: 443" >> config.yaml
	@echo "    protocol: TCP" >> config.yaml

.PHONY: deploy-local
deploy-local: create-config
	@if ! kind get clusters | grep -q '^kind$$'; then \
		kind create cluster --config=config.yaml; \
	fi; \

	@if ! kubectl get namespace | grep -q '^ingress-nginx$$'; then \
		kubectl apply -f https://raw.githubusercontent.com/kubernetes/ingress-nginx/main/deploy/static/provider/kind/deploy.yaml; \
	fi; \
	kubectl wait --namespace ingress-nginx --for=condition=ready pod --selector=app.kubernetes.io/component=controller --timeout=90s; \

	@if $(GITHUB_TOKEN) == "" || $(GITHUB_USER) == ""; then \
		echo "GITHUB_TOKEN and GITHUB_USER must be set"; \
		exit 1; \
	fi; \
	make docker-build -e GITHUB_TOKEN=$(GITHUB_TOKEN) GITHUB_USER=$(GITHUB_USER) TAG=$(TAG); \
	kind load docker-image $(TAG); \
	helm upgrade --install --create-namespace -n local policies-service helm-chart/application/ -f helm-chart/application/values-local.yaml

.PHONY: db-cleanup
db-cleanup:
	@if [ -f /tmp/pg-port-forward.pid ]; then \
		kill $$(cat /tmp/pg-port-forward.pid) || true; \
		rm /tmp/pg-port-forward.pid; \
	fi
	helm uninstall local-db || true
	kubectl delete pvc -l app.kubernetes.io/instance=local-db -n default || true

.PHONY: db-create
## The db-create target sets up a local PostgreSQL database in Kubernetes using Helm, waits for it to be ready,
## forwards the port, creates a test database, writes connection details to .local.env, and runs database migrations.
db-create:
	@kubectl config use-context kind-lhte >/dev/null
	@if ! helm list -A | grep -q '^local-db'; then \
		helm install local-db oci://registry-1.docker.io/bitnamicharts/postgresql --wait >/dev/null; \
	fi
	@kubectl wait --for=condition=ready pod -l app.kubernetes.io/name=postgresql -n default --timeout=120s >/dev/null
	@kubectl port-forward svc/local-db-postgresql 5555:5432 -n default >/tmp/pg-port-forward.log 2>&1 & \
		echo $$! > /tmp/pg-port-forward.pid; \
		sleep 5
	@export PGHOST=localhost; \
	export PGPORT=5555; \
	export PGUSER=postgres; \
	export PGPASSWORD=$$(kubectl get secret local-db-postgresql -n default -o jsonpath="{.data.postgres-password}" | base64 -d); \
	echo "export PGHOST=$$PGHOST"; \
	echo "export PGPORT=$$PGPORT"; \
	echo "export PGUSER=$$PGUSER"; \
	echo "export PGPASSWORD=$$PGPASSWORD"; \
	echo "export PGDATABASE=test"; \
	psql "host=$$PGHOST port=$$PGPORT user=$$PGUSER password=$$PGPASSWORD sslmode=disable" -c "CREATE DATABASE test;" 2>/dev/null || true; \
	echo "PGHOST=$$PGHOST" > ./.local.env; \
    echo "PGPORT=$$PGPORT" >> ./.local.env; \
    echo "PGUSER=$$PGUSER" >> ./.local.env; \
    echo "PGPASSWORD=$$PGPASSWORD" >> ./.local.env; \
    echo "PGDATABASE=test" >> ./.local.env; \
    make run ARGS="dbmigration ./.local.env"
