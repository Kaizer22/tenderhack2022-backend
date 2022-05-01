# strategy-management-service
### (Made for TenderHack2022)

Microservice to run users automated trading strategies.
Functionality:
- Basic Logging
- Auto generating Swagger docs
- Business logic
- Configuring using envs
- Profiling with pprof
- Usage of go-pg ORM
- JWT authentication & authorization
- Access control based on policy descriptor `basic_policy.csv`


To get the Swagger page go to: `/swagger/index.html`

| Var name                    | Var description                                                                                                              | Default value |
|-----------------------------|------------------------------------------------------------------------------------------------------------------------------|---------------|
| GIN_MODE                    | Run mode for Gin framework. For more info visit the Gin repository.                                                          | debug         |
| LOG_LEVEL                   | Logging level.                                                                                                               | DEBUG         |
| LISTEN_ADDRESS              | Services' port.                                                                                                              | 8080          |
| POSTGRES_HOST               |                                                                                                                              | localhost     |
| POSTGRES_PORT               |                                                                                                                              | 5432          |
| POSTGRES_DB                 | Postgres database. Should be created in advance. After the service started, migrations will be applied (two tables created). | store         |
| POSTGRES_USERNAME           |                                                                                                                              | postgres      |
| POSTGRES_PASSWORD           |                                                                                                                              | postgres      |
| POSTGRES_SSL_MODE           |                                                                                                                              | disable       |
| POSTGRES_CONNECTION_TIMEOUT |                                                                                                                              | 10            |
| TOKEN_TTL                   | Security access token is valid for that period of time (value - time.Duration)                                               | 30m           |
| ACCESS_SECRET               | Private key for access token encryption                                                                                      | -             |
| REFRESH_SECRET              | Private key for refresh token encryption                                                                                     | -             |
