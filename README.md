Simple Go Gin Starter include :
- JWT Auth
- Postgres Database
- Swagger

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

## How to run with docker-compose
Pastikan file .env sudah sesuai:
bash
cp .env.example .env
configurasi file .env, lihat .env.example

# Pastikan konfigurasi database di .env sesuai dengan yang di docker-compose.yaml

Build dan jalankan aplikasi:
bash
# Build image dan start container
docker-compose up -d

# Lihat log migration
docker-compose logs -f migrate
Atau jalankan migrasi secara terpisah:
bash
# Hanya jalankan migrasi
docker-compose run --rm migrate
Setelah migrasi selesai, aplikasi akan berjalan di:
http://localhost:8080


# How to create admin user :
- create in swagger with master password, if not set master password, it will create user with role `user`
http://localhost:8080/swagger/index.html#/auth/post_auth_register
