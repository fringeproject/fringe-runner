Fringe Runner
=============

The runner fetches new assets over the wire.

## Requirements

Fringe runner is written in Go.


## How to use

### Configuration

Fringe-Runner uses environment variables as configuration. It's possible to
create a `.env` in the current directory. This file will be loaded at the startup.

```
LOG_LEVEL=<debug/info/...>
HTTP_PROXY=<proxy URL>
VERIFY_CERT=<true/false>
```


### Module

This command interacts with the fringe modules.

You can list the modules with the `-l/--list` argument:

```bash
fringe-runner module --list | jq
```

You can execute a module manually with the following command:

```bash
fringe-runner module --exec <module_slug> <asset_value>
```


## Modules to fetch assets from publics API

Here is a list of module's slugs available to query public resources:

- [`bufferover`](https://dns.bufferover.run/)
- [`certspotter`](https://sslmate.com/certspotter/): API key is not supported yet
- [`crtsh`](crt.sh/)
- [`dnsdumpster`](https://dnsdumpster.com/): API key is not supported yet
- [`hackertarget`](https://hackertarget.com/)
- [`securitytrails`](https://securitytrails.com/): need an API key sets as `SECURITYTRAILS_API_KEY`
- [`shodan`](https://www.shodan.io/): may use an API key sets as `SHODAN_API_KEY`
- [`sublist3r`](https://github.com/aboul3la/Sublist3r)
- [`threatcrowd`](https://www.threatcrowd.org/)
- [`urlscan`](https://urlscan.io/)
- [`virustotal`](https://www.virustotal.com/): use the unofficial API (ui)
