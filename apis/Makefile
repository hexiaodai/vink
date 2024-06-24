
ifeq ($(shell uname),Darwin)
SEDI=sed -i ""
else
SEDI=sed -i
endif

.PHONY: gen-proto
gen-proto: clean-proto
	./gen.sh
	@for f in $$(find . -name "*.pb.go"); do \
      if grep github.com/golang/protobuf/ptypes/struct $$f &>/dev/null; then \
          $(SEDI) 's#github.com/golang/protobuf/ptypes/struct#google.golang.org/protobuf/types/known/structpb#g' $$f; \
      fi; \
    done

.PHONY: clean-proto
clean-proto:
	./clean.sh

.PHONY: gen-labels
gen-labels:
	go run ./tools/annotations_prep/main.go --input ./label/labels.yaml --output ./label/labels.gen.go --collection_type label

.PHONY: gen-annotations
gen-annotations:
	go run ./tools/annotations_prep/main.go --input ./annotation/annotations.yaml --output ./annotation/annotations.gen.go --collection_type annotation

gen: gen-labels gen-annotations gen-proto

.PHONY: gen
