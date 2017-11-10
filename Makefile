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

dep-install:
	$(EXEC_POPULIN) bash -c "dep ensure -vendor-only"

bash-geography: 
	$(EXEC_POPULIN) bash

run-geography: 
	@$(EXEC_POPULIN) bash -c "go build && ./popul.in"

test-geography: 
	@$(EXEC_POPULIN) bash -c "godog"

lint-geography: 
	@$(EXEC_POPULIN) bash -c "gometalinter.v1 --config gometalinter.json ./..."

import-geography-fixtures: 
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/importer && importer data/geography/fixtures"

import-geography-data: 
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/importer && importer data/geography/real"
