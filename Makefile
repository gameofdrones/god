god: god.go thunder.go
	go build  -o $@ $^

test: god
	sudo ./$<
