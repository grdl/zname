package zname

import (
	"fmt"
)

type Config struct {
	RebuildCache bool
	Word         string
	CachePath    string
}

func Run(cfg *Config) error {
	if cfg.RebuildCache {
		err := RebuildCache(cfg.CachePath)
		if err != nil {
			return err
		}
	}

	if cfg.Word == "" {
		return nil
	}

	cache, err := OpenCache(cfg.CachePath)
	if err != nil {
		return err
	}

	records, err := FindByWord(cache, cfg.Word)
	if err != nil {
		return err
	}

	fmt.Println(records)

	return nil
}
