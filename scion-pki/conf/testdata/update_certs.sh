#! /bin/bash

set -ueo pipefail

PLAYGROUND=$(pwd)/../../../tools/cryptoplayground
PUBDIR="$SAFEDIR"
KEYDIR="$SAFEDIR"

loc="Zürich"
IA="1-ff00:0:110"

. "$PLAYGROUND/crypto_lib.sh"

navigate_pubdir
basic_conf && sensitive_conf && regular_conf
prepare_ca
sed -i \
    -e 's/{{.Country}}/CH/g' \
    -e "s/{{.State}}/$loc/g" \
    -e "s/{{.Location}}/$loc/g" \
    -e "s/{{.Organization}}/$loc/g" \
    -e "s/{{.OrganizationalUnit}}/$loc InfoSec Test Squad/g" \
    -e "s/{{.ISDAS}}/$IA/g" \
    basic.cnf
for cnf in *.cnf
do
    sed -i \
    -e "s/{{.ShortOrg}}/$loc/g" \
    $cnf
done

# Generate certificates
in_docker 'navigate_pubdir && gen_sensitive && gen_regular'
