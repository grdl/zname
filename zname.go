package zname

import (
	"os"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
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

	printTable(records)

	return nil
}

func printTable(records []Record) {
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
