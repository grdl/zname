package zname

import (
	"os"
	"path/filepath"
	"strings"
	"zname/cache"

	"github.com/jedib0t/go-pretty/v6/table"
)

type Config struct {
	RebuildCache bool
	Word         string
	CachePath    string
}

func (c *Config) validate() error {
	// if CachePath starts with '~', expand it to user home dir
	if c.CachePath[0] == '~' {
		home, err := os.UserHomeDir()
		if err != nil {
			return err
		}

		c.CachePath = filepath.Join(home, c.CachePath[1:])
	}

	return nil
}

type Zname struct {
	config *Config
}

func New(config *Config) (*Zname, error) {
	err := config.validate()
	if err != nil {
		return nil, err
	}

	return &Zname{
		config: config,
	}, nil
}

func (z *Zname) Run() error {
	c, err := cache.Open(z.config.CachePath)
	if err != nil {
		return err
	}

	err = z.rebuildCacheIfNeeded(c)
	if err != nil {
		return nil
	}

	records, err := c.FindByWord(z.config.Word)
	if err != nil {
		return err
	}

	printTable(records)

	return nil
}

func (z *Zname) rebuildCacheIfNeeded(c *cache.Cache) error {
	if z.config.RebuildCache {
		err := c.Rebuild()
		if err != nil {
			return err
		}
	}

	return nil
}

func printTable(records []cache.Record) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetStyle(table.StyleLight)
	t.Style().Options.SeparateColumns = false

	t.AppendHeader(table.Row{"Record", "Target"})

	for _, record := range records {
		t.AppendRow(table.Row{
			strings.TrimRight(record.Name, "."),
			strings.TrimRight(record.Target, "."),
		})
	}
	t.SortBy([]table.SortBy{
		{Name: "Record", Mode: table.Asc},
	})

	t.AppendSeparator()
	t.Render()
}
