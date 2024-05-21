# GOPATH_DIR := $(GOPATH)/src/github.com/open-telemetry/opentelemetry-proto
GENDIR := gen
# GOPATH_GENDIR := $(GOPATH_DIR)/$(GENDIR)

PROTOC := protoc --proto_path=${PWD}/opentelemetry-proto

PROTO_FILES := $(wildcard opentelemetry-proto/opentelemetry/proto/collector/logs/v1/*.proto )


PROTO_GEN_GO_DIR ?= model/proto


define exec-command
$(1)

endef

PROTO_INCLUDES := \
	-Iopentelemetry-proto/opentelemetry/proto/common/v1 \
	-Iopentelemetry-proto/opentelemetry/proto/resource/v1 \
	-Iopentelemetry-proto/opentelemetry/proto/logs/v1 \
	-Iopentelemetry-proto/opentelemetry/proto/collector/logs/v1 
# Remapping of std types to gogo types (must not contain spaces)
PROTO_GOGO_MAPPINGS := $(shell echo \
		Mgoogle/protobuf/descriptor.proto=github.com/gogo/protobuf/types \
		Mgoogle/protobuf/timestamp.proto=github.com/gogo/protobuf/types \
		Mgoogle/protobuf/duration.proto=github.com/gogo/protobuf/types \
		Mgoogle/protobuf/empty.proto=github.com/gogo/protobuf/types \
		Mgoogle/api/annotations.proto=github.com/gogo/googleapis/google/api \
		Mcommon.proto=logger/model \
	| $(SED) 's/  */,/g')


# DO NOT DELETE EMPTY LINE at the end of the macro, it's required to separate commands.
define print_caption
  @echo "🏗️ "
  @echo "🏗️ " $1
  @echo "🏗️ "

endef

define proto_compile
  $(call print_caption, "Processing $(2) --> $(1)")

  $(PROTOC) \
    $(PROTO_INCLUDES) \
    --go_out=$(strip $(1)) \
	--go_opt=Mopentelemetry-proto/opentelemetry/proto/common/v1/common.proto=logger/model/proto \
    $(3) $(2)

endef

# $(call proto_compile, $(PROTO_GEN_GO_DIR), opentelemetry-proto/opentelemetry/proto/resource/v1/resource.proto)
# $(call proto_compile, $(PROTO_GEN_GO_DIR), opentelemetry-proto/opentelemetry/proto/common/v1/common.proto)
.PHONY: protoc
protoc:
	$(PROTOC) $(PROTO_INCLUDES) --go_out=model --go_opt=Mresource.proto=proto/resource/v1 resource.proto
	$(PROTOC) $(PROTO_INCLUDES) --go_out=model --go_opt=Mcommon.proto=proto/common/v1 common.proto
	$(PROTOC) $(PROTO_INCLUDES) --go_out=model --go_opt=Mlogs.proto=proto/logs/v1 logs.proto
	$(PROTOC) $(PROTO_INCLUDES) --go_out=model --go_opt=Mlogs_service.proto=proto/v1 logs_service.proto






# Generate gRPC/Protobuf implementation for Go.
.PHONY: proto-model
proto-model:
	rm -rf ./$(PROTO_GEN_GO_DIR)
	mkdir -p ./$(PROTO_GEN_GO_DIR)
	$(foreach file,$(PROTO_FILES),$(call exec-command,$(PROTOC) --go_out=./$(PROTO_GEN_GO_DIR) $(file)))

	
print: 
	@echo $(PROTO_FILES)



protoc -Iopentelemetry-proto/opentelemetry/proto/common/v1 --go_out=model/proto --go_opt=Mopentelemetry-proto/opentelemetry/proto/common/v1/common.proto=logger/model/proto opentelemetry-proto/opentelemetry/proto/common/v1/common.proto

# protoc -Iopentelemetry-proto/opentelemetry/proto/common/v1 --go_out=model --go_opt=Mcommon.proto=logger/model/proto common.proto


