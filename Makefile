TARGET = $(CURDIR)/target

REPO ?= quay.io/$(USER)

override IMAGE := tryit-editor:latest

clean:
	@rm -rf $(TARGET)

build: clean
	@mkdir -p $(TARGET)
	@go build -o $(TARGET)

test: clean
	@go clean -testcache
	@go test ./...

test-coverage: clean
	@mkdir -p $(TARGET)
	@go test ./... -coverprofile=$(TARGET)/coverage.out
	go tool cover -html=$(TARGET)/coverage.out

image-build:
	@podman build --layers=true -t $(IMAGE) .

image-run:
	@podman run -it --rm -p 8080:8080 $(IMAGE)

image-push:
	@podman tag $(IMAGE) $(REPO)/$(IMAGE)
	@podman push $(REPO)/$(IMAGE)

openshift-deploy:
	@oc apply -f openshift-deployment.yaml
	@oc get route cool -n tryit-editor -o=jsonpath='{.spec.host}'
	@echo ""
