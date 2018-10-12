build:
	docker build . -t rproxy
	docker build . -f ./test/Dockerfile -t rproxy_test

run: build
	docker rm -f rproxy_redis_1 || true
	docker rm -f rproxy_rproxy_1 || true

	docker-compose up -d redis
	docker-compose up -d rproxy

test: run
	docker-compose up test