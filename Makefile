god: god.go context.go
	go build  -o $@ $^

test: god
	./$<
