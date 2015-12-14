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

```
#!/bin/sh

DOMAINS=""

EXPIRE_DATE=`date --date='1 week'`

for i in `letsencrypt-expiring-certs -expire "$EXPIRE_DATE"`;do
	DOMAINS="$DOMAINS -d $i"
done


echo $DOMAINS
DIR=/tmp/letsencrypt-auto
mkdir -p $DIR && /home/ajnasz/src/letsencrypt/letsencrypt-auto --renew certonly --server https://acme-v01.api.letsencrypt.org/directory -a webroot --webroot-path=$DIR --agree-tos $DOMAINS
service nginx reload
```
