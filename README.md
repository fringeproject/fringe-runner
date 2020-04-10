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
