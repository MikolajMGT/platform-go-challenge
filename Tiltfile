# -*- mode: Python -*-

load('ext://restart_process', 'docker_build_with_restart')

### CASSANDRA

k8s_yaml('assets/k8s/cassandra.yaml')
k8s_resource(
  'cassandra',
)

### ASSETS

compile_cmd = 'GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -tags netgo -v -o ./assets/build/assets-service cmd/main.go'

local_resource(
  'go-compile-assets',
  compile_cmd,
  deps=['./cmd', './internal', './pkg', './cfg'],
)

docker_build_with_restart(
  'assets-service',
  '.',
  entrypoint=['/usr/local/bin/app/assets-service'],
  dockerfile='Dockerfile.localbuild',
  only=[
    './assets/build/assets-service',
  ],
  live_update=[
    sync('assets/build/assets-service', '/usr/local/bin/app/assets-service'),
  ],
)

k8s_yaml('assets/k8s/assets.yaml')
k8s_resource(
  'assets',
  resource_deps=['go-compile-assets']
)
