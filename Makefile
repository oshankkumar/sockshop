SHELL 		:= /bin/bash
BIN_DIR 	?= ./bin

GOFLAGS     := -v -trimpath
LDFLAGS     := -buildid= -extldflags "-f no-PIC -static"
TAGS        := osusergo netgo static_build

LINKCOLOR="\033[34;1m"
ENDCOLOR="\033[0m"
BINCOLOR="\033[37;1m"

define MAKE_GO_BUILD
	@printf '    %b %b...\n' $(LINKCOLOR)BUILDING$(ENDCOLOR) $(BINCOLOR)$(1)$(ENDCOLOR) 1>&2
	@CGO_ENABLED=0 go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/local/$(1) ./cmd/$(1)
	@printf '    %b Build complete.\n\n' $(LINKCOLOR)DONE$(ENDCOLOR)
endef

define MAKE_GO_BUILD_LINUX
	@printf '    %b %b...\n' $(LINKCOLOR)BUILDING$(ENDCOLOR) $(BINCOLOR)$(1)$(ENDCOLOR) 1>&2
	@GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(GOFLAGS) -tags '$(TAGS)' -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/linux-amd64/$(1) ./cmd/$(1)
	@printf '    %b Build complete.\n\n' $(LINKCOLOR)DONE$(ENDCOLOR)
endef

build:
	$(call MAKE_GO_BUILD,sockshop)