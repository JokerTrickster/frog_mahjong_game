shell:
	docker-compose run --service-ports app bash

build: clean
	docker-compose build

up:
	docker-compose up

down:
	docker-compose down --remove-orphans

clean: down
	find . -name '*.pyc' -exec rm -f {} +
	find . -name '*.pyo' -exec rm -f {} +
	find . -name '*~' -exec rm -f {} +
	find . -name '__pycache__' -exec rm -fr {} +

run-golang:
	docker-compose run --service-ports app ./app
