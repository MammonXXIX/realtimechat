load('ext://restart_process', 'docker_build_with_restart')

k8s_yaml('./infra/development/k8s/app-config.yaml')
k8s_yaml('./infra/development/k8s/api-gateway-deployment.yaml')
k8s_yaml('./infra/development/k8s/authentication-service-deployment.yaml')

docker_build_with_restart(
    'realtimechat/api-gateway',
    '.',
    entrypoint=['/app/build/api-gateway'],
    dockerfile='./infra/development/docker/api-gateway.Dockerfile',
    only=[
        './backend/services/api-gateway',
        './backend/shared',
        './backend/go.mod',
        './backend/go.sum',
    ],
    live_update=[
        sync('./backend/services/api-gateway', '/app/services/api-gateway'),
        sync('./backend/shared', '/app/shared'),
    ],
)
k8s_resource('api-gateway', port_forwards=8081, labels="services")

docker_build_with_restart(
    'realtimechat/authentication-service',
    '.',
    entrypoint=['/app/build/authentication-service'],
    dockerfile='./infra/development/docker/authentication-service.Dockerfile',
    only=[
        './backend/services/authentication-service',
        './backend/shared',
        './backend/go.mod',
        './backend/go.sum',
    ],
    live_update=[
        sync('./backend/services/authentication-service', '/app/services/authentication-service'),
        sync('./backend/shared', '/app/shared'),
    ],
)
k8s_resource('authentication-service', port_forwards=8082, labels="services")

