# Tool to get domains of expiring certificates

The tool helps to renew letsencrypt certiicates from cron

## Install

The tool is written in go, so you will need a go compiler.

It's often shipped with Linux distribution, for example on debian you can install the _golang_ package:

```sh
$ sudo apt-get install golang
```

Then to install the tool run the following command:

```
$ export GOPATH="$HOME/gocode";
$ go get github.com/Ajnasz/letsencrypt-expiring-certs
```

This will download and also compile the code:

In the _$HOME/gocode/bin/_ folder you will find the binary.

In the _gocode/src/github.com/Ajnasz/letsencrypt-expiring-certs/_ folder you will find the source code.

In the folder of the source code you can build tool again by using the `go build` command.

If you use other OS then Linux or other distribution, read the documenatition: https://golang.org/doc/install

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

You can download prebuild binary from the following link: http://bit.ly/1m4Xus1
sha1: http://bit.ly/1mkl9V9

Download the file, check if the file not damaged:

```sh
$ wget http://bit.ly/1m4Xus1
$ wget http://bit.ly/1mkl9V9
$ sha1sum -c sha1.sum
```

Unzip and make the file executable:

```sh
$ gunzip letsencrypt-expiring-certs.2.0.1.gz
$ chmod +x letsencrypt-expiring-certs.2.0.1
$ mv letsencrypt-expiring-certs.2.0.1 letsencrypt-expiring-certs
```
