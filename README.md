# Install

To use the provider, you will need to setup your go environment first. Once the
environment is set, get the terraform-provider-abiquo code:

```
go get github.com/abiquo/terraform-provider-abiquo
```

Currently, the provider has been tested with terraform v0.11.3, so, you need to
change the tag the terraform repository is pointing to in your go environment.

# Configuring the provider

The provider currently support two configurations, Basic Auth and Oauth.
Also, you will need to configure your OpenSSL setup to validate your deployment
certificate. Refer to the OpenSSL documentation for the platform where you will
be running the provider.

## Basic Auth

```
provider "abiquo" {
  endpoint       = "https://fqdn:443/api"
  username       = "admin"
  password       = "xabiquo"
}
```

## Oauth

```
provider "abiquo" {
  endpoint       = "https://fqdn:443/api"
  consumerkey    = "5336cd80-d17b-488a-8917-518a12ee366a"
  consumersecret = "nuDmkp1t4qmcyxGVfVsujmVqJ5VexeLIymvBA5Oy"
  token          = "7ea0959c-82f1-4013-ab2b-6648999f3915"
  tokensecret    = "TgYSC9Y4TX3r+p9q3F8DhcJ3J9FFXOCmPD6pAKw1G31wTUAtlTgZTMJjDT/jS2F4K2DUYX6Py641PLeBkKMntS+GdKkO09ajkil9ZH67Fa0="
}
```

# Examples

Check the examples folder to get an idea of how to use the provider and the
available resources.
