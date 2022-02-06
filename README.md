# File Store
File Store is a simple HTTP file store server, using embedded Key/Value Store (Badger).
The Store allows a simple operation like GET, ADD, UPDATE, REMOVE, LIST files from the server through a client.

----

## Build

Execute the following command to build the binaries
```
go mod download
make build
```

## Run

### Server
```
./bin/server/store
```
```
Usage:
  store [options]
  
Options:
  --config string   Configuration File (default "config.yaml")
  --debug           Debug Mode
  --host string     Server Hostname
  --log string      Logger Mode - (debug, info, warn, error, fatal, panic) (default "info")
  --port int        Server Port (default 8080)
  --tls             Enable TLS
```

#### Configuration File (optional)
Server can be configured using a config.yaml file
```
store --config <config file>
```
##### config.yaml
```
server:
  host: localhost
  port: 8080
  debug: true
  tls: false
  log: debug

cors:
  allowOrigins:
    - https://localhost:8080
  allowMethods:
    - DELETE
    - GET
    - POST
    - PUT
    - PATCH
  maxAge: 86400

database:
  diskless: false
  encryption: false
  cacheSize: 0
  path: temp/database
```

### Client
```
./bin/client/store
```
#### Client Config
The client requires the server endpoint configuration. This only required once.
```
./bin/store config
Store URL: http://<store-host>:<port>
```
```
Usage:
  store [command]

Available Commands:
  add         Add files
  completion  Generate the autocompletion script for the specified shell
  config      Configure store
  count       Word count
  frequency   Word frequency
  get         Get File
  help        Help about any command
  list        List files
  remove      Remove files
  update      Update/Create files
  version     Store version

Flags:
  -h, --help      help for store
  -v, --version   version for store

Use "store [command] --help" for more information about a command.
```

## Docker
The server and the embedded database can be run in a container.
```
make build
make build-docker
```
```
docker run --rm --name store -p 8080:8080 -v $(pwd)/database:/database store:latest
```


