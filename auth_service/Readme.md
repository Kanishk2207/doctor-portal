mysql command:

docker run --network=ecofy-net --name my_mysql --restart always -e MYSQL_ROOT_PASSWORD=Kanishk_22 \
-v ~/ecofy-mysql:/var/lib/mysql \
-p 127.0.0.1:3307:3306 -d mysql:8.0


postgres command:

docker run -d  --network=doctor-portal --name my_postgres   -e POSTGRES_USER=dev_user   -e POSTGRES_PASSWORD=Kanishk_22   -e POSTGRES_DB=mydatabase   -v ~/doctor-portal-psql:/var/lib/postgresql/data   -p 127.0.0.1:5433:5432   postgres:latest


nodemon command:

nodemon --signal SIGTERM


docker run --name local_mysql --restart always -e MYSQL_ROOT_PASSWORD=Kanishk_22 \
-p 127.0.0.1:3307:3306 -d mysql:8.0

