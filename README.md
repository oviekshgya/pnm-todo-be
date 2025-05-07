## SHAGYA-TECH SERVICE

-----------------------------
## Tech Stack

**Language:** Go

**Framework:** Fiber

**Database:** MongoDb



## Running Tests

To run tests, run the following command with air toml runner

```bash$
 $ go install github.com/air-verse/air@latest
```

```bash$
 $ air init
```

```bash$
 $ air
```

Too Unit Testing

```bash$
$ go test .\test\user_controller_test.go
```
Too run Local

```bash$ cd build/build-base
$ docker run -d --name rabbitmq \
  -p 5672:5672 -p 15672:15672 \
  rabbitmq:3.12-management
$ docker build -t build-base -f Dockerfile .
$ cd ...
$ docker-compose up --build --remove-orphans -d

```
## Deployment

To deploy this project

```bash$ cd build/build-base
Prose Development Aplikasi Dengan CICD Push main
```

``$ docker build -t pnm-todo-be .
``

``$ docker run -p 3000:3000 pnm-todo-be
``
