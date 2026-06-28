# Golang-Fiber-Tunnel-Wifi-Project

## Running Guide

**1. Install Depedencies**
```bash
go mod tidy 
```
or 
```bash
go mod download
```

**2. Run Migration**
```bash
go run cmd/migratedatabase/main.go
```

**3. Run Sedder**
```bash
go run cmd/runseeder/main.go
```

**4. Chane file name .env.example to .env and set up**

**5. Run Apps**
```bash
go run main.go
```
or with docker-compose
```bash
    docker-compose up --build
```