Сборка
```
docker compose build --no-cache
docker compose up -d
```

Поднять, опустить миграции
```
docker compose exec app sh -c 'migrate -path /app/migrations -database "$POSTGRES_URL" up'
docker compose exec app sh -c 'migrate -path /app/migrations -database "$POSTGRES_URL" down'
```


```
protoc --go_out=./proto --go_opt=paths=source_relative \
       --go-grpc_out=./proto --go-grpc_opt=paths=source_relative \
       proto/file_service.proto
```

