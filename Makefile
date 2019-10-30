run:
	@gofmt -l -w .
	@go run main.go

package:
	@go build main.go

