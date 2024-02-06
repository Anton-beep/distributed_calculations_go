# Distributed Calculations
Distributed calculations written in Go language

# API Documentation
Generate documentation (swagger):
- [install swag](https://github.com/swaggo/swag)
- ```shell
swag init -g .\cmd\calculationServer\main.go -o .\docs\storage\
swag init -g .\cmd\storage\main.go -o .\docs\storage\
```