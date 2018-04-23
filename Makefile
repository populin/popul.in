EXEC=docker exec -it 
EXEC_POPULIN=$(EXEC) populin-geography-api

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
	@$(EXEC_POPULIN) bash -c "go get -u github.com/alecthomas/gometalinter && go install github.com/alecthomas/gometalinter && gometalinter --install"
	@$(EXEC_POPULIN) bash -c "go install github.com/DATA-DOG/godog"

update:
	$(EXEC_POPULIN) bash -c "dep ensure -update"
	@$(EXEC_POPULIN) bash -c "go install github.com/DATA-DOG/godog"

bash: 
	$(EXEC_POPULIN) bash

run: 
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/cmd/geography && geography"
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/cmd/politics && politics"

doc: 
	@echo "documentation available on http://localhost:6060/pkg/github.com/populin/popul.in"
	@$(EXEC_POPULIN) bash -c "godoc -http=\":6060\""

test:
	@$(EXEC_POPULIN) bash -c "cd cmd/geography && go test -v -race"
	@$(EXEC_POPULIN) bash -c "cd cmd/politics && go test -v -race"

lint: 
	@$(EXEC_POPULIN) bash -c "gometalinter ./..."

fix:
	@$(EXEC_POPULIN) bash -c "gofmt -s -w ."
	@$(EXEC_POPULIN) bash -c "goimports -w ."

import-fixtures: 
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/cmd/geojson_importer && geojson_importer data/geography/fixtures"

import-data: 
	@$(EXEC_POPULIN) bash -c "go install github.com/populin/popul.in/cmd/geojson_importer && geojson_importer data/geography/real"

cloc:
	@cloc --exclude-list-file=.clocignore .
