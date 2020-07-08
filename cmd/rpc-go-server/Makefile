gen-fixtures:
	@echo "==> generate fixtures/todo_no_types_gen.go.fix"
	@go run main.go -schema ../../examples/todo/schema.json > fixtures/todo_no_types_gen.go.fix
	@echo "==> generate fixtures/todo_types_gen.go.fix"
	@go run main.go -types github.com/apex/rpc/cmd/rpc-go-server/api -schema ../../examples/todo/schema.json > fixtures/todo_types_gen.go.fix

.PHONY: gen-fixtures