test:
	go test ./...

mock-gen:
	go install
	buf generate