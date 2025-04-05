package tasks

import (
	"github.com/suzmii/ACMBot/internal/context"
	"github.com/suzmii/ACMBot/pkg/model/provider"
)

func CodeforcesProfileHandler(ctx *context.Context) error {
	return context.NewChainContext(ctx).
		Then(getHandlerFromParams).
		Then(getCodeforcesUserByHandle).
		Then(getRenderedCodeforcesUserProfile).
		Then(sendPicture).
		Execute()
}

func CodeforcesRatingHandler(ctx *context.Context) error {
	return context.NewChainContext(ctx).
		Then(getHandlerFromParams).
		Then(getCodeforcesUserByHandle).
		Then(getRenderedCodeforcesRatingChanges).
		Then(sendPicture).
		Execute()
}

func AtcoderProfileHandler(ctx *context.Context) error {
	return context.NewChainContext(ctx).
		Then(getHandlerFromParams).
		Then(getAtcoderUserByHandle).
		Then(getRenderedAtcoderUserProfile).
		Then(sendPicture).
		Execute()
}

func RaceHandler(provider provider.RaceProvider) Task {
	return func(ctx *context.Context) error {
		ctx.StepValue = provider
		return context.NewChainContext(ctx).
			Then(getRaceFromProvider).
			Then(sendRace).
			Execute()
	}
}

func TextHandler(text string) Task {
	return func(ctx *context.Context) error {
		ctx.StepValue = text
		return sendText(ctx)
	}
}

func BindCodeforcesUserHandler(ctx *context.Context) error {
	return context.NewChainContext(ctx).
		Then(getHandlerFromParams).
		Then(bindCodeforcesUser).
		Execute()
}

func QQGroupRankHandler(ctx *context.Context) error {
	return qqGroupRank(ctx)
}
