.PHONY: contract contract-lint contract-bundle contract-openapi-ts

OPENAPI_DIR := openapi
OPENAPI_SPEC := $(OPENAPI_DIR)/openapi.yaml
OPENAPI_TS_INPUT := ./$(OPENAPI_SPEC)
OPENAPI_REDOCLY := $(OPENAPI_DIR)/redocly.yaml
OPENAPI_BUNDLE := $(OPENAPI_DIR)/dist/openapi.json
OPENAPI_TS_DIR := $(OPENAPI_DIR)/generated/ts
OPENAPI_TS_TMP := $(OPENAPI_DIR)/generated/openapi-ts-tmp

contract: contract-lint contract-bundle contract-openapi-ts

contract-lint:
	bunx @redocly/cli@2.25.4 lint --config $(OPENAPI_REDOCLY) $(OPENAPI_SPEC)

contract-bundle:
	mkdir -p $(dir $(OPENAPI_BUNDLE))
	bunx @redocly/cli@2.25.4 bundle $(OPENAPI_SPEC) --config $(OPENAPI_REDOCLY) --ext json -o $(OPENAPI_BUNDLE)

contract-openapi-ts:
	rm -rf $(OPENAPI_TS_DIR) $(OPENAPI_TS_TMP)
	mkdir -p $(dir $(OPENAPI_TS_TMP))
	bunx @hey-api/openapi-ts@0.95.0 -i $(OPENAPI_TS_INPUT) -o $(OPENAPI_TS_TMP)
	mv $(OPENAPI_TS_TMP) $(OPENAPI_TS_DIR)
