#!/bin/bash
source config.sh

# Check if Postgres container is running
if [ "$(docker ps -q -f name=my_postgres)" ]; then
    echo "Postgres container is already running."
else
    echo "Starting Postgres container..."
    docker run -d  --network=doctor-portal --name my_postgres \
    -e POSTGRES_USER=dev_user -e POSTGRES_PASSWORD=Kanishk_22 \
    -e POSTGRES_DB=mydatabase -v ~/doctor-portal-psql:/var/lib/postgresql/data \
    -p 127.0.0.1:5433:5432 postgres:latest

    # Wait for Postgres to be ready
    echo "Waiting for Postgres to start..."
    sleep 10

    # Create authservicedb and patientservicedb
    echo "Creating authservicedb and patientservicedb..."
    docker exec -i my_postgres psql -U dev_user -d mydatabase <<EOF
    CREATE DATABASE authservicedb;
    CREATE DATABASE patientservicedb;
EOF
fi

# Check if the service container is running
if [ "$(docker ps -q -f name=$CONTAINER_NAME)" ]; then
    echo "Stopping and removing the existing container..."
    docker stop $CONTAINER_NAME
    docker rm $CONTAINER_NAME
fi

# Check if the service container exists but is not running
if [ "$(docker ps -aq -f status=exited -f name=$CONTAINER_NAME)" ]; then
    echo "Removing the existing stopped container..."
    docker rm $CONTAINER_NAME
fi

# Run the service container
echo "Starting auth service"
docker run -it --network=doctor-portal -p 8081:8080 -p 50052:50051 --restart always \
--name $CONTAINER_NAME -v $DIR:/app \
--env-file .env $IMAGE_NAME:$IMAGE_TAG sh -c "cd cmd/auth_service/ && bash"
