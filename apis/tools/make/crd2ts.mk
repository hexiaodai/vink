crd2ts.generate: ## Generate OpenAPI specs
crd2ts.generate:
	@$(LOG_TARGET)
	@crdtoapi -i crd-to-openapi/ -o sdks/ts/openapi/openapi.yaml
	@npx openapi-typescript sdks/ts/openapi/openapi.yaml -o sdks/ts/openapi/openapi-schema.d.ts
	@sed -i '' 's/metadata\?: Record<string, never>/metadata\?: Record<string, any>/g' sdks/ts/openapi/openapi-schema.d.ts
