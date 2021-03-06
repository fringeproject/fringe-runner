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

You can list the modules with `-L/--list-modules`:

```bash
fringe-runner module --list-modules | jq
```

You can execute a module manually:

```bash
fringe-runner module -m module_slug -a asset_value
```

The `asset` argument can also be a file containing assets to execute the same
module on all the assets from the file.

```bash
fringe-runner module -m crtsh -a assets.txt
```

### Parse a local file

This command parse a file separated by new lines:

```bash
fringe-runner parse -p <file_path> | jq ".[].value"
```


## Use cases

### Get website technologies with Wappalyzer

You want to use the `wappalyzer` module on a URL to identify technologies on a
website:

```bash
fringe-runner module -m wappalyzer https://fringeproject.com | jq .[].value
```

This will returns a list of technologies and their versions.


### Take screenshots a list of hostnames

The file `hostnames.txt` contains a list of hostname (1 by line) and you want
to take a screenshot of the webservers runnings on those:

First, configure the screenshot renderer (see [`config.yml`](./config.yml)) then
type the following command:

```bash
fringe-runner module -m http-probe -a hostnames.txt -w workflows/screenshot.yml
```

The runner execute the module `http-probe` on each line of the `hostnames.txt`
file. This module checks for web-servers on HTTP (80) and HTTPS (443). Then the
workflow `screenshot.yml` take a live screenshot for every listening web-server.


## Available modules

Here is a list of module's slugs available to query public resources:

- [`alienvault`](https://alienvault.com)
- [`bufferover`](https://dns.bufferover.run/)
- [`censys`](https://censys.io/api): need API key (`censys_api_id`) and secret (`censys_api_secret`)
- [`certspotter`](https://sslmate.com/certspotter/): API key is not supported yet
- [`crtsh`](crt.sh/)
- [`dnsdumpster`](https://dnsdumpster.com/): API key is not supported yet
- [`github-subdomains`](https://developer.github.com/v3/search/#search-code): need API token (`github_api_token`)
- [`hackertarget`](https://hackertarget.com/)
- [`securitytrails`](https://securitytrails.com/): need an API key sets as `securitytrails_api_key`
- [`shodan`](https://www.shodan.io/): may use an API key sets as `shodan_api_key`
- [`sublist3r`](https://github.com/aboul3la/Sublist3r)
- [`threatcrowd`](https://www.threatcrowd.org/)
- [`threatminer`](https://www.threatminer.org/)
- [`urlscan`](https://urlscan.io/)
- [`virustotal`](https://www.virustotal.com/): use the unofficial API (ui)
- [`whoisxmlapi`](https://reverse-ip.whoisxmlapi.com/api/documentation/making-requests): need API key (`whoisxmlapi_key`)
- [`yougetsignal`](https://www.yougetsignal.com/)


The following modules are still in progress:

- `nessus`: Add a new scan to a module instance (`nessus_endpoint`, `nessus_username` and `nessus_password`)
