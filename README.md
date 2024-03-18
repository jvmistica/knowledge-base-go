# knowledge-base-go
![Build status](https://github.com/jvmistica/knowledge-base-go/workflows/build/badge.svg)
A collection of knowledge base APIs.

## Environment Variables
```
export POSTGRES_HOST=<postgres_host>
export POSTGRES_PORT=<postgres_port>
export POSTGRES_USER=<user>
export POSTGRES_PASS=<password>
export POSTGRES_DB=<postgres_db>
```

## Usage
To run without seeding the database:
```
go run main.go
```

To run and seed the database:
```
go run main.go --seed=true
```
