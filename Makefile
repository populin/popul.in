EXEC=docker exec -it 
EXEC_POPULIN=$(EXEC) populin-api 

.PHONY: build start stop restart install bash run run test lint fix import-fixtures import-data

build:
	@docker-compose build 

start: 
	@docker-compose up -d --remove-orphans --force-recreate

stop: 
	@docker-compose kill || true
	@docker-compose rm --force || true

restart: stop start

install:
	$(EXEC_POPULIN) bash -c "dep ensure -vendor-only -v"

bash: 
	$(EXEC_POPULIN) bash

run: 
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/cmd/api && api"

doc: 
	@echo "documentation available on http://localhost:6060/pkg/github.com/populin/popul.in"
	@$(EXEC_POPULIN) bash -c "godoc -http=\":6060\""

test: 
	@$(EXEC_POPULIN) bash -c "cd cmd/api && godog"

lint: 
	@$(EXEC_POPULIN) bash -c "gometalinter.v1 --config gometalinter.json ./..."

fix:
	@$(EXEC_POPULIN) bash -c "gofmt -s -w ."
	@$(EXEC_POPULIN) bash -c "goimports -w ."

import-fixtures: 
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/cmd/geojson_importer && geojson_importer data/geography/fixtures"

import-data: 
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/cmd/geojson_importer && geojson_importer data/geography/real"
