package embedstore

import (
	"testing"

	"github.com/tzapio/tzap/pkg/types"
)

func TestSortSearchResults(t *testing.T) {
	searchResults := []types.SearchResult{
		{
			Vector: types.Vector{
				ID:        "1",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file1",
					Start:    500,
					End:      750,
				},
			},
		},
		{
			Vector: types.Vector{
				ID:        "2",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file1",
					Start:    100,
					End:      350,
				},
			},
		},
		{
			Vector: types.Vector{
				ID:        "3",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file1",
					Start:    300,
					End:      550,
				},
			},
		},
		{
			Vector: types.Vector{
				ID:        "4",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file2",
					Start:    100,
					End:      350,
				},
			},
		},
	}

	want := []types.SearchResult{
		{
			Vector: types.Vector{
				ID:        "1",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file1",
					Start:    100,
					End:      750,
				},
			},
		},
		{
			Vector: types.Vector{
				ID:        "4",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file2",
					Start:    100,
					End:      350,
				},
			},
		},
	}

	got := TightenSearchResults(searchResults)

	if len(got.Results) != len(want) {
		t.Errorf("Expected %d sorted search results, got %d", len(want), len(got.Results))
	}

	for i := range want {
		if want[i].Vector.Metadata.Filename != got.Results[i].Vector.Metadata.Filename ||
			want[i].Vector.Metadata.Start != got.Results[i].Vector.Metadata.Start ||
			want[i].Vector.Metadata.End != got.Results[i].Vector.Metadata.End {
			t.Errorf("Expected search result %+v, got %+v", want[i].Vector.Metadata, got.Results[i].Vector.Metadata)
		}
	}
}

func TestGroupConsecutiveMetadata(t *testing.T) {
	searchResults := []types.SearchResult{
		{
			Vector: types.Vector{
				ID:        "1",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file1",
					Start:    100,
					End:      350,
				},
			},
		},
		{
			Vector: types.Vector{
				ID:        "2",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file1",
					Start:    300,
					End:      550,
				},
			},
		},
		{
			Vector: types.Vector{
				ID:        "3",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file1",
					Start:    500,
					End:      750,
				},
			},
		},
		{
			Vector: types.Vector{
				ID:        "4",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file1",
					Start:    700,
					End:      950,
				},
			},
		},
	}

	want := []types.SearchResult{
		{
			Vector: types.Vector{
				ID:        "",
				TimeStamp: 0,
				Metadata: types.Metadata{
					Filename: "file1",
					Start:    100,
					End:      950,
				},
			},
		},
	}

	got := groupConsecutiveMetadata(searchResults)

	if len(got) != len(want) {
		t.Errorf("Expected %d consecutive metadata groups, got %d", len(want), len(got))
	}

	for i := range want {
		if want[i].Vector.Metadata.Filename != got[i].Vector.Metadata.Filename ||
			want[i].Vector.Metadata.Start != got[i].Vector.Metadata.Start ||
			want[i].Vector.Metadata.End != got[i].Vector.Metadata.End {
			t.Errorf("Expected metadata group %+v, got %+v", want[i].Vector.Metadata, got[i].Vector.Metadata)
		}
	}
}
