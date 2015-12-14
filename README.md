# Tool to get domains of expiring certificates

The tool helps to renew letsencrypt certiicates from cron

## Usage

List domains with expired certificates
```
$ ./letsencrypt-expiring-certs -expire="`date`"
```

List domains of certificates which will expire in 2 weeks
```
$ ./letsencrypt-expiring-certs -expire="`date --date='2 weeks'`"
```

List domains of certificates which will expire at 03/05/2016 11:28
```sh
$ ./letsencrypt-expiring-certs -expire="`date --date '03/05/2016 11:28'`"
```

Use the `-certs-path` option to change the default path of certificates (/etc/letsencrypt/live)

Use the `-pem-name` if the nem of your pem files is different then the default (fullchain.pem)


## How to use in shell script

The following shellscript will renew the certificates which will expire in 1 week. It assumes that the letsencrypt-expiring-certs binary can be found in the PATH.

```sh
#!/bin/sh

DOMAINS=""

EXPIRE_DATE=`date --date='1 week'`
LETSENCRYPT="/path/to/letsencrypt/letsencrypt-auto"

for i in `letsencrypt-expiring-certs -expire "$EXPIRE_DATE"`;do
	DOMAINS="$DOMAINS -d $i"
done


echo $DOMAINS
DIR=/tmp/letsencrypt-auto
mkdir -p $DIR && $LETSENCRYPT --renew certonly --server https://acme-v01.api.letsencrypt.org/directory -a webroot --webroot-path=$DIR --agree-tos $DOMAINS
service nginx reload
```

## Download

You can download prebuild binary from the following link: http://bit.ly/1Z8pIQx
sha1: http://bit.ly/1mkl9V9

Download the file, check if the file not damaged:

```sh
$ wget http://bit.ly/1Z8pIQx
$ wget http://bit.ly/1mkl9V9
$ sha1sum -c sha1.sum
```

Unzip and make the file executable:

```sh
$ gunzip letsencrypt-expiring-certs.1.0.gz
$ chmod +x letsencrypt-expiring-certs.1.0
$ mv letsencrypt-expiring-certs.1.0 letsencrypt-expiring-certs
```
