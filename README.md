# zname

Search Route53 DNS zones.

## Installation

Grab the [latest release](https://gitlab.com/zapier/zname).

Or build locally:
```
CGO_ENABLED=1 go build -o zname cmd/zname/main.go 
```
## Usage

Search for `<WORD>` in Route53 record names or targets:

```
zname <WORD> [-r]

Flags:
    -r  Rebuild the cache file
```
