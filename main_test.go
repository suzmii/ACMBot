package main_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/stretchr/testify/require"
	pkgapi "github.com/suzmii/ACMBot/api"
	"github.com/suzmii/ACMBot/config"
	"github.com/suzmii/ACMBot/consts"
	"github.com/suzmii/ACMBot/database"
	pkghandler "github.com/suzmii/ACMBot/handler"
	pkgrender "github.com/suzmii/ACMBot/render"
	"github.com/suzmii/ACMBot/util/logx"
)

var (
	err     error
	cfg     *config.Config
	api     *pkgapi.API
	store   database.Store
	render  *pkgrender.Render
	handler *pkghandler.Handler
)

func TestMain(m *testing.M) {
	cfg = config.LoadConfig()
	cfg.Logger.Level = "trace"
	logx.Init(cfg.Logger)
	api, err = pkgapi.New(cfg.API)
	if err != nil {
		panic(err)
	}
	store, err = database.New(context.Background(), cfg.Database.DSN)
	if err != nil {
		panic(err)
	}
	render, err = pkgrender.NewRender(cfg.Render)
	if err != nil {
		panic(err)
	}
	handler = pkghandler.NewHandler(cfg.Handler, store, api, render)
	os.Exit(m.Run())
}

func TestGetRacesFromDB(t *testing.T) {
	for _, i := range []string{
		consts.RaceResourceCodeforces,
		consts.RaceResourceAtcoder,
		consts.RaceResourceLeetcode,
		consts.RaceResourceLuogu,
		consts.RaceResourceNowcoder,
	} {
		t.Log(i)
		races, err := handler.GetUpcomingRace(t.Context(), i)
		require.NoError(t, err)
		spew.Dump(races)
	}

}

var testUsername = "2c2048d2"

func TestCodeforcesRatingRecords(t *testing.T) {
	bytes, err := handler.GetCodeforcesRatingImage(t.Context(), testUsername)
	require.NoError(t, err)
	err = os.WriteFile(fmt.Sprintf("codeforces_rating_%s.png", testUsername), bytes, 0644)
	require.NoError(t, err)
}

func TestCodeforcesProfile(t *testing.T) {
	bytes, err := handler.GetCodeforcesUserProfileImage(t.Context(), testUsername)
	require.NoError(t, err)
	err = os.WriteFile(fmt.Sprintf("codeforces_profile_%s.png", testUsername), bytes, 0644)
	require.NoError(t, err)
}

func BenchmarkCodeforcesRatingRecords(b *testing.B) {
	ctx := context.Background()

	_, err := handler.GetCodeforcesRatingImage(ctx, testUsername)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := handler.GetCodeforcesRatingImage(ctx, testUsername)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkCodeforcesProfile(b *testing.B) {
	ctx := context.Background()

	_, err := handler.GetCodeforcesUserProfileImage(ctx, testUsername)
	require.NoError(b, err)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := handler.GetCodeforcesUserProfileImage(ctx, testUsername)
		if err != nil {
			b.Fatal(err)
		}
	}
}
