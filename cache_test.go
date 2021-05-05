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
func TestFindByWord(t *testing.T) {
	zones := []Zone{
		{
			ID:   "1",
			Name: "example.com",
			Records: []Record{
				{
					Name:   "example.com",
					Target: "192.168.0.2",
					ZoneID: "1",
					Type:   "A",
				}, {
					Name:   "foo.example.com",
					Target: "someservice.local",
					ZoneID: "1",
					Type:   "A",
				}, {
					Name:   "bar.example.com",
					Target: "other-service.local",
					ZoneID: "1",
					Type:   "A",
				},
			},
		}, {
			ID:   "2",
			Name: "data.example.com",
			Records: []Record{
				{
					Name:   "data.example.com",
					Target: "127.0.0.1",
					ZoneID: "2",
					Type:   "A",
				}, {
					Name:   "foo.data.example.com",
					Target: "192.168.1.1",
					ZoneID: "2",
					Type:   "A",
				}, {
					Name:   "bar1.data.example.com",
					Target: "bar-service.local",
					ZoneID: "2",
					Type:   "A",
				},
			},
		},
	}

	tests := []struct {
		word string
		want []string
	}{
		{
			word: "foo",
			want: []string{
				"foo.data.example.com",
				"foo.example.com",
			},
		}, {
			word: "1",
			want: []string{
				"example.com",
				"data.example.com",
				"foo.data.example.com",
				"bar1.data.example.com",
			},
		}, {
			word: "xxx",
			want: []string{},
		},
	}

	for i, test := range tests {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			db, err := openDB("file::memory:?mode=memory")
			require.NoError(t, err)

			for _, zone := range zones {
				err := zone.Save(db)
				require.NoError(t, err)
			}

			found, err := FindByWord(db, test.word)
			require.NoError(t, err)

			assert.Len(t, test.want, len(found))
			for _, record := range found {
				assert.Contains(t, test.want, record.Name)
			}
		})
	}
}
