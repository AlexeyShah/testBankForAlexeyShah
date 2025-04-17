## Обновляем сваггер доку
```
cd src/bank-service
go get -u github.com/swaggo/swag/cmd/swag
go get -u github.com/swaggo/swag
go get -u github.com/swaggo/gin-swagger
go get -u github.com/swaggo/files
go get -u github.com/swaggo/http-swagger

swag init -g ./cmd/server/main.go -o docs --parseDependency
```
http://localhost:8080/swagger/index.html


## Собираем докер
```
cd src/bank-service
docker build -t cr.yandex/crpr6h01ptp35hb82u6b/bank-service:dev .
```

## Локальная отладка
```
Поддерживается Windows VisualCode, конфиг - launch bank-service
.vscode/launch.json
```