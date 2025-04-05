package tasks

import (
	"context"
	"github.com/suzmii/ACMBot/internal/util/ctxUtil"
	"github.com/suzmii/ACMBot/pkg/model/bot"
	"github.com/suzmii/ACMBot/pkg/model/provider"
)

func CodeforcesProfileHandler(ctx context.Context) error {
	return ctxUtil.NewChainContext(ctx).
		Then(getHandlerFromParams).
		Then(getCodeforcesUserByHandle).
		Then(getRenderedCodeforcesUserProfile).
		Then(sendPicture).
		Execute()
}

func CodeforcesRatingHandler(ctx context.Context) error {
	return ctxUtil.NewChainContext(ctx).
		Then(getHandlerFromParams).
		Then(getCodeforcesUserByHandle).
		Then(getRenderedCodeforcesRatingChanges).
		Then(sendPicture).
		Execute()
}

func AtcoderProfileHandler(ctx context.Context) error {
	return ctxUtil.NewChainContext(ctx).
		Then(getHandlerFromParams).
		Then(getAtcoderUserByHandle).
		Then(getRenderedAtcoderUserProfile).
		Then(sendPicture).
		Execute()
}

func RaceHandler(provider provider.RaceProvider) bot.Handler {
	return func(ctx context.Context) error {
		return ctxUtil.NewChainContext(ctxUtil.Set(ctx, provider)).
			Then(getRaceFromProvider).
			Then(sendRace).
			Execute()
	}
}

func TextHandler(text string) bot.Handler {
	return func(ctx context.Context) error {
		_, err := sendText(ctxUtil.Set(ctx, textMessage(text)))
		return err
	}
}

func BindCodeforcesUserHandler(ctx context.Context) error {
	return ctxUtil.NewChainContext(ctx).
		Then(getHandlerFromParams).
		Then(bindCodeforcesUser).
		Execute()
}

func QQGroupRankHandler(ctx context.Context) error {
	_, err := qqGroupRank(ctx)
	return err
}
