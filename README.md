# sdx-encrypted-submitter

A command line utility to copy a files contents to a rabbit exchange

Building and Running

Install Go and ensure that your GOPATH env is set ( usually it's ~/go).

```
go get  github.com/ONSdigital/sdx-encrypted-submitter/
cd $GOPATH/src/github.com/ONSdigital/sdx-encrypted-submitter
go get -u github.com/streadway/amqp
go get -u gopkg.in/yaml.v2
go get -u gopkg.in/square/go-jose.v2
go build
go install
```

### Preparing Keys

This requires that new keys are generated and added to sdx-decrypt so that this utility does not 
need knowledge of existing eq and sdx keys. To do this we need to generate 2 pairs of keys, one pair that 
is supplemental to existing eq keys and the other which is supplemental to sdx keys.
To generate the key pairs we need to create keys with the correct name format.
```
sdc-<service>-submission-<key-use>-<key-type>-<version>.pem
```
Where <service> is set to 'submitter' for upstream keys and 'sdx' for sdx keys.
*Note 'submitter' is used in place of 'eq' so that we do not clash with eq keys.*
<key-use> is set to 'signing' for upstream keys and 'encryption' for sdx keys.
<key-type> is set to public or private.
<version> is set to a version number. Note to avoid collisions with sdx keys the sdx version number starts at 1000.

```
openssl genrsa out <private_key_name>.pem 4096 
openssl rsa -pubout -in <private_key_name>.pem -out <public_key_name>.pem
```

So in total we have 4 pem files, 2 public, 2 private. 
Now we have to import the correct ones into the keys.yaml file of sdx-decrypt. 
We want the upstream public key pem file and the sdx private key pem file loaded into the keys.yml file
Create a folder and place these files in there then following the instructions from sdc-cryptography ReadMe
```
generate_keys.py <path to files to add>
```
The keys that where not added are the ones that we reference on the command line with :
```
-s <path to upstream private pem file>
-e <path to sdx public pem file>
```

### Notes: 
* Dep support likely to follow. (https://github.com/golang/dep)
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
-e|If specified the keyfile used to encrypt the source file (up stream private key)|Yes
-s|If specified the keyfile used to sign the source file (sdx public key)|Yes
-h|Get a list of command line arguments. Prevents processing , other parameters ignored
 


### Examples
```
sdx-encrypted-submitter -n RabbitUser  -p RabbitPassword  -f SourceFile -e encryption_key_pem_file -s signing_key_pem_file
```

Reads the file encrypts and signs the contents and writes to the exchange specified in the sdx-encrypted-submitter.yml file, using the routing key specified there.

```
sdx-encrypted-submitter -h  
```
Shows command line help

```
port: 5672
host: "localhost"
exchange: "message"
vhost: "%2f"
routingkey: "survey"

```
Example sdx-encrypted-submitter.yml working locally on developer machine


### License

Copyright (c) 2017 Crown Copyright (Office for National Statistics)

Released under MIT license, see [LICENSE](LICENSE) for details.