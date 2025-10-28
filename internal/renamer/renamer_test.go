package renamer

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

func TestPadNumber(t *testing.T) {
	type Case struct {
		Num      string
		Width    int
		Expected string
	}

	cases := []Case{
		{Num: "5", Width: 2, Expected: "05"},
		{Num: "12", Width: 2, Expected: "12"},
		{Num: "123", Width: 2, Expected: "123"},
		{Num: "", Width: 3, Expected: "000"},
		{Num: "0", Width: 2, Expected: "00"},
		{Num: "7", Width: 1, Expected: "7"},
		{Num: "7", Width: 0, Expected: "7"},
		{Num: "007", Width: 5, Expected: "00007"},
		{Num: "9", Width: -3, Expected: "9"},
	}

	for _, tc := range cases {
		got := padNumber(tc.Num, tc.Width)
		require.Equal(t, tc.Expected, got, "num=%q width=%d", tc.Num, tc.Width)
	}
}

func TestExtractOtherTags(t *testing.T) {
	type Case struct {
		Raw      string
		Expected []string
	}

	cases := []Case{
		{
			Raw:      "[Lilith-Raws] Kubo-san wa Mob wo Yurusanai - 04 [Baha][WEB-DL][1080p][AVC AAC][CHT][MP4].mp4",
			Expected: []string{"[Lilith-Raws]", "[Baha]", "[WEB-DL]", "[1080p]", "[AVC AAC]", "[CHT]", "[MP4]"},
		},
		{
			Raw:      "[BeanSub&FZSD][Chainsaw_Man][02][GB][1080P][x264_AAC].mp4",
			Expected: []string{"[BeanSub&FZSD]", "[Chainsaw_Man]", "[02]", "[GB]", "[1080P]", "[x264_AAC]"},
		},
		{
			Raw:      "S01E07",
			Expected: []string{},
		},
		{
			Raw:      "",
			Expected: []string{},
		},
	}

	for _, tc := range cases {
		got := extractOtherTags(tc.Raw)
		require.Equal(t, tc.Expected, got, "raw=%q", tc.Raw)
	}
}
