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

build-linux:
	$(call MAKE_GO_BUILD_LINUX,sockshop)

generate-dep-graph:
	godepgraph -i github.com/go-chi/chi/v5/middleware,github.com/go-chi/chi/v5,github.com/go-sql-driver/mysql,github.com/google/uuid,github.com/jmoiron/sqlx,go.uber.org/zap,github.com/prometheus/client_golang/prometheus,github.com/prometheus/client_golang/prometheus/promauto,github.com/prometheus/client_golang/prometheus/promhttp  -s ./cmd/sockshop | dot -Tpng -o godepgraph.png

run:
	docker compose -f docker-compose.yml up -d --wait

clean:
	docker compose down