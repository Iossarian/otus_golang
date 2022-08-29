docker-compose -f deployments/docker-compose.tests.yaml up --build
EXIT_CODE=$?
docker-compose -f deployments/docker-compose.tests.yaml down
docker-compose --env-file deployments/.env -f deployments/docker-compose.yaml down
exit ${EXIT_CODE}