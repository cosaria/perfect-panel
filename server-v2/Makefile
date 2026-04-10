.PHONY: contract contract-lint contract-bundle contract-openapi-ts

OPENAPI_DIR := openapi
OPENAPI_SPEC := ../server-v2/$(OPENAPI_DIR)/openapi.yaml
OPENAPI_REDOCLY := ../server-v2/$(OPENAPI_DIR)/redocly.yaml
OPENAPI_BUNDLE := ../server-v2/$(OPENAPI_DIR)/dist/openapi.json
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
	tmp_bundle="$$(mktemp /tmp/server-v2-openapi-bundle.XXXXXX.json)"; \
	trap 'rm -f "$$tmp_bundle"' EXIT; \
	node -e 'const fs=require("fs"); const input="$(OPENAPI_BUNDLE)"; const output=process.argv[1]; const doc=JSON.parse(fs.readFileSync(input,"utf8")); const refMap=new Map([["#/components/parameters/idempotencyKey","#/components/parameters/idempotency_key"],["#/components/parameters/pageSize","#/components/parameters/page_size"],["#/components/parameters/nodeTimestamp","#/components/parameters/node_timestamp"],["#/components/parameters/nodeNonce","#/components/parameters/node_nonce"],["#/components/parameters/nodeSignature","#/components/parameters/node_signature"],["#/components/schemas/Problem","#/components/schemas/problem"],["#/components/schemas/ValidationProblem","#/components/schemas/validation_problem"],["#/components/schemas/PaginationMeta","#/components/schemas/pagination_meta"],["#/components/schemas/Money","#/components/schemas/money"]]); const walk=(value)=>{ if (Array.isArray(value)) return value.map(walk); if (value && typeof value === "object") { if (typeof value.$$ref === "string" && refMap.has(value.$$ref)) return { ...value, $$ref: refMap.get(value.$$ref) }; const out={}; for (const [key, nested] of Object.entries(value)) out[key]=walk(nested); return out; } return value; }; const sanitized=walk(doc); for (const [group, keys] of Object.entries({parameters:["idempotencyKey","pageSize","nodeTimestamp","nodeNonce","nodeSignature"],schemas:["Problem","ValidationProblem","PaginationMeta","Money"]})) { const bucket=sanitized.components?.[group]; if (!bucket) continue; for (const key of keys) delete bucket[key]; } fs.writeFileSync(output, JSON.stringify(sanitized, null, 2));' "$$tmp_bundle"; \
	cd $(WEB_DIR) && bunx @hey-api/openapi-ts@0.95.0 -i "$$tmp_bundle" -o $(OPENAPI_TS_TMP); \
	mv $(OPENAPI_TS_TMP) $(OPENAPI_TS_DIR)
