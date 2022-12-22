# ip2location redirect

Plugins for getting the country code of the client ip from ip2location database and redirect to specific url when (no)match the target countries


Configuration
---
To configure this plugin you should add its configuration to the Traefik dynamic configuration as explained [here](https://doc.traefik.io/traefik/getting-started/configuration-overview/#the-dynamic-configuration). The following snippet shows how to configure this plugin with the File provider in TOML and YAML:

Static:
```yaml
experimental:
  pilot:
    token: xxx

  plugins:
    ip2location_redirect:
      modulename: github.com/jiangwennn/traefik-plugin-ip2location-redirect
      version: v0.1.0
```

Dynamic:
```yaml
http:
  middlewares:
    my-plugin:
      plugin:
        ip2location_redirect:
          regions: ["CN"],
          redirectUrl: "https://github.com"
          noMatch: false # optional
          permanent: false # optional
          fromHeader: X-User-IP # optional
          disableErrorHeader: false # optional
```


Options
---

### Regions 
`regions` | `[]string` | *Required* 

Array of country codes in 2 letters, eg: `CN`,`UK`

### RedirectUrl 

`redirectUrl` | `string` | *Required* 

The url address redirect to

### NoMatch

`noMatch` | `bool` | *Default*: `false`

If `true`, redirect action performs when the ip doesn't belong to any of the regions configured

### Permanent

`permanent` | `bool` | *Default*: `false`

If `true`, redirect status code will be `301 Moved Permanently`, otherwise `302 Found`

### FromHeader

`fromHeader` | `string` | *Default*: `empty`

If defined, IP address will be obtained from this HTTP header. `Remote-Addr` By default.

### DisableErrorHeader

`disableErrorHeader` | `bool` | *Default*: `false`


Errors
---

If any error occurred, this error will be placed to X-IP2LOCATION-REDIRECT-ERROR header


Reference
---

- traefik-plugin-ip2location

    https://github.com/negasus/traefik-plugin-ip2location

- Country Codes

    https://www.countrycode.org/