package render

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/suzmii/ACMBot/config/subconfig"
	"github.com/suzmii/ACMBot/testutil"
)

var render *Render

func TestMain(m *testing.M) {
	var err error
	render, err = NewRender(subconfig.Render{
		Headless: false,
		PoolSize: 1,
	})
	if err != nil {
		panic(err)
	}
	os.Exit(m.Run())
}

func TestRank(t *testing.T) {
	t.Log("start render")
	image, err := render.Rank(t.Context(), []RankUser{
		{
			Rank:   1,
			Handle: "Tourist",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 4000,
			Level:  string(Rating2Level(4000)),
		},
		{
			Rank:   2,
			Handle: "LGM",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 3500,
			Level:  string(Rating2Level(3500)),
		},
		{
			Rank:   3,
			Handle: "IGM",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 2800,
			Level:  string(Rating2Level(2800)),
		},
		{
			Rank:   4,
			Handle: "GM",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 2500,
			Level:  string(Rating2Level(2500)),
		},
		{
			Rank:   5,
			Handle: "IM",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 2350,
			Level:  string(Rating2Level(2350)),
		},
		{
			Rank:   6,
			Handle: "Master",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 2200,
			Level:  string(Rating2Level(2200)),
		},
		{
			Rank:   7,
			Handle: "CM",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 2000,
			Level:  string(Rating2Level(2000)),
		},
		{
			Rank:   8,
			Handle: "Expert",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 1800,
			Level:  string(Rating2Level(1800)),
		},
		{
			Rank:   9,
			Handle: "Specialist",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 1500,
			Level:  string(Rating2Level(1500)),
		},
		{
			Rank:   10,
			Handle: "Pupil",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 1300,
			Level:  string(Rating2Level(1300)),
		},
		{
			Rank:   11,
			Handle: "Newbie",
			Avatar: "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
			Rating: 1000,
			Level:  string(Rating2Level(1000)),
		},
	})
	require.NoError(t, err)
	testutil.WriteArtifact(t, []string{"..", ".cache", "render-tests", "rank.png"}, image)
}

func TestProfileV2(t *testing.T) {
	image, err := render.ProfileV2(t.Context(), CodeforcesUserProfile{
		Avatar:    "https://userpic.codeforces.org/3363857/title/dfaf69f954d867bd.jpg",
		Handle:    "tourist.v2",
		MaxRating: 3979,
		FriendOf:  78231,
		Rating:    3900,
		Level:     Rating2Level(3900),
		Solved:    4287,
		SolvedData: []CodeforcesUserSolvedData{
			{Range: "800+", Percent: 38.4},
			{Range: "1400+", Percent: 32.1},
			{Range: "2000+", Percent: 20.8},
			{Range: "2600+", Percent: 8.7},
		},
	})
	require.NoError(t, err)
	testutil.WriteArtifact(t, []string{"..", ".cache", "render-tests", "codeforces_profile_v2.png"}, image)
}

func TestAtcoderProfile(t *testing.T) {
	image, err := render.AtcoderProfile(t.Context(), AtcoderUserProfile{
		Avatar:           "https://img.atcoder.jp/assets/top/img/logo_bk.svg",
		Handle:           "rng_58",
		MaxRating:        4123,
		PromotionMessage: "12x Heuristic",
		Rating:           4011,
		Level:            AtcoderRating2Level(4011, "1"),
		Solved:           1934,
		SolvedData: []AtcoderUserSolvedData{
			{Range: "~399", Percent: 14.2},
			{Range: "400~799", Percent: 27.8},
			{Range: "800~1199", Percent: 25.4},
			{Range: "1200~1599", Percent: 18.6},
			{Range: "1600~1999", Percent: 9.1},
			{Range: "2000+", Percent: 4.9},
		},
		Time: time.Now().Format("2006-01-02 15:04:05"),
	})
	require.NoError(t, err)
	testutil.WriteArtifact(t, []string{"..", ".cache", "render-tests", "atcoder_profile.png"}, image)
}
