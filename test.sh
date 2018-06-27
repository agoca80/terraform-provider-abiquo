#!/bin/sh -x

cd ../../hashicorp/terraform
make testacc TEST=../../abiquo/terraform-provider-abiquo TESTARGS="$* -coverprofile $HOME/coverage"
echo go tool cover -html=$HOME/coverage

