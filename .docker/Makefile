CI_COMMIT_SHORT_SHA ?= $(shell git rev-parse --short HEAD)

all: docker-base

.PHONY: docker-base
docker-base:
	docker build -f base.dockerfile -t base:${CI_COMMIT_SHORT_SHA} --progress=plain .
	docker image tag base:${CI_COMMIT_SHORT_SHA} registry-sd.vseinstrumenti.net/wtis/vetum/base:${CI_COMMIT_SHORT_SHA}
