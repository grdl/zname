# Zname

[![build](https://github.com/grdl/zname/actions/workflows/build.yml/badge.svg)](https://github.com/grdl/zname/actions/workflows/build.yml)
[![release](https://github.com/grdl/zname/actions/workflows/release.yml/badge.svg)](https://github.com/grdl/zname/actions/workflows/release.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/grdl/zname)](https://goreportcard.com/report/github.com/grdl/zname)

It's called "Zname" because it searches through "CNAMEs" ;)

## Installation

Grab the [latest release](https://github.com/grdl/zname/releases/latest).

Or use Homebrew:
```
brew install grdl/tap/zname
```


## Usage

Search for `<WORD>` in DNS record names or targets. Currently, only AWS Route53 is supported.

```
Zname - search through your cloud DNS records.

Usage:
  zname <WORD> [flags]

Flags:
  -p, --cache-path string   Path to the local cache file (default "~/.zname.cache")
  -h, --help                Print this help and exit
  -r, --rebuild-cache       Rebuild the local cache
  -v, --version             Print version and exit
```
