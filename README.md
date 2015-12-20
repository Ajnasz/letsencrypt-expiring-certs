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

You can see an example cron script in the scripts folder, see the scritps/cron.sh
Make sure every path is correct in the script.

## Download

You can download prebuild binary from the following link: http://bit.ly/1TaQKD6
sha1: http://bit.ly/1mkl9V9

Download the file, check if the file not damaged:

```sh
$ wget hhttp://bit.ly/1TaQKD6
$ wget http://bit.ly/1mkl9V9 # sha1 sum
$ sha1sum -c sha1.sum
```

Unzip and make the file executable:

```sh
$ gunzip letsencrypt-expiring-certs.2.0.0.gz
$ chmod +x letsencrypt-expiring-certs.2.0
$ mv letsencrypt-expiring-certs.2.0.0 letsencrypt-expiring-certs
```
