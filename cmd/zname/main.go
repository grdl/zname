package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"zname"
)

const (
	defaultCacheFile = ".zname.cache"
	cachePathEnv     = "ZNAME_CACHE"
)

var rebuildCache = flag.Bool("r", false, "Rebuild cache")

func main() {
	cfg := parseFlags()

	if err := zname.Run(cfg); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func parseFlags() *zname.Config {
	flag.Parse()

	var cachePath string
	if path, ok := os.LookupEnv(cachePathEnv); !ok {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			cachePath = defaultCacheFile
		} else {
			cachePath = filepath.Join(homeDir, defaultCacheFile)
		}
	} else {
		cachePath = path
	}

	return &zname.Config{
		RebuildCache: *rebuildCache,
		Word:         flag.Arg(0),
		CachePath:    cachePath,
	}
}
