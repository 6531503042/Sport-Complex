sudo systemctl start mongod

go run main.go ./env/dev/.env.auth

go run main.go ./env/dev/.env.booking

go run main.go ./env/dev/.env.user

go run main.go ./env/dev/.env.facility

go run main.go ./env/dev/.env.payment