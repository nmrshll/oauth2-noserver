.PHONY: example
example:
	go run example/example.go

embed:
	embedmd -w README.md