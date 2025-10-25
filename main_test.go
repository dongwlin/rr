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
		{
			Raw:             "[BeanSub&FZSD][Chainsaw_Man][02][GB][1080P][x264_AAC].mp4",
			ExpectedEpisode: 2,
		},
		{
			Raw:             "S01E07",
			ExpectedEpisode: 7,
		},
		{
			Raw:             "EP12 - The Final Battle",
			ExpectedEpisode: 12,
		},
		{
			Raw:             "第03集 - 新的开始",
			ExpectedEpisode: 3,
		},
		{
			Raw:             "No episode info here",
			ExpectedEpisode: -1,
		},
	}

	for _, tc := range testCases {

		ep := extractEpisode(tc.Raw)
		require.Equal(t, tc.ExpectedEpisode, ep)
	}

}
