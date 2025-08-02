go_build := CGO_ENABLED=0 go build
binaries_path := bin

_pre_build:
	mkdir -p  $(binaries_path)

#build: @ Build pokedex
build: _pre_build
	$(go_build)	-o $(binaries_path)/pokedex ./cmd/pokedex/main.go