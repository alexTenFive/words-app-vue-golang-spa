# app

## Get project
```
go get -u "github.com/alexTenFive/words-app-vue-golang-spa"
```
### Compiles and minifies for production front-end
```
./pkg/http/web/app/yarn run build
```

### Run back-end
```
go build && ./words-app-vue-golang-spa
```
Backend serving on port :9090

### Backend routes
```
POST /api/send
GET /api/results
```

### Address
Front will be on addresses:
Input text - 
```
http://localhost:8080/
```
Get results - 
```
http://localhost:8080/results
```
