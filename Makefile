all:
	go build "./..."
	go test "./..."
	staticcheck "./..."
	g++ -c -std=c++20 -fmodules-ts -xc++-system-header ctime \
		-xc++-system-header array \
		-xc++-system-header cerrno \
		-xc++-system-header cstdint \
		-xc++-system-header cstdlib \
		-xc++-system-header iostream \
		-xc++-system-header tuple \
		-xc++-system-header vector
	cat program.bpl | go run ./bin/main.go cpp
	clang-format -i a.bpl.cpp
	g++ -o main -std=c++20 -fmodules-ts c.cpp a.bpl.cpp
	cat c.bpl | go run ./bin/main.go query
