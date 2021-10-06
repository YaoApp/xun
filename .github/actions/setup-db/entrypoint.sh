#!/bin/sh

docker_run="docker run"

startMySQL() {
    VERSION=$1
    echo "Start MySQL $VERSION"
    docker_run="$docker_run -e MYSQL_RANDOM_ROOT_PASSWORD=true -e MYSQL_USER=$INPUT_USER -e MYSQL_PASSWORD=$INPUT_PASSWORD"
    docker_run="$docker_run -e MYSQL_DATABASE=$INPUT_DB"
    docker_run="$docker_run -d -p 3306:3306 mysql:$VERSION --port=3306"
    docker_run="$docker_run --character-set-server=utf8mb4 --collation-server=utf8mb4_general_ci"
    sh -c "$docker_run"

    DB_HOST="tcp(127.0.0.1:3306)/$INPUT_DB?charset=utf8mb4&parseTime=True&loc=Local"
    DB_USER=$INPUT_USER
    echo "DB_HOST=$DB_HOST" >> $GITHUB_ENV
    echo "DB_USER=$DB_USER" >> $GITHUB_ENV
    echo "DB_DRIVER=mysql" >> $GITHUB_ENV
    echo "$DB_HOST"
}

startPostgres() {
    VERSION=$1
    echo "Start Postgres $VERSION"
    docker_run="$docker_run --name postgres_$VERSION"
    docker_run="$docker_run -e POSTGRES_DB=$INPUT_DB"
    docker_run="$docker_run -e POSTGRES_USER=$INPUT_USER"
    docker_run="$docker_run -e POSTGRES_PASSWORD=$INPUT_PASSWORD"
    docker_run="$docker_run -d -p 5432:5432 postgres:$VERSION"

    # waiting for postgres ready
    timeout 90s bash -c "until docker exec postgres_$VERSION pg_isready ; do sleep 5 ; done"

    DB_HOST="127.0.0.1/$INPUT_DB?sslmode=disable"
    DB_USER=$INPUT_USER
    echo "DB_HOST=$DB_HOST" >> $GITHUB_ENV
    echo "DB_USER=$DB_USER" >> $GITHUB_ENV
    echo "DB_DRIVER=postgres" >> $GITHUB_ENV
    echo "$DB_HOST"
}

startSQLite3() {
    echo "Start SQLite3"
    echo "DSN=file:$INPUT_DB.db?cache=shared&mode=memory" >> $GITHUB_ENV
    echo "DB_DRIVER=sqlite3" >> $GITHUB_ENV
}

# MySQL8.0, MySQL5.7, Postgres9.6, Postgres14, SQLite3
case $INPUT_KIND  in 
MySQL8.0)
    startMySQL 8.0
    ;;
MySQL5.7)
    startMySQL 5.7
    ;;
MySQL5.6)
    startMySQL 5.6
    ;;
Postgres9.6)
    startPostgres 9.6
    ;;
Postgres14)
    startPostgres 14
    ;;
SQLite3)
    startSQLite3 
    ;;
esac
