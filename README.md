## Architecture Folder structure

scalable-backend-architecture/
├── cmd/
│   └── api/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── database/
│   │   └── database.go
│   ├── models/
│   │   └── base_model.go
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── cors.go
│   │   ├── logger.go
│   │   └── error.go
│   ├── user/
│   │   ├── model.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   └── dto.go
│   └── router/
│       └── router.go
├── pkg/
│   ├── response/
│   │   └── response.go
│   └── utils/
│       ├── hash.go
│       └── token.go
├── .env.example
├── .env
├── .gitignore
├── go.mod
├── Makefile
└── README.md