## Groxy
Little proxy manager builded in go

### Features
- Simple goroutine http server 
- Config file based on .yml
- Proxy pass headers and real ips
- Timeouts
- Rate Limit
- Max Content-Lenght
- Healthchecks


### Setup
Install go packages
```
go mod download
```

Setup your `config.yaml` file based on `config.example.yaml`

Run go server
```
go run ./main.go 
```

For  testing purposes:
```
go test ./... -v -cover
```

If you wanna test with a internal server, config the yml with localhost:3000
```
cd ./example-server
bun run index.ts
```

### Author
- [Alexis Garcia](https://github.com/alexisg24)

