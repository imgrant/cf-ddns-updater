# cf-ddns-updater

Yet another Cloudflare dynamic DNS record update tool

![Go](https://github.com/imgrant/cf-ddns-updater/workflows/Go/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/imgrant/cf-ddns-updater)](https://goreportcard.com/report/github.com/imgrant/cf-ddns-updater)

This is a simple Go utility for updating a DNS record via the [Cloudflare v4 API](https://api.cloudflare.com/). It’s designed for dynamic DNS use cases, and so resolves the current external IP address using the [ipify API](https://www.ipify.org/) for determining public IP addresses, and then points a DNS A record at it. You may or may not find it useful.

## :zap: Quick start

1. Download the appropriate binary from [Releases](https://github.com/imgrant/cf-ddns-updater/releases)
2. Create a [configuration file](#configuration)
3. Run `cf-ddns-updater`

## :toolbox: Requirements

- :house: A valid domain, hosted at Cloudflare
- :ticket: A Cloudflare [API token](https://dash.cloudflare.com/profile/api-tokens), with `Zone:Read` and `DNS:Edit` permissions for the domain you wish to update
- :name_badge: The DNS record that you wish to update must already exist
- :signal_strength: Internet connectivity

## :rocket: Usage

### Command line options

```
  -c    Config file (default: "/home/user/.config/cf-ddns-updater/config.json")
  -n    Check mode (dry run, default: off)
```

The default configuration filename is `config.json` (see below), the exact location is dependent on the operating system; the standard user configuration directories are used (e.g., following the XDG Base Directory Specification for Linux/macOS, which is usually `$HOME/.config`, and `%APPDATA%` for Windows).

The output from `cf-ddns-updater -h` will show the default location in effect. If your configuration file is named  otherwise or located elsewhere, use the `-c` command line parameter to specify the path/filename to use.

To perform a dry run, where the public IP address is determined but the DNS record is not actually updated, use the `-n` flag to enable check mode.


### Configuration

Create a [JSON configuration file](config.example.json)  with the hostname to update, and a valid Cloudflare API token:

```json
{
  "FQDN": "myip.example.com",
  "APIToken": "<Your Cloudflare API token>"
}
```

:bulb: N.b. See [requirements](#toolbox-requirements) above: the DNS record (type `A`) for the FQDN must already exist. The API token must also have sufficient permissions to read the corresponding domain zone and edit the DNS record.

### Example output

#### Updating the DNS record to match the current public IP address

```shell
$ ./cf-ddns-updater
2022/04/21 21:47:38 Current IP address is: 203.0.113.2
2022/04/21 21:47:40 myip.example.com points to 203.0.113.1, the record will be updated.
2022/04/21 21:47:40 Setting myip.example.com => 203.0.113.2 ...
2022/04/21 21:47:41 Success!
```

#### When the DNS record already points to the current public IP address

```shell
$ ./cf-ddns-updater
2022/04/21 21:52:26 Current IP address is: 203.0.113.2
2022/04/21 21:52:28 myip.example.com points to current IP address, no change is needed.
```

#### When check mode is enabled

```shell
$ ./cf-ddns-updater -n
2022/04/21 21:55:08 Current IP address is: 203.0.113.3
2022/04/21 21:55:11 myip.example.com points to 203.0.113.2, the record will be updated.
2022/04/21 21:55:11 Check mode is active, no changes will be made.
```

## :construction: Building

To build the binary yourself, you need [Go](https://go.dev/) (version 1.17 or later).

First, clone the repository:

```shell
git clone git@github.com:imgrant/cf-ddns-updater.git
```

Then, change into the cloned directory and use `go` to compile the source code:

```shell
cd cf-ddns-updater; go build
```
