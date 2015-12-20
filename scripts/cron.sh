#!/bin/sh

EXPIRE_DATE=`date --date='1 week'`
LETSENCRYPT="/path/to/letsencrypt/letsencrypt-auto" # TODO change to the path of letsencrypt-auto
LETSENCRYPT_EXPIRING_CERTS="/usr/local/sbin/letsencrypt-expiring-certs" # TODO change to the path of letsencrypt-expiring-certs
DRY=1 # TODO Change to 0 in production

DIR=/tmp/letsencrypt-auto
crete_temp_folder() {
	echo "Create temporary folder"
	if dry_mode;then
		return
	fi

	mkdir -p $DIR

	MKDIRSTATUS=$?

	if [ "$MKDIRSTATUS" -ne "0" ];then
		exit $MKDIRSTATUS
	fi
}

reload_services() {
	echo "Reload nginx"
	if ! is_dry;then
		service nginx reload
	fi
}

line_to_domains() {
	LINE=$@

	DOMAINS=""
	for DOMAIN in $LINE; do
		DOMAINS="$DOMAINS -d $DOMAIN"
	done

	echo $DOMAINS
}

is_dry() {
	[ "$DRY" -ne "0" ]
}

renew_domains() {
	DOMAINS=$@

	if is_dry;then
		echo $LETSENCRYPT --renew certonly --server https://acme-v01.api.letsencrypt.org/directory -a webroot --webroot-path=$DIR --agree-tos $DOMAINS
	else
		$LETSENCRYPT --renew certonly --server https://acme-v01.api.letsencrypt.org/directory -a webroot --webroot-path=$DIR --agree-tos $DOMAINS
	fi

	RENEW_STATUS=$?

	if [ "$RENEW_STATUS" -ne "0" ];then
		>&2 echo "Renew failed"
		exit $RENEW_STATUS
	fi
}

check_letsencrypt() {
	if [ ! -x "$LETSENCRYPT_EXPIRING_CERTS" ];then
		>&2 echo "letsencrypt-expiring-certs not executable, set LETSENCRYPT_EXPIRING_CERTS variable to the correct path"
		exit 2;
	fi
}

get_lines() {
	LINES=`$LETSENCRYPT_EXPIRING_CERTS -expire "$EXPIRE_DATE"`

	GET_LINES_STATUS=$?

	if [ "$GET_LINES_STATUS" -ne "0" ];then
		exit $GET_LINES_STATUS
	fi

	echo $LINES
}

main() {
	if is_dry;then
		echo "Running in dry mode, change DRY to 0 in the script on live environment"
	fi

	check_letsencrypt
	LINES=`get_lines` || exit $?

	echo $LINES | while read LINE; do
		DOMAINS=`line_to_domains $LINE`
		if [ ! -z "$DOMAINS" ];then
			renew_domains $DOMAINS
		fi
	done

	reload_services
}

main
