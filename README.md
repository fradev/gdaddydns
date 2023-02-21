# gdaddydns
Small util to manipulate GO Daddy DNS entries via CMD line

It allows to list,add,delete DNS entries using go daddy API.

It needs a simple configuration file (by default ~/.gdaddydns.json) with the following syntax

```
{
  "domains": [
    {"name": "example.com", "api_key": "EXAMPLE_KEY", "api_secret": "EXAMPLE_SECRET"},
    {"name": "me.com", "api_key": "ME_KEY", "api_secret": "ME_SECRET"},
    {"name": "xxxx.net", "api_key": "XXXX_KEY", "api_secret": "XXXX_SECRET"}
   ]
}
```

gdaddydns help should be self explicative

## Examples 

`
gdaddydns add --domain example.com --name localhost --data 127.0.0.1 --type A
`

`
gdaddydns list --domain example.com  --type MX
`

`
gdaddydns del --domain example.com  --name www --type A
`

