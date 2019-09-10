# Bemobi

## Client:

- index.html (Vuejs, Axios and Bootstrap 4)

## Prerequisites:
- GNU Make
- Docker 19.03++
- docker-compose 1.23++

---
## Enviroment:
```
APP_HOST=localhost
APP_PORT=8080
    
MONGODB_DATABASE=bemobi
MONGODB_HOST=localhost
MONGODB_PORT=27017
MONGODB_USER=
MONGODB_PASS=
  ```
  
## Build:
```
$ make build
```

## Run:
```
$ make run
```

## Test:
```
$ make test
```

## Documentation:

### Shortener URL:
```
POST http://<APP_HOST>:<APP_PORT>/create?url=<string>&[CUSTOM_ALIAS}
```

---
### Retrieve URL
```
GET http://<APP_HOST>:<APP_PORT>/{alias}
```

---
### MoreVisited
```
GET http://<APP_HOST>:<APP_PORT>/moreVisited
```
