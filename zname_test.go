package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(t *testing.T) {
	client, err := New()
	require.NoError(t, err)

	zones, err := client.GetZones()
	require.NoError(t, err)

	for _, zone := range zones {
		fmt.Printf("Scraping %s zone...\n", *zone.Name)

		records, err := client.GetRecords(*zone.Id)
		require.NoError(t, err)

		fmt.Printf("\tFound %d records\n", len(records))

	}

}
