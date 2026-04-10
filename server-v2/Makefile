.PHONY: contract contract-lint contract-bundle contract-openapi-ts

OPENAPI_DIR := openapi
OPENAPI_SPEC := ../server-v2/$(OPENAPI_DIR)/openapi.yaml
OPENAPI_REDOCLY := ../server-v2/$(OPENAPI_DIR)/redocly.yaml
OPENAPI_BUNDLE := ../server-v2/$(OPENAPI_DIR)/dist/openapi.json
OPENAPI_TS_INPUT := ../server-v2/$(OPENAPI_DIR)/openapi.yaml
OPENAPI_TS_DIR := ../server-v2/$(OPENAPI_DIR)/generated/ts
OPENAPI_TS_TMP := ../server-v2/$(OPENAPI_DIR)/generated/openapi-ts-tmp
WEB_DIR := ../web

contract: contract-lint contract-bundle contract-openapi-ts

contract-lint:
	cd $(WEB_DIR) && bunx @redocly/cli@2.25.4 lint --config $(OPENAPI_REDOCLY) $(OPENAPI_SPEC)

contract-bundle:
	mkdir -p $(dir $(OPENAPI_BUNDLE))
	cd $(WEB_DIR) && bunx @redocly/cli@2.25.4 bundle $(OPENAPI_SPEC) --config $(OPENAPI_REDOCLY) --ext json -o $(OPENAPI_BUNDLE)

contract-openapi-ts:
	rm -rf $(OPENAPI_TS_DIR) $(OPENAPI_TS_TMP)
	mkdir -p $(dir $(OPENAPI_TS_TMP))
	cd $(WEB_DIR) && bunx @hey-api/openapi-ts@0.95.0 -i $(OPENAPI_TS_INPUT) -o $(OPENAPI_TS_TMP)
	mv $(OPENAPI_TS_TMP) $(OPENAPI_TS_DIR)
