# Generate Mocks
```
mockery --name Cacher --dir ./src/seabackend/infrastructure/ --output ./src/seabackend/infrastructure/mocks
```

Generate code tagged with `go:generate` comments:
```
go generate ./...
```

# Swagger
Generate empty spec:
```
swagger init spec --title "gosea" --description "gosea - SEA Go Sample Application" --version 1.0.0 --scheme http ./src/seaswagger/
```

Swagger Online-Editor:
https://studio.apicur.io/