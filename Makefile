god: god.go thunder.go server.go
	go build  -o $@ $^

test: god
	sudo ./$<
