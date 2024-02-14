# Distributed Calculations
Distributed calculations written in Go language

# Configure
***Before building and using calculation server and storage you must create `.env` file*** (specifies environmental variables, i.e. config). To do this you can use `.env.template` (`calculationServer/.env.template`, `storage/.env.template`) and just put your values there.
## CalculationServer
- `STORAGE_URL` - URL of storage server
- `NUMBER_OF_CALCULATORS` - Number of calculators (workers) that will be created
- `SEND_ALIVE_DURATION` - Duration of sending alive message to storage server

## Storage
- `POSTGRESQL_USER` - User for database
- `POSTGRESQL_PASSWORD` - Password for database
- `POSTGRESQL_DB` - Database name
- `POSTGRESQL_HOST` - Host of database
- `POSTGRESQL_PORT` - Port of database
- `RESET_POSTGRESQL` - If `TRUE` then database will be reset on start of the storage server
- `CHECK_SERVER_DURATION` - Duration of checking if calculation server is alive

# Database Start
***Docker is required! ([install](https://docs.docker.com/engine/install/))***

```shell
docker run --name <name for docker container> -p 5432:5432 -e  POSTGRES_USER=<your database user> -e POSTGRES_PASSWORD=<your password for database> -d postgres:16
```

You can also specify local directory for database storage using: `-v <local path>:/var/lib/postgresql/data postgres:16`

*Based on https://hub.docker.com/_/postgres*

You can also start docker somehow else.

When docker is running, you need to reset it, before program can use it, so on the first start of `storage/main.go` set in your `storage/.env` `RESET_POSTGRESQL=TRUE`. After the first start, do not forget to set `RESET_POSTGREQL` to `FALSE`, or it will wipe data.

# Build
**Storage:**
```shell
cd storage
go build .
```

**Calculation Server:**
```shell
cd calculationServer
go build .
```

# API Documentation
Documentation is always available (without rebuild of the documentation) available at http://localhost:8080/swagger/index.html
Generate documentation (swagger):
[install swag](https://github.com/swaggo/swag)
````shell
cd storage
swag fmt
swag init
````

# Tests
For storage testing database is required (see **Database Start** section), also do not forget to change `calculationServer/tests/config_test.go` and `storage/tests/config_test.go` to specify where is postgresql database, number of calculators, and secret key.\
To run tests:
````shell
cd calculationServer
go test -v ./tests/...
cd ..
cd storage
go test -v ./tests/...
````