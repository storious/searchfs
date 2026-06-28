.PHONY: fmt lint test build clean run

fmt:
	zig fmt build.zig src/*.zig

test:
	zig build test -Doptimize=Debug

build:
	zig build

run:
	zig build run

lint: fmt test build

clean:
	rm -rf .zig-cache zig-out
