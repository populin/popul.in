EXEC=docker exec -it 
EXEC_POPULIN=$(EXEC) populin-api 

build:
	@docker-compose build 

start: 
	@docker-compose up -d --remove-orphans --force-recreate

stop: 
	@docker-compose kill || true
	@docker-compose rm --force || true

restart: stop start

install:
	$(EXEC_POPULIN) bash -c "dep ensure -vendor-only"

bash: 
	$(EXEC_POPULIN) bash

run: 
	@$(EXEC_POPULIN) bash -c "go build && ./popul.in"

test: 
	@$(EXEC_POPULIN) bash -c "godog"

lint: 
	@$(EXEC_POPULIN) bash -c "gometalinter.v1 --config gometalinter.json ./..."

fix:
	@$(EXEC_POPULIN) bash -c "gofmt -s -w ."
	@$(EXEC_POPULIN) bash -c "goimports -w ."

import-fixtures: 
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/importer && importer data/geography/fixtures"

import-data: 
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/importer && importer data/geography/real"
