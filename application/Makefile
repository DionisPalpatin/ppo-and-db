export UNIT_SUCCESS=0
export INTEGRATION_SUCCESS=0

unit-test:
	rm -rf allure-unit/allure-results
	mkdir -p allure-unit/allure-results
	go test -shuffle on ./application/internal/business-logic/unit-tests ./application/internal/data-access/unit-tests --parallel 8 -tags=unit
	if [ $$? -eq 0 ]; then export UNIT_SUCCESS=1; else export UNIT_SUCCESS=0; fi
	cp environment.properties allure-unit/allure-results
	cp -r allure-unit/allure-results/* allure-results/

integration-test:
	rm -rf allure-integration/allure-results
	mkdir -p allure-integration/allure-results
	if [ "$$UNIT_SUCCESS" -ne 1 ]; then echo "Skipping integration tests due to failed unit tests." && export INTEGRATION_SUCCESS=0 && exit 0; fi
	go test -shuffle on ./application/integration-tests --parallel 8 -tags=integration
	if [ $$? -eq 0 ]; then export INTEGRATION_SUCCESS=1; else export INTEGRATION_SUCCESS=0; fi
	cp environment.properties allure-integration/allure-results
	cp -r allure-integration/allure-results/* allure-results/

e2e-test:
	rm -rf allure-e2e/allure-results
	mkdir -p allure-e2e/allure-results
	if [ "$$INTEGRATION_SUCCESS" -ne 1 ]; then echo "Skipping e2e tests due to failed integration tests." && exit 0; fi
	go test -shuffle on ./application/e2e-tests --parallel 8 -tags=e2e
	cp environment.properties allure-e2e/allure-results
	cp -r allure-e2e/allure-results/* allure-results/

allure:
	rm -rf allure-reports/
	allure generate allure-results -c -o allure-reports

get-report: unit-test integration-test e2e-test allure

show-report: unit-test integration-test e2e-test allure
	allure open allure-reports

.PHONY: all
all: get-report show-report
