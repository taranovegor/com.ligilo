# com.ligilo
Сut and direct

## Set up
### Requirements
- [Docker](https://docs.docker.com/engine/install/)
- [docker-compose](https://docs.docker.com/compose/gettingstarted/)
- [MariaDB 11](https://mariadb.com/kb/en/mariadb-11-0-0-release-notes/)

### Configuration
Copy an instance of the environment file and save it as a file `.env`
```shell
cp .env.dist .env
```
... and configure as you need. Environment variables are described in comments

Build or pull docker images of application
```shell
make container-build
```
```shell
make container-pull
```

### Launch
Run the application using the Make tool
```shell
make start
```
