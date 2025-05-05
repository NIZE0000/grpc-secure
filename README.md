# POC GRPC FOR MIGRATION

## Generate the protobuf file
```
protoc --go_out=. --go-grpc_out=. proto/migrate.proto
```

## Generate the server.cert and server.key
<!-- ```
openssl req -x509 -newkey rsa:4096 -nodes -keyout cert/server.key -out cert/server.crt -days 365 -subj "/CN=localhost" -config "nul"
``` -->

```
openssl req -x509 -nodes -newkey rsa:4096 -keyout cert/server.key -out cert/server.crt -days 365 -config cert/san.cnf -extensions req_ext
```

## Architecture 

```css
[ Source System ]
     │
     ▼
[ Kafka Topic ]
     │
     ▼
[ Go Worker Service ]  →  [ gRPC Destination Service ] → [ Database ]

```