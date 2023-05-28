# gin-gonic REST API using MongoDB
Test Gin Gonic REST API using MongoDB
 
## How to start the Gin-gonic application
```shell
go run main.go
docker-compose up
```

## CURL commands to interact with the REST API
### Get all the interviews
```shell
curl -X GET 'localhost:8080/interviews'
```

### Create a new interview
```shell
curl -X POST 'localhost:8080/interviews' \
--data '{"subject": "นัดสัมภา test_backend", "detail": "knowledge: JAVA, HTML, CSS, GO Lang, Javascript", "create_by": "admin"}'
```

### Update a book entry
```shell
curl -X PUT 'localhost:8080/interviews/53fbf4615c3b9f41c381b6a3' \
--data '{"subject": "นัดสัมภา ai", "detail": "knowledge: mssql, mongodb, oracle", "updated_by":"admin", "status":"In Progress" }'
```

### Delete a book entry
```shell
curl -X DELETE 'localhost:8080/interviews/53fbf4615c3b9f41c381b6a3'
```

## Required Docker images

Create a MongoDB container
```shell
docker pull mongo
docker create --name mongodb -it -p 27017:27017 mongo
```
