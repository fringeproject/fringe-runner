# Fringe Runner

`fringe-runner` is a tool to fetch new assets over the wire and add it the
[FringeProject](https://fringeproject.com).


## Installation

To easy-install the binary run the following command:

```bash
go get github.com/fringeproject/fringe-runner
```

If you want to install manually or for an other environment, please read the
documentation [here](https://docs.fringeproject.com/runner/)

### Docker
A Docker image of the latest build is available on [DockerHub](https://hub.docker.com/r/fringeproject/fringe-runner):

```bash
docker run -it --rm fringeproject/fringe-runner:latest <cmd>
```

You can also build the image yourself:

```bash
docker build .
---> <image_id>
docker run -it <image_id> <cmd>
```


## How to use

### Configuration

The configuration of `fringe-runner` is store in a YAML file `config.yml`.
Please check the [`config.yml`](./config.yml) provided in the repository and the
[documentation](https://docs.fringeproject.com/runner/#configuration) for more
information.


### Module

This command interacts with the fringe modules.

You can list the modules with `-l/--list`:

```bash
fringe-runner module --list | jq
```

You can execute a module manually `-e/--exec`:

```bash
fringe-runner module --exec <module_slug> <asset_value>
```

### Example

You want to use the `wappalyzer` module on a URL to identify technologies on a
website:

```bash
fringe-runner module -e wappalyzer https://fringeproject.com | jq .[].value
```

This will returns a list of technologies and their versions.


## Available modules

Here is a list of module's slugs available to query public resources:

- [`alienvault`](https://alienvault.com)
- [`bufferover`](https://dns.bufferover.run/)
- [`censys`](https://censys.io/api): need API key (`censys_api_id`) and secret (`censys_api_secret`)
- [`certspotter`](https://sslmate.com/certspotter/): API key is not supported yet
- [`crtsh`](crt.sh/)
- [`dnsdumpster`](https://dnsdumpster.com/): API key is not supported yet
- [`hackertarget`](https://hackertarget.com/)
- [`securitytrails`](https://securitytrails.com/): need an API key sets as `securitytrails_api_key`
- [`shodan`](https://www.shodan.io/): may use an API key sets as `shodan_api_key`
- [`sublist3r`](https://github.com/aboul3la/Sublist3r)
- [`threatcrowd`](https://www.threatcrowd.org/)
- [`threatminer`](https://www.threatminer.org/)
- [`urlscan`](https://urlscan.io/)
- [`virustotal`](https://www.virustotal.com/): use the unofficial API (ui)
- [`yougetsignal`](https://www.yougetsignal.com/)
