# cdncheck

This tool is a CLI wrapper for [ProjectDiscovery](https://github.com/projectdiscovery)'s [cdncheck](https://github.com/projectdiscovery/cdncheck) library - "Helper library that checks if a given IP belongs to known CDN ranges (akamai, cloudflare, incapsula and sucuri)".

## Installation

```
GO111MODULE=on go get -u -v github.com/enenumxela/cdncheck/cmd/cdncheck
```

## Usage

```
cdncheck -h
```

```
USAGE:
  cdncheck [OPTIONS]

OPTIONS:
  -iL, --input-list   input IP list (use `iL -` to read from stdin)
   -c, --concurrency  number of concurrent threads (default: 10)
```