default: vet test

vet:
	echo "==> Running vet..."
	go vet `go list ./... | grep -v /vendor/`

test:
	echo "==> Running tests..."
	go test -cover -v `go list ./... | grep -v /vendor/`
