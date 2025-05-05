## Generate the protobuf file
```
protoc --go_out=. --go-grpc_out=. proto/migrate.proto
```

## Generate the server.cert and server.key
```
openssl req -x509 -newkey rsa:4096 -nodes -keyout cert/server.key -out cert/server.crt -days 365 -subj "/CN=localhost" -config "nul"
```