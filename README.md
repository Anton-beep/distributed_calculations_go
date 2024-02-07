# Distributed Calculations
Distributed calculations written in Go language

# Database Start
***Docker is required! ([install](https://docs.docker.com/engine/install/))***

```shell
docker run --name <name for database> -p 5432:5432 -e  POSTGRES_USER=<your database user> -e POSTGRES_PASSWORD=<your password for database> -d postgres:16
```

You can also specify local directory for database storage using: `-v <local path>:/var/lib/postgresql/data postgres:16`

*Based on https://hub.docker.com/_/postgres*

You can also start docker somehow else, it must be on http://localhost:5432
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
Generate documentation (swagger):
[install swag](https://github.com/swaggo/swag)
````shell
cd calculationServer
swag fmt 
swag init
cd ..
cd storage
swag fmt
swag init
````
Documentation is available at http://localhost:8080/swagger/index.html

# Tests And Benchmarks
For storage testing database is required (see **Database Start** section), also do not forget to change `calculationServer/tests/config_test.go` and `storage/tests/config_test.go` to specify where is postgresql database, number of calculators, and secret key.\
To run tests:
````shell
cd calculationServer
go test -v ./tests/...
cd ..
cd storage
go test -v ./tests/...
````