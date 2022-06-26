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

### Create TLS certificate
Certificates will be in home user directory in subdirectory '**cert**'.  
Inside this directory we will create two files:
- **localhost.crt** - self-signed-certificate
- **localhost.key** - private key  

Creation command:  
`openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout cert/localhost.key -out cert/localhost.crt`  

# Docker - how to use mend as docker container
Firstly you must install docker on your computer.  

## Create container  
Go to root diectory of the project and:  
`docker build --tag mend .`

## Start mend as container
When all was OK:  
`docker run -p 8010:8010 mend`

## Check   
Run 'Postman', select GET command with address 'https://localhost:8010/user/count'.
As result you should see:  
`{
    "count": 0
}`
