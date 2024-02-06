# Distributed Calculations
Distributed calculations written in Go language

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
