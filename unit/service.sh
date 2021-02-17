#!/bin/bash
StartMySQL5.7() {
    docker pull mysql:5.7.25
    docker run --name=mysql5.7 -d -p 3308:3306 -e MYSQL_ROOT_PASSWORD=123456 mysql:5.7.25
    Waiting "MySQL 5.7" "ready for connections" 30
    docker exec mysql5.7 mysql -uroot -p123456 -e "CREATE DATABASE xun CHARACTER SET utf8 COLLATE utf8_general_ci"
    docker exec mysql5.7 mysql -uroot -p123456 -e "CREATE USER xun@'%' IDENTIFIED BY '123456'"
    docker exec mysql5.7 mysql -uroot -p123456 -e "GRANT SELECT ON xun.* TO 'xun'@'%'";
}

IsReady() {
    checkstr=$1
    res=$(docker logs mysql5.7 2>&1 | grep "$checkstr")
    if [ "$res" == "" ]; then
        echo "0"
    else
        echo "1"
    fi
}

Waiting() {
    name=$1
    checkstr=$2
    let timeout=$3
    echo -n "Starting $name ."
    isready=$(IsReady "$checkstr")
    timing=0
    while  [ "$isready" == "0" ];
    do
        sleep 1
        isready=$(IsReady "$checkstr")
        let timing=${timing}+1
        echo -n "."
        if [ $timing -eq $timeout ]; then
            echo " failed. timout($timeout)" >&2
            exit 1
        fi
    done
    echo "done"
}

command=$1
case $command in
    mysql5.7) StartMySQL5.7;;
    *) $(echo "please input command" >&2) ;;
esac