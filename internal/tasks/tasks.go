package tasks

import (
	"fmt"
	"github.com/suzmii/ACMBot/internal/bot/message"
	"github.com/suzmii/ACMBot/internal/context"
	"github.com/suzmii/ACMBot/internal/manager"
	"github.com/suzmii/ACMBot/pkg/errs"
	"github.com/suzmii/ACMBot/pkg/model/provider"
	"github.com/suzmii/ACMBot/pkg/model/race"
)

type Task func(ctx *context.Context) error

// getHandlerFromParams nil -> []string
func getHandlerFromParams(ctx *context.Context) error {
	params := ctx.Params()
	var handles []string

	for _, handle := range params {
		for _, c := range handle {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '.' || c == '-') {
				return errs.ErrIllegalHandle
			}
		}
		handles = append(handles, handle)
	}

	ctx.StepValue = handles
	return nil
}

// getCodeforcesUserByHandle []string -> *manager.CodeforcesUser
func getCodeforcesUserByHandle(ctx *context.Context) error {
	handles, ok := ctx.StepValue.([]string)
	if !ok {
		return errs.NewInternalError("错误的参数类型")
	}

	if len(handles) == 0 {
		return errs.ErrNoHandle
	}

	if len(handles) > 1 {
		ctx.Send(message.Text("太多handle惹，我只查询`" + handles[0] + "`的哦"))
	}

	user, err := manager.GetUpdatedCodeforcesUser(handles[0])
	if err != nil {
		return err
	}

	ctx.StepValue = user
	return nil
}

// getRenderedCodeforcesUserProfile *manager.CodeforcesUser -> []byte
func getRenderedCodeforcesUserProfile(ctx *context.Context) error {
	user, ok := ctx.StepValue.(*manager.CodeforcesUser)
	if !ok {
		return errs.NewInternalError("错误的参数类型")
	}

	pic, err := user.ToRenderProfileV2().ToImage()
	if err != nil {
		return err
	}

	ctx.StepValue = pic
	return nil
}

// getRenderedCodeforcesRatingChanges *manager.CodeforcesUser -> []byte
func getRenderedCodeforcesRatingChanges(ctx *context.Context) error {
	user, ok := ctx.StepValue.(*manager.CodeforcesUser)
	if !ok {
		return errs.NewInternalError("错误的参数类型")
	}

	pic, err := user.ToRenderRatingChanges().ToImage()
	if err != nil {
		return err
	}

	ctx.StepValue = pic
	return nil
}

// getRaceFromProvider model.RaceProvider -> []model.Race
func getRaceFromProvider(ctx *context.Context) error {
	provider_, ok := ctx.StepValue.(provider.RaceProvider)
	if !ok {
		return errs.NewInternalError("错误的参数类型")
	}

	races, err := provider_()
	if err != nil {
		return err
	}

	ctx.StepValue = races
	return nil
}

// getAtcoderUserByHandle []string -> *manager.AtcoderUser
func getAtcoderUserByHandle(ctx *context.Context) error {
	handles, ok := ctx.StepValue.([]string)
	if !ok {
		return errs.NewInternalError("错误的参数类型")
	}

	if len(handles) == 0 {
		return errs.ErrNoHandle
	}

	if len(handles) > 1 {
		ctx.Send(message.Text("太多handle惹，我只查询`" + handles[0] + "`的哦"))
	}

	user, err := manager.GetUpdatedAtcoderUser(handles[0])
	if err != nil {
		return err
	}

	ctx.StepValue = user
	return nil
}

// getRenderedAtcoderUserProfile *manager.AtcoderUser -> []byte
func getRenderedAtcoderUserProfile(ctx *context.Context) error {
	user, ok := ctx.StepValue.(*manager.AtcoderUser)
	if !ok {
		return errs.NewInternalError("错误的参数类型")
	}

	pic, err := user.ToRenderProfile().ToImage()
	if err != nil {
		return err
	}

	ctx.StepValue = pic
	return nil
}

// bindCodeforcesUser []string -> nil
func bindCodeforcesUser(ctx *context.Context) error {
	if ctx.Platform != context.PlatformQQ {
		ctx.Send(message.Text(errs.ErrBadPlatform.Error()))
		return nil
	}

	handles, ok := ctx.StepValue.([]string)
	if !ok {
		return errs.NewInternalError("错误的参数类型")
	}

	if len(handles) != 1 {
		ctx.Send(message.Text(errs.ErrImDedicated.Error()))
		return nil
	}

	caller := ctx.GetCallerInfo()

	if caller.Group.ID == 0 {
		ctx.Send(message.Text(errs.ErrGroupOnly.Error()))
		return nil
	}

	var qqBind = manager.QQBind{
		QQGroupID:        uint64(caller.Group.ID),
		QQName:           caller.NickName,
		QID:              uint64(caller.ID),
		CodeforcesHandle: handles[0],
	}

	if err := manager.BindQQAndCodeforcesHandler(qqBind); err != nil {
		return err
	}

	ctx.Send(message.Text("绑定成功 " + caller.NickName + " -> " + handles[0]))

	ctx.StepValue = nil

	return nil
}

// qqGroupRankHandler nil -> nil
func qqGroupRank(ctx *context.Context) error {
	if ctx.Platform != context.PlatformQQ {
		ctx.Send(message.Text(errs.ErrBadPlatform.Error()))
		return nil
	}

	caller := ctx.GetCallerInfo()

	if caller.Group.ID == 0 {
		ctx.Send(message.Text(errs.ErrGroupOnly.Error()))
		return nil
	}

	group := manager.QQGroup{
		QQGroupName: caller.Group.Name,
		QQGroupID:   uint64(caller.Group.ID),
	}

	rank, err := manager.GetGroupRank(group)
	if err != nil {
		return errs.NewInternalError(err.Error())
	}

	msg := caller.Group.Name + "\n"
	for _, user := range rank.QQUsers {
		msg += fmt.Sprintf("#%d %s %d\n", user.RankInGroup, user.QName, user.CodeforcesRating)
	}

	ctx.Send(message.Text(msg))

	ctx.StepValue = nil

	return nil
}

// sendPicture []byte -> nil
func sendPicture(ctx *context.Context) error {
	pic, ok := ctx.StepValue.([]byte)
	if !ok {
		return errs.NewInternalError("错误的参数类型")
	}

	ctx.Send(message.Image(pic))
	ctx.StepValue = nil
	return nil
}

// sendRace []model.Race -> nil
func sendRace(ctx *context.Context) error {
	races, ok := ctx.StepValue.([]race.Race)
	if !ok {
		return errs.NewInternalError("错误的参数类型")
	}

	if len(races) == 0 {
		ctx.Send(message.Text("没有获取到相关数据..."))
		return nil
	}

	ctx.Send(message.Races(races))
	ctx.StepValue = nil
	return nil
}

func sendText(ctx *context.Context) error {
	text, ok := ctx.StepValue.(string)
	if !ok {
		return errs.NewInternalError("错误的参数类型")
	}
	ctx.Send(message.Text(text))
	ctx.StepValue = nil
	return nil
}
