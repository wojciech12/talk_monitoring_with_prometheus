SERVICE_NAME=order-manager
DOCKER_PROJECT_NAME=talk-golang-prom-red

docker_build:
	docker build -t wojciech12/$(DOCKER_PROJECT_NAME) .

start:
	docker-compose -f docker-compose.yaml --project-name $(DOCKER_PROJECT_NAME) up -d

stop:
	docker-compose -f docker-compose.yaml --project-name $(DOCKER_PROJECT_NAME) stop

rebuild:
	docker-compose -f docker-compose.yaml --project-name $(DOCKER_PROJECT_NAME) build

# reload the config from the disk
prometheus_reload_config:
	curl 127.0.0.1:9090/-/reload -X POST

srv_random_trafic:
	curl 127.0.0.1:8080/hello ; \
	curl 127.0.0.1:8080/world ;

srv_show_metrics:
	curl 127.0.0.1:8080/metrics
	
srv_random_trafic_complex_failed_db:
	curl 127.0.0.1:8080/complex?is_db_error=True

srv_random_trafic_complex_slow_db_and_svc:
	curl '127.0.0.1:8080/complex?db_sleep=3&srv_sleep=2'

srv_random_trafic_complex_failed_third_party:
	curl 127.0.0.1:8080/complex?is_srv_error=True

srv_wrk_random:
	 wrk -t8 -c12 -d360s -s tools/load_generator/mixed.lua http://127.0.0.1:8080

ensure_formatting:
	goimports -w *.go && goimports -w */*.go

utils_ln_gopath:
	mkdir -p $$GOPATH/src/github.com/wojciech12; \
	ln -s $$(pwd) $$GOPATH/src/github.com/wojciech12/$(SERVICE_NAME)