.PHONY: \
	all check \
	build test fmt lint clean \
	searchfs-build searchfs-test searchfs-fmt searchfs-lint \
	gdfs-build gdfs-test gdfs-fmt gdfs-lint \
	zigkv-build zigkv-test zigkv-fmt zigkv-lint

all: build

check: fmt lint test

#
# All projects
#

build: searchfs-build gdfs-build zigkv-build

test: searchfs-test gdfs-test zigkv-test

fmt: searchfs-fmt gdfs-fmt zigkv-fmt

lint: searchfs-lint gdfs-lint zigkv-lint

clean:
	$(MAKE) -C searchfs clean
	$(MAKE) -C gdfs clean
	$(MAKE) -C zigkv clean

#
# SearchFS
#

searchfs-build:
	$(MAKE) -C searchfs build

searchfs-test:
	$(MAKE) -C searchfs test

searchfs-fmt:
	$(MAKE) -C searchfs fmt

searchfs-lint:
	$(MAKE) -C searchfs lint

#
# GDFS
#

gdfs-build:
	$(MAKE) -C gdfs build

gdfs-test:
	$(MAKE) -C gdfs test

gdfs-fmt:
	$(MAKE) -C gdfs fmt

gdfs-lint:
	$(MAKE) -C gdfs lint

#
# ZigKV
#

zigkv-build:
	$(MAKE) -C zigkv build

zigkv-test:
	$(MAKE) -C zigkv test

zigkv-fmt:
	$(MAKE) -C zigkv fmt

zigkv-lint:
	$(MAKE) -C zigkv lint

