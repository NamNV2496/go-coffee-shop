# 1. How to run

        1. `cd deployments/` open cmd: `docker-compose up`
        2. Access `https://localhost:9001` and create `coffee` bucket (like picture in below)
        3. `cd ../`
        4. run counter-service `make run-counter`
        5. run kitchen-service `make run-kitchen`
        6. run product-service `make run-product`
        7. run authorization-service `make run-authorization`
        8. run batch-service `make run-batch`

![alt text](docs/minio.png)

# 2. Technologies

        0. Docker
        1. routes - "github.com/gin-gonic/gin"
        2. gRPC - "google.golang.org/grpc"
        3. redis - "github.com/go-redis/redis/v8"
        4. cron job - "github.com/go-co-op/gocron/v2"
        5. kafka - "github.com/IBM/sarama"
        6. minio - "github.com/minio/minio-go"
        7. gorp_migration - "github.com/rubenv/sql-migrate"
        8. goqu - "github.com/doug-martin/goqu/v9"
        9. google wire - "github.com/google/wire"
        10. ratelimit - "golang.org/x/time/rate"
        11. bycrypt - "golang.org/x/crypto/bcrypt"
        12. gopdf - "github.com/jung-kurt/gofpdf"


# 3. Flow

Project structure

        1. authorization-service: for login/register authen, author
        2. counter-service: for receptionist to make order, find items and payment
        3. kitchen-service: kitchen process (in here I made it easily to view order and update finished order)
        4. product-service: manage all product of coffee store
        5. batch-service: to run end of day to make report (you can trigger it by manual)

![alt text](docs/flow.png)

### Sequence diagram

![alt text](docs/sequence_diagram.png)

# 4. How's it work

## Please import `coffee-project.postman_collection.json` to your postman and run it

# 5. Migration sql (auto run when start)

1. Create sql file
2. Run UP to call up migration

```go
migrator := NewMigrator(db)
err = migrator.Up(context.Background())
fmt.Println("Migration up")
if err != nil {
    fmt.Println("failed to execute database up migration")
}
```
3. Result

gorp_migrations to save sql file. Check changelog and build if the file is changed (like liquibase in Java spring boot)

![gorp_migrations](docs/migrateDB.png)

# 6. Authorization

When you call login the API will response 2 token:
- token (main token)
- refreshToken (token with longer expired time. It can help renew token when the main token is expired)

I you meet httpCode = 407 please call API `/api/v1/renewToken`

```json
{
    "code": 200,
    "message": {
        "refreshToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTczMDMzMDksInJvbGVzIjpbIm1lbWJlciJdLCJ1c2VySWQiOiJuYW1udiJ9.PX7DvULwCDwYsy_0UdUx21AXLBS4WV-x--E1wFv-LtI",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTcyOTYzMzksImlhdCI6MTcxNzI5NjMwOSwicm9sZXMiOlsibWVtYmVyIl0sInVzZXJJZCI6Im5hbW52In0.TEH-1oX5PrXpCNBWf4EoPVpCGVAkb1DVezzMtYEC8fU"
    },
    "status": "OK"
}
```

![authorizations](docs/authorizations.png)

<!-- # 9. Dockerfile

locate in `root` dictory

```bash
docker build -t counter-service -f cmd/counter/Dockerfile .

docker build -t kitchen-service -f cmd/kitchen/Dockerfile .

docker build -t product-service -f cmd/product/Dockerfile .
``` -->
