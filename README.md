# URL Shortener

This is a URL shortener built using [Go](https://golang.org/). The main objective was to learn how to use the
[Hexagonal Architecture](https://en.wikipedia.org/wiki/Hexagonal_architecture_(software)). The service currently
supports the following databases

- [Redis](https://redis.io)
- [MongoDB](https://www.mongodb.com)

## Run Locally

Clone the project

```bash
  git clone git@github.com:jwambugu/hexagonal-architecture-url-shortner.git
```

Go to the project directory

```bash
  cd hexagonal-architecture-url-shortner
```

Set the environment variables (Redis)

```bash
  export MS_REPOSITORY=redis
  export REDIS_URL=redis-url-goes-here
```

Set the environment variables (MongoDB)

```bash
  export MS_REPOSITORY=mongo
  export MONGO_URL=mongo-url-goes-here
  export MONGO_DB=redirects
```

Start the server

```go
  go run main.go
```

## API Reference

#### Get Shortened URL

```http
  GET /{code}
```

| Parameter | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `code` | `string` | **Required**. |

#### Shorten URL

```http
  POST /
```

| Parameter | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `url`      | `string` | **Required**. URL to shorten |

## References

- [Tensor Programming](https://youtu.be/rQnTtQZGpg8?list=PLJbE2Yu2zumAixEws7gtptADSLmZ_pscP)

  