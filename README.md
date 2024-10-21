# Task Manager

## Sample Request and Response
### List Tasks
* Request
```
GET  http://localhost:8080/v1/tasks
Content-Type: application/json
```

* Response
```
HTTP/1.1 200 OK
Content-Type: application/json

[
  {
    "id": "7494d1aa-21f1-4003-8504-70602e167839",
    "name": "task-name",
    "status": 0
  }
]
```

### Create a Task
* Request
```
POST http://localhost:8080/v1/tasks
Content-Type: application/json

{
    "name": "task-name",
    "status": 0
}
```

* Response
```
HTTP/1.1 201 Created
Content-Type: application/json

{
    "id": "7494d1aa-21f1-4003-8504-70602e167839",
    "name": "task-name",
    "status": 0
}
```

### Update a Task
* Request
```
PUT  http://localhost:8080/v1/tasks/7494d1aa-21f1-4003-8504-70602e167839
Content-Type: application/json

{
    "name": "task-name",
    "status": 1
}
```

* Response
```
HTTP/1.1 200 OK
Content-Type: application/json

{
    "id": "7494d1aa-21f1-4003-8504-70602e167839",
    "name": "task-name",
    "status": 1
}
```

### Delete a Task
* Request
```
DELETE  http://localhost:8080/v1/tasks/7494d1aa-21f1-4003-8504-70602e167839
Content-Type: application/json
```

* Response
```
HTTP/1.1 204 No Content
Content-Type: application/json
```

## Run API in Docker
1. Build the Docker image: `docker build -t go-crud-api .`
2. Run the Docker container: `docker run -p 8080:8080 go-crud-api`
3. Verify API is Running: `curl http://localhost:8080/v1/tasks`