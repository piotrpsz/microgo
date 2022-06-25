# mend

The program is example of microservice in Go with using Gin framework.  
The server accept only **TLS** connection (self-signed certificate).  
You can use **REST**-api to check how works all **CRUD** elements (Create-Read-Update-Delete).

# External dependencies
- [gin] go get -u github.com/gin-gonic/gin
- [logrus] go get -u github.com/sirupsen/logrus

# Prepare enviromment
Here are few information how to prepare local emvironment befor start microservice.

## Working directory of the microservice 
All data of the microservice are located in user home directory in subdirectory '**.mend**'.

`cd`  
`mkdir .mend`    
`cd .mend`    

In this directory will be located log file of the microservice.  
And certificates used on TLS.

### Create TLS certificate
Certificates will be in subdirectory '**cert**'.  
Inside this directory we will create two files:
- **localhost.crt** - self-signed-certificate
- **localhost.key** - private key  

Creation command:  
`openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout cert/localhost.key -out cert/localhost.crt`  


# Swagger - How to check available API
You can see available API using swagger. The app is prepared for this salution.  
Firstly you must install swagger. I explain how do it, only for persons which have installed Go.

## Install current version from source
We downloads the original code:  
`git clone https://github.com/go-swagger/go-swagger`  
`cd go-swagger/cmd/swagger`  
`go build swagger.go`  
`cp swagger $GOPATH/bin`  

When finished go to the root directory of the project and:  
`Projects/Go/mend % swagger generate spec -o ./swagger.json`  
`swagger serve ./swagger.json`  


