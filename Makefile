.PHONY: build run test clean

build:
	docker-compose build

run:
	docker-compose up

stop:
	docker-compose down

clean:
	docker-compose down -v
	rm -rf service-a/tmp
	rm -rf service-b/tmp

