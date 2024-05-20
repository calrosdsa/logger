# GOPATH_DIR := $(GOPATH)/src/github.com/open-telemetry/opentelemetry-proto
GENDIR := gen
# GOPATH_GENDIR := $(GOPATH_DIR)/$(GENDIR)

PROTOC := protoc --proto_path=${PWD}

PROTO_FILES := $(wildcard opentelemetry-proto/opentelemetry/proto/logs/v1/*.proto )


PROTO_GEN_GO_DIR ?= model


define exec-command
$(1)

endef


# Generate gRPC/Protobuf implementation for Go.
.PHONY: gen-go
proto-model:
	rm -rf ./$(PROTO_GEN_GO_DIR)
	mkdir -p ./$(PROTO_GEN_GO_DIR)
	$(foreach file,$(PROTO_FILES),$(call exec-command,$(PROTOC) --go_out=./$(PROTO_GEN_GO_DIR) $(file)))

	
print: 
	@echo $(PROTO_FILES)