package render

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/suzmii/ACMBot/config/subconfig"
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
	err = os.WriteFile("rank.png", image, 0644)
	require.NoError(t, err)
}
