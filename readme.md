#  🚀 Tickitz – Golang + Gin Backend API----
 
>  A robust and performant RESTful API built using **Golang**, **Gin**, **Swagger**, and **Redis**. Designed for speed, scalability, and simplicity.

----

## 📸 Preview
> Swagger UI for Dokumentation: [`/swagger/index.html`](http://localhost:8080/swagger/index.html)

## 🚀 Features
- 🔐 JWT Authentication (Login & Register)
- 🧠 Redis Caching 
- 📘 Swagger Auto-Generated API Docs
- 🧾 CRUD for resources
- 🌐 RESTful API design
- 🗂️ MVC architecture
- 📦 PostgreSQL integration
- 🐳 Dockerized (Redis + PostgreSQL)
- 🧵 Graceful structured logging


## 🛠️ Tech Stack
![Go](https://img.shields.io/badge/-Go-00ADD8?logo=go&logoColor=white&style=for-the-badge)
![Gin](https://img.shields.io/badge/-Gin-00ADD8?logo=go&logoColor=white&style=for-the-badge)
![PostgreSQL](https://img.shields.io/badge/-PostgreSQL-4169E1?logo=postgresql&logoColor=white&style=for-the-badge)
![Docker](https://img.shields.io/badge/-Docker-2496ED?logo=docker&logoColor=white&style=for-the-badge)
![Swagger](https://img.shields.io/badge/Swagger-UI-85EA2D?logo=swagger&logoColor=black&style=for-the-badge)
![Redis](https://img.shields.io/badge/Redis-Database-DC382D?logo=redis&logoColor=white&style=for-the-badge)
![Postman](https://img.shields.io/badge/Postman-FF6C37?logo=postman&logoColor=white&style=for-the-badge)


##  🔐 .env Configuration
```
DBUSER=youruser
DBPASS=yourpass
DBHOST=localhost
DBPORT=yourport
DBNAME=tickitz

JWT_SECRET=your_jwt_secret

REDISUSER=youruser
REDISPASS=yourpass
REDISPORT=yourport

```

## 📦 How to Install & Run
First, clone this repository: 

https://github.com/federus1105/Backend-Golang-Weekly10.git
```bash
cd Backend-Golang-Weekly10
```
### Install Dependencies
```go
go mod tidy
```
### Run Project
```go
go run .\cmd\main.go 
```
### Run Swagger
```go
swag init -g ./cmd/main.go
```

## 📬 Postman Collection

You can try out the API using the Postman collection below:

🔗 [Tickitz API Postman Collection](https://federusrudi-9486783.postman.co/workspace/9cd45016-f25d-441e-8c5a-10f1070df09d/collection/48098195-225adccd-0cce-4652-9e86-4dd2ae598ae5?action=share&source=copy-link&creator=48098195)

> Make sure the server is running at `http://localhost:8080`


## 👨‍💻 Made by
### 📬 fedeursrudi@gmail.com
