# app

## Get project
```
go get -u "github.com/alexTenFive/words-app-vue-golang-spa"
```
### Compiles and minifies for production front-end
if pkg/http/web/app/dist - doesn't exists
```
./pkg/http/web/app/yarn run build
```

### Run back-end
```
go build && ./words-app-vue-golang-spa
```
Site serving on port :9090

### Backend routes
```
POST /api/send
GET /api/results
```

### Address
Front will be on addresses:     
Input text - 
```
http://localhost:9090/
```
Get results - 
```
http://localhost:9090/results
```
