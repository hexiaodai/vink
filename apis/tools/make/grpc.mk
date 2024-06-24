##@ grpc

.PHONY: grpc.generate
grpc.generate: ## Generated client and server code.
grpc.generate: grpc.clean
	@$(LOG_TARGET)
	buf generate --timeout 10m -v \
	--path common/ \
	--path management/

	@for d in common/ management/; do \
		for f in $$(find $$d -name "*.proto"); do \
			protoc --validate_out="paths=source_relative,lang=go:." $$f; \
		done \
	done

	@for d in sdks/ts/management; do \
		for f in $$(find $$d -type f -name "*.ts"); do \
    		if [ "$$(uname)" = "Darwin" ]; then \
      			$(SEDI) -r 's#(^type Base.*)#/* vink modified */ export \1#g' $$f; \
    		else \
      			$(SEDI) -r 's#(^type Base.*)#/* vink modified */ export \1#g' $$f; \
    		fi \
  		done \
	done

PATTERNS := .validate.go _deepcopy.gen.go .gen.json gr.gen.go .pb.go _json.gen.go .pb.gw.go .swagger.json .deepcopy.go

.PHONY: grpc.clean
grpc.clean: ## Clean generated code.
	@$(LOG_TARGET)
	@for p in $(PATTERNS); do \
    	rm -f $(ROOT_DIR)/**/**/**/*"$$p"; \
    	rm -f $(ROOT_DIR)/**/**/*"$$p"; \
    	rm -f $(ROOT_DIR)/**/*"$$p"; \
	done

	@find $(ROOT_DIR)/sdks/ts | grep -v  package.json | awk "NR != 1" | xargs rm -rf
