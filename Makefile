### UTILS

tidy-all::
	go mod tidy

create-configs-local:
	kubectl --namespace dev delete configmap local-config 2>/dev/null; true
	kubectl --namespace dev create configmap local-config \
	--from-file=config.yaml

### LOCAL ENV

init-local-env::
	eval assets/scripts/init-local.sh

start-local-env: tidy-all create-configs-local
	tilt up

stop-local-env:
	tilt down

### RUNTIME

start:
	docker compose up

test:
	go test ./...

stop:
	docker compose down
