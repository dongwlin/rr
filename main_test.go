package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestExtractEpisode(t *testing.T) {

	type Case struct {
		Raw             string
		ExpectedEpisode int
	}

	testCases := []Case{
		{
			Raw:             "[Lilith-Raws] Kubo-san wa Mob wo Yurusanai - 04 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4",
			ExpectedEpisode: 4,
		},
	}

	for _, tc := range testCases {

		ep := extractEpisode(tc.Raw)
		require.Equal(t, tc.ExpectedEpisode, ep)
	}

}
