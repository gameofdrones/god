god: god.go thunder.go server.go motion.go
	go build  -o $@ $^

test: god
	sudo ./$<
