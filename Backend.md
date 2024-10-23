# Backend

<h2>📦 Packages</h2>

```bash
go get github.com/labstack/echo/v4
go get github.com/labstack/echo/v4/middleware
go get github.com/go-playground/validator/v10
go get github.com/joho/godotenv
go get go.mongodb.org/mongo-driver/mongo
go get github.com/golang-jwt/jwt/v5
go get github.com/stretchr/testify
```

<h2>📃 Start App in Terminal</h2>

```bash
go run main.go ./env/dev/.env.auth
```
```bash
go run main.go ./env/dev/.env.booking
```
```bash
go run main.go ./env/dev/.env.user
```
```bash
go run main.go ./env/dev/.env.gym
```
```bash
go run main.go ./env/dev/.env.football
```
```bash
go run main.go ./env/dev/.env.swimming
```
```bash
go run main.go ./env/dev/.env.badminton
```
```bash
go run main.go ./env/dev/.env.payment
```
```bash
go run main.go ./env/dev/.env.notification
```

<p>Migration</p>

<p>dev</p>

```bash
go run ./pkg/database/script/migration.go ./env/dev/.env.user && \
go run ./pkg/database/script/migration.go ./env/dev/.env.auth && \
go run ./pkg/database/script/migration.go ./env/dev/.env.booking && \
go run ./pkg/database/script/migration.go ./env/dev/.env.gym && \
go run ./pkg/database/script/migration.go ./env/dev/.env.swimming && \
go run ./pkg/database/script/migration.go ./env/dev/.env.badminton && \
go run ./pkg/database/script/migration.go ./env/dev/.env.football && \
go run ./pkg/database/script/migration.go ./env/dev/.env.payment
```

<h2>🍰 Generate a Proto File Command</h2>
<p>User</p>

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    modules/player/proto/userPb.proto
```

<p>Auth</p>

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    modules/auth/proto/authPb.proto
```

<p>Gym</p>

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    modules/booking/proto/gymPb.proto
```

<p>Swimming</p>

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    modules/swimming/proto/swimmingPb.proto
```

<p>Badminton</p>

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    modules/badminton/proto/badmintonPb.proto
```

<p>Football</p>

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    modules/football/proto/footballPb.proto
```


<h2>🐳 Docker Build</h2>

```bash
docker build -t 6531503042/Sport-Complexp:latest -f build/auth/Dockerfile .
```