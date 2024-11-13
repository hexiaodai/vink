##@ grpc

.PHONY: grpc.generate
grpc.generate: ## Generated client and server code.
grpc.generate: grpc.clean
	@$(LOG_TARGET)
	buf generate --timeout 10m -v \
	--path types/ \
	--path management/

	@for d in types/ management/; do \
		for f in $$(find $$d -name "*.proto"); do \
			protoc --validate_out="paths=source_relative,lang=go:." \
			$$f; \
		done \
	done

	@for d in types/ management/; do \
		for f in $$(find $$d -name "*.proto"); do \
			npx protoc --ts_out $(ROOT_DIR)/sdks/ts \
			$$f; \
		done \
	done

PATTERNS := .validate.go _deepcopy.gen.go _grpc.pb.go .gen.json gr.gen.go .pb.go _json.gen.go .pb.gw.go .swagger.json .deepcopy.go

.PHONY: grpc.clean
grpc.clean: ## Clean generated code.
	@$(LOG_TARGET)
	@for p in $(PATTERNS); do \
    	rm -f $(ROOT_DIR)/**/**/**/*"$$p"; \
    	rm -f $(ROOT_DIR)/**/**/*"$$p"; \
    	rm -f $(ROOT_DIR)/**/*"$$p"; \
	done

	@find $(ROOT_DIR)/sdks/ts \
	| grep -v package.json \
	| grep -v label \
	| grep -v annotation \
	| awk "NR != 1" \
	| xargs rm -rf

.PHONY: grpc.release.ts-sdk
grpc.release.ts-sdk: ## Release the js sdk.
	@$(LOG_TARGET)
	if [ -z "$(VERSION)" ]; then echo "VERSION is not set"; exit 1; fi
	$(SEDI) "s/\"version\": \".*\"/\"version\": \"$(VERSION)\"/g" $(ROOT_DIR)/sdks/ts/package.json
	npm --version
	node --version
	cd sdks/ts; \
	npm config set //registry.npmjs.org/:_authToken $(NPM_TOKEN) && \
	npm publish --access=public

.PHONY: grpc.release
grpc.release: ## Release the grpc code.
grpc.release: grpc.release.ts-sdk
