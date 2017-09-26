# sdx-encrypted-submitter

A command line utility to copy a files contents to a rabbit exchange

Building and Running

Install Go and ensure that your GOPATH env is set ( usually it's ~/go).

```
go get  github.com/ONSdigital/sdx-encrypted_submitter/
cd $GOPATH/src/sdx-encrypted-submitter
go get -u github.com/streadway/amqp
go get -u gopkg.in/yaml.v2
go build
```

### Notes: 
* govender support likely to follow. (https://github.com/kardianos/govendor)
* No unit tests yet
* Assumes exchanges , queues and bindings are in place . Does NOT create them if they are not.


### Usage

Arguments for use are split between sdx-encrypted-submitter.yml and command line arguments

#### Config file settings (sdx-encrypted-submitter.yml) 

Name | Description
-----|--------- 
port|The rabbit port to use 
host|The ip of the rabbit queue host
exchange|The name of the exchange to post to
vhost|The virtual rabbit host to use 
routingkey|The routing key to use when saving the message

#### Command line arguments

Name | Description | Mandatory
-----|---------|----------
-n|Name to use when connecting to rabbit|Yes
-p|Password to use when connecting to rabbit|Yes
-f|The source file holding the message to send to rabbit|Yes
-e|If specified the keyfile used to encrypt the source file (not yet implemented)|No
-s|If specified the keyfile used to sign the source file (not yet implemented)|No
-h|Get a list of command line arguments. Prevents processing , other parameters ignored
 


### Examples
```
sdx-submitter -n RabbitUser  -p RabbitPassword  -f EncryptedAndSignedFile
Reads the file and writes to the exchange specified in the sdx-submitter.yml file, using the routing key specified there.
```
```
sdx-submitter -h  
Shows command line help
```

### Coming soon to a terminal near you ........
```
sdx-submitter -n RabbitUser  -p RabbitPassword  -f EncryptedAndSignedFile -e EncryptFile
Reads the file content and encrypt using the key file content 

```



