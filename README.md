Fringe Runner
=============

The runner fetches new assets over the wire.

## Requirements

Fringe runner is written in Go.


## How to use

### module

This command interacts with the fringe modules.

You can list the modules with the `-l/--list` argument:

```bash
fringe-runner module --list
```

You can execute a module manually with the following command:

```bash
fringe-runner module --exec <module_slug> <asset_value>
```
