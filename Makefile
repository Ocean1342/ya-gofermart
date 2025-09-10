include .env
define goose_env_export
	export GOOSE_DBSTRING="host=0.0.0.0 port=5432 user=ya password=ya database=ya  sslmode=disable"
	export GOOSE_DRIVER=postgres
endef

migration_create:
	cd ./migrations && \
	$(GOPATH)/bin/goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) create table sql

migration_up:
	cd /home/ant/go-work/go-musthave-diploma-tpl/migrations && \
	$(GOPATH)/bin/goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) up

migration_drop_all:
	cd /home/ant/go-work/go-musthave-diploma-tpl/migrations && \
	$(GOPATH)/bin/goose $(GOOSE_DRIVER) $(GOOSE_DBSTRING) down-to 0

mockgen_accrual_system:
	$(GOPATH)/bin/mockgen -source=./internal/accrual-system/inteface.go  -destination=./internal/accrual-system/system_mock.go -package=accrual_system
