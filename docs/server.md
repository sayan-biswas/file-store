# Store Server

The store server uses the Gin web framework to spin up the server and the router
 
It provides the following endpoints/operations on the store.

## Endpoints

- **/**
  - **GET** - Used for health/liveliness probe, and also to identify store URL from client
- **/store**
  - **GET** - Fetch a file from the server
    - Ex: /store?file=*filename*
  - **POST** - Push a file from the server
    - Ex: /store (*data sent as array of bytes in multipart form*)
  - **PUT** - Update an existing file in server. Creates a new file when the file doesn't exist on the server
    - Ex: /store (*data sent as array of bytes as multipart form*)
  - **DELETE** - Remove a file from server
    - Ex: /store?file=*filename*
- **/store/list**
  - **GET** - Get a list of file and details from the store
- **/store/check/file**
  - **GET** - Check if a file exists on the server
    - Ex: /store/check/file?file=*filename*
- **/store/check/sha**
  - **GET** - Check if same data exist on the server
    - Ex: /store/check/sha?sha=*base16 encoded 256 bit checksum of the data*
- **/store/count**
  - **GET** - Get the total word count from all the files on the server
- **/store/frequency**
  - **GET** - Get frequency of word in ascending/descending order from all the files on the store

## Configuration

The server can be configured in the following methods (order is priority wise)..
- Command line flags
- Environment Variables
- YAML config file
- Default Values

### Command Line

The below code start the server at the giver **PORT** with **DEBUG** mode on

```
store --port 4000 --debug
``` 
Command line flags has the highest priority and hence overrides all other configuration methods 

**Full list of flags**
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

### Environment Variables

Any configurable property can be configured using environment variables. The variable name should be in the follwing syntax **STORE_** *PROPERTY*

The below code will set the loggin mode to show only **info** and the port to **4000**

```
export STORE_LOG=info
export STORE_PORT=4000
```

### Configuration File

The server can be configured using a configuration file where the properties declared in YAML syntax.

Below is the full config file

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

By defalut config file is searched in the below mentioned path with the name **config.yaml**
- root
- /etc/store/config
- $HOME/store/
- ./config

The location of the config can also be provided through the **--config** flag.
```
store --config /usr/store/config.yaml
``` 

