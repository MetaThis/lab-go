# lab-go

Example of a simple microservice in Go, demonstrating some patterns that proven useful.

## Setup

Clone into the directory structure expected by Go.

`$GOPATH/src/github.com/MetaThis/lab-go`

```
cd $GOPATH/src/github.com/MetaThis/lab-go
go get
```
Run tests
```
go test -v -cover
```
Run server locally
```
go build
./lab-go
```
With the server running, you can use Postman or cURL to make requests.
```
POST http://localhost:8080/lab/instrument/1/samples
```
Posting a body such as 
```
[{"id":1},{"id":2},{"id":999}]
``` 

