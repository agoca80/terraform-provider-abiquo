#!/bin/sh -x

TS=$(date +%s)
cd ../../hashicorp/terraform
make testacc TEST=../../abiquo/terraform-provider-abiquo TESTARGS="$* -coverprofile /tmp/coverage-$TS"
go tool cover -html=/tmp/coverage-$TS

