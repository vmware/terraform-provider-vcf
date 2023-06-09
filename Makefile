TEST?=$$(go list ./... |grep -v 'vendor')
PKG_NAME=internal

documentation:
	go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --examples-dir=./examples
	
fmtcheck:
	@sh -c "'$(CURDIR)/scripts/gofmtcheck.sh'"

lint:
	@echo "==> Checking source code against linters..."
	@golangci-lint run ./$(PKG_NAME)/...

build: fmtcheck
	go install

init:
	go build -o terraform-provider-vcf
	terraform init

test: fmtcheck
	go test -i $(TEST) || exit 1
	echo $(TEST) | \
		xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 240m -parallel=4
