.PHONY: contract contract-lint contract-bundle contract-openapi-ts

OPENAPI_DIR := openapi
OPENAPI_SPEC := $(OPENAPI_DIR)/openapi.yaml
OPENAPI_REDOCLY := $(OPENAPI_DIR)/redocly.yaml
OPENAPI_BUNDLE := $(OPENAPI_DIR)/dist/openapi.json
OPENAPI_TS_DIR := $(OPENAPI_DIR)/generated/ts
OPENAPI_TS_TMP := $(OPENAPI_DIR)/generated/openapi-ts-tmp

contract: contract-lint contract-bundle contract-openapi-ts

contract-lint:
	bunx @redocly/cli lint --config $(OPENAPI_REDOCLY) $(OPENAPI_SPEC)

contract-bundle:
	mkdir -p $(dir $(OPENAPI_BUNDLE))
	bunx @redocly/cli bundle $(OPENAPI_SPEC) --config $(OPENAPI_REDOCLY) --ext json -o $(OPENAPI_BUNDLE)

contract-openapi-ts:
	rm -rf $(OPENAPI_TS_DIR) $(OPENAPI_TS_TMP)
	mkdir -p $(dir $(OPENAPI_TS_TMP))
	cd ../web && bunx @hey-api/openapi-ts@0.95.0 -i ../server-v2/$(OPENAPI_BUNDLE) -o ../server-v2/$(OPENAPI_TS_TMP)
	mv $(OPENAPI_TS_TMP) $(OPENAPI_TS_DIR)
