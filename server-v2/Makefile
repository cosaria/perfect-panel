.PHONY: contract contract-lint contract-bundle contract-openapi-ts

OPENAPI_DIR := openapi
OPENAPI_SPEC := ../server-v2/$(OPENAPI_DIR)/openapi.yaml
OPENAPI_REDOCLY := ../server-v2/$(OPENAPI_DIR)/redocly.yaml
OPENAPI_BUNDLE := ../server-v2/$(OPENAPI_DIR)/dist/openapi.json
OPENAPI_TS_DIR := ../server-v2/$(OPENAPI_DIR)/generated/ts
WEB_DIR := ../web
REDOCLY_BIN := $(WEB_DIR)/node_modules/.bin/redocly
OPENAPI_TS_BIN := $(WEB_DIR)/node_modules/.bin/openapi-ts

contract: contract-lint contract-bundle contract-openapi-ts

contract-lint:
	$(REDOCLY_BIN) lint --config $(OPENAPI_REDOCLY) $(OPENAPI_SPEC)

contract-bundle:
	mkdir -p $(dir $(OPENAPI_BUNDLE))
	$(REDOCLY_BIN) bundle $(OPENAPI_SPEC) --config $(OPENAPI_REDOCLY) --ext json -o $(OPENAPI_BUNDLE)

contract-openapi-ts:
	rm -rf $(OPENAPI_TS_DIR)
	mkdir -p $(dir $(OPENAPI_TS_DIR))
	$(OPENAPI_TS_BIN) -i $(OPENAPI_BUNDLE) -o $(OPENAPI_TS_DIR)
