fmt:
	gofmt -w .

vet:
	go vet ./...

test:
	go test ./...

check: fmt vet test
