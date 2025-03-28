# Tool to get domains of expiring certificates

The tool helps to check if the certificates are expired or will expire soon.

## Usage

### Check remote domains

Use the `-domains` flag to check remote domains
In that case the tool won't check the filesystem, but the domains you provide.

An example to check the expiration of the example.org and wikipedia.org domains:

```sh
$ ./letsencrypt-expiring-certs -domains example.org,wikipedia.org -expire "$(date -Is --date 2 year)" -print-date -quit
*.example.org   2026-01-15T23:59:59Z    x509: certificate has expired or is not yet valid: current time 2027-03-28T09:42:53+02:00 is after 2026-01-15T23:59:59Z
example.org     2026-01-15T23:59:59Z    x509: certificate has expired or is not yet valid: current time 2027-03-28T09:42:53+02:00 is after 2026-01-15T23:59:59Z
*.wikipedia.org 2025-10-17T23:59:59Z    x509: certificate has expired or is not yet valid: current time 2027-03-28T09:42:53+02:00 is after 2025-10-17T23:59:59Z
wikimedia.org   2025-10-17T23:59:59Z    x509: certificate has expired or is not yet valid: current time 2027-03-28T09:42:53+02:00 is after 2025-10-17T23:59:59Z
mediawiki.org   2025-10-17T23:59:59Z    x509: certificate has expired or is not yet valid: current time 2027-03-28T09:42:53+02:00 is after 2025-10-17T23:59:59Z
# ...
```

### Check local letsencrypt certificates

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

## Build

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

## Renew certificates

You can see an example cron script in the scripts folder, see the scritps/cron.sh
Make sure every path is correct in the script.
