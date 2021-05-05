package main

import (
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDBSaveAndLoad(t *testing.T) {
	tests := []struct {
		zones []Zone
	}{
		{
			zones: []Zone{
				{
					ID:      "id1",
					Name:    "name1",
					Records: []Record{},
				},
			},
		}, {
			zones: []Zone{
				{
					ID:      "id1",
					Name:    "name1",
					Records: []Record{},
				}, {
					ID:      "id2",
					Name:    "name2",
					Records: []Record{},
				},
			},
		}, {
			zones: []Zone{
				{
					ID:   "id1",
					Name: "name1",
					Records: []Record{
						{
							Name:   "rec1",
							Type:   "A",
							Target: "target1",
							ZoneID: "id1",
						}, {
							Name:   "rec2",
							Type:   "A",
							Target: "target2",
							ZoneID: "id1",
						}, {
							Name:   "rec3",
							Type:   "A",
							Target: "target3",
							ZoneID: "id1",
						},
					},
				}, {
					ID:      "id2",
					Name:    "name2",
					Records: []Record{},
				},
			},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			db, err := openDB("file::memory:?mode=memory")
			require.NoError(t, err)

			for _, zone := range test.zones {
				err := zone.Save(db)
				require.NoError(t, err)
			}

			found, err := FindAllZones(db)
			require.NoError(t, err)

			assert.Equal(t, test.zones, found)
		})
	}
}
