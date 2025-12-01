### Flow Logic Pattern :
Router → Handler
Handler → Repository → Database

### Install swag
```bash
go install github.com/swaggo/swag/cmd/swag@latest
```

### How to run :
```bash
go run cmd/server/main.go
```

### Generate docs :
```bash
swag init -g cmd/server/main.go -o ./docs/
```

### Using Makefile (easy way)
If you have a `Makefile`, you can run:

```bash
make run
```
it will generate docs and run the app

```bash
make clean
```
it will migrate the database with clean option (drop all tables first)

```bash
make migrate
```
it will migrate the database

```bash
make seed
```
it will seed the database

```bash
make build
```
it will build the app

```bash
make test
```
it will run the test


## Links
- [Swagger](http://localhost:8080/swagger/index.html)
