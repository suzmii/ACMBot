package tasks

import (
	"context"
	"fmt"
	"github.com/suzmii/ACMBot/internal/bot/message"
	"github.com/suzmii/ACMBot/internal/logic/manager"
	"github.com/suzmii/ACMBot/internal/util/ctxUtil"
	"github.com/suzmii/ACMBot/pkg/errs"
	"github.com/suzmii/ACMBot/pkg/model/bot"
)

// getHandlerFromParams nil -> handles
func getHandlerFromParams(ctx context.Context) (context.Context, error) {
	params, ok := ctxUtil.Get[Params](ctx)
	if !ok {
		return ctx, errs.NewInternalError("missing Params in ctx")
	}
	var handles_ handles

	for _, handle := range params {
		for _, c := range handle {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '_' || c == '.' || c == '-') {
				return ctx, errs.ErrIllegalHandle
			}
		}
		handles_ = append(handles_, handle)
	}

	return ctxUtil.Set[handles](ctx, handles_), nil
}

// getCodeforcesUserByHandle []string -> *manager.CodeforcesUser
func getCodeforcesUserByHandle(ctx context.Context) (context.Context, error) {
	api := ctxUtil.MustGet[apiCaller](ctx)

	handles, ok := ctxUtil.Get[handles](ctx)
	if !ok {
		return ctx, errs.NewInternalError("missing handles in ctx")
	}

	if len(handles) == 0 {
		return ctx, errs.ErrNoHandle
	}

	if len(handles) > 1 {
		api.Send(message.Text("太多handle惹，我只查询`" + handles[0] + "`的哦"))
	}
	user, err := manager.GetUpdatedCodeforcesUser(handles[0])
	if err != nil {
		return ctx, err
	}

	return ctxUtil.Set[codeforcesUser](ctx, user), nil
}

// getRenderedCodeforcesUserProfile *manager.CodeforcesUser -> []byte
func getRenderedCodeforcesUserProfile(ctx context.Context) (context.Context, error) {
	user, ok := ctxUtil.Get[codeforcesUser](ctx)
	if !ok {
		return ctx, errs.NewInternalError("missing codeforces user in ctx")
	}

	pic, err := user.ToRenderProfileV2().ToImage()
	if err != nil {
		return ctx, err
	}

	return ctxUtil.Set[picMessage](ctx, pic), nil
}

// getRenderedCodeforcesRatingChanges *manager.CodeforcesUser -> []byte
func getRenderedCodeforcesRatingChanges(ctx context.Context) (context.Context, error) {
	user, ok := ctxUtil.Get[codeforcesUser](ctx)
	if !ok {
		return ctx, errs.NewInternalError("missing codeforces user in ctx")
	}

	pic, err := user.ToRenderRatingChanges().ToImage()
	if err != nil {
		return ctx, err
	}

	return ctxUtil.Set[picMessage](ctx, pic), nil
}

// getRaceFromProvider model.RaceProvider -> []model.Race
func getRaceFromProvider(ctx context.Context) (context.Context, error) {
	provider, ok := ctxUtil.Get[raceProvider](ctx)
	if !ok {
		return ctx, errs.NewInternalError("missing race provider in ctx")
	}

	races_, err := provider()
	if err != nil {
		return ctx, err
	}

	return ctxUtil.Set[races](ctx, races_), nil
}

// getAtcoderUserByHandle []string -> *manager.AtcoderUser
func getAtcoderUserByHandle(ctx context.Context) (context.Context, error) {
	api := ctxUtil.MustGet[apiCaller](ctx)

	handles, ok := ctxUtil.Get[handles](ctx)
	if !ok {
		return ctx, errs.NewInternalError("missing handles in ctx")
	}

	if len(handles) == 0 {
		return ctx, errs.ErrNoHandle
	}

	if len(handles) > 1 {
		api.Send(message.Text("太多handle惹，我只查询`" + handles[0] + "`的哦"))
	}

	user, err := manager.GetUpdatedAtcoderUser(handles[0])
	if err != nil {
		return ctx, err
	}

	return ctxUtil.Set[atcoderUser](ctx, user), nil
}

// getRenderedAtcoderUserProfile *manager.AtcoderUser -> []byte
func getRenderedAtcoderUserProfile(ctx context.Context) (context.Context, error) {
	user, ok := ctxUtil.Get[atcoderUser](ctx)
	if !ok {
		return ctx, errs.NewInternalError("错误的参数类型")
	}

	pic, err := user.ToRenderProfile().ToImage()
	if err != nil {
		return ctx, err
	}

	return ctxUtil.Set[picMessage](ctx, pic), nil
}

// bindCodeforcesUser []string -> nil
func bindCodeforcesUser(ctx context.Context) (context.Context, error) {
	platform := ctxUtil.MustGet[platform](ctx)
	api := ctxUtil.MustGet[apiCaller](ctx)

	if platform != bot.PlatformQQ {
		api.Send(message.Text(errs.ErrBadPlatform.Error()))
		return ctx, nil
	}

	handles, ok := ctxUtil.Get[handles](ctx)
	if !ok {
		return ctx, errs.NewInternalError("missing handles in ctx")
	}

	if len(handles) != 1 {
		api.Send(message.Text(errs.ErrImDedicated.Error()))
		return ctx, nil
	}

	caller := api.GetCallerInfo()

	if caller.Group.ID == 0 {
		api.Send(message.Text(errs.ErrGroupOnly.Error()))
		return ctx, nil
	}

	var qqBind = manager.QQBind{
		QQGroupID:        uint64(caller.Group.ID),
		QQName:           caller.NickName,
		QID:              uint64(caller.ID),
		CodeforcesHandle: handles[0],
	}

	if err := manager.BindQQAndCodeforcesHandler(qqBind); err != nil {
		return ctx, err
	}

	api.Send(message.Text("绑定成功 " + caller.NickName + " -> " + handles[0]))

	return ctx, nil
}

// qqGroupRankHandler nil -> nil
func qqGroupRank(ctx context.Context) (context.Context, error) {
	platform := ctxUtil.MustGet[platform](ctx)
	api := ctxUtil.MustGet[apiCaller](ctx)

	if platform != bot.PlatformQQ {
		api.Send(message.Text(errs.ErrBadPlatform.Error()))
		return ctx, nil
	}

	caller := api.GetCallerInfo()

	if caller.Group.ID == 0 {
		api.Send(message.Text(errs.ErrGroupOnly.Error()))
		return ctx, nil
	}

	group := manager.QQGroup{
		QQGroupName: caller.Group.Name,
		QQGroupID:   uint64(caller.Group.ID),
	}

	rank, err := manager.GetGroupRank(group)
	if err != nil {
		return ctx, errs.NewInternalError(err.Error())
	}

	msg := caller.Group.Name + "\n"
	for _, user := range rank.QQUsers {
		msg += fmt.Sprintf("#%d %s %d\n", user.RankInGroup, user.QName, user.CodeforcesRating)
	}

	api.Send(message.Text(msg))

	return ctx, nil
}

// sendPicture []byte -> nil
func sendPicture(ctx context.Context) (context.Context, error) {
	api := ctxUtil.MustGet[apiCaller](ctx)

	pic, ok := ctxUtil.Get[picMessage](ctx)
	if !ok {
		return ctx, errs.NewInternalError("missing picMessage in ctx")
	}

	api.Send(message.Image(pic))
	return ctx, nil
}

// sendRace []model.Race -> nil
func sendRace(ctx context.Context) (context.Context, error) {
	api := ctxUtil.MustGet[apiCaller](ctx)

	races, ok := ctxUtil.Get[races](ctx)
	if !ok {
		return ctx, errs.NewInternalError("missing races in ctx")
	}

	if len(races) == 0 {
		api.Send(message.Text("没有获取到相关数据..."))
		return ctx, nil
	}

	api.Send(message.Races(races))
	return ctx, nil
}

func sendText(ctx context.Context) (context.Context, error) {
	api := ctxUtil.MustGet[apiCaller](ctx)

	text, ok := ctxUtil.Get[textMessage](ctx)
	if !ok {
		return ctx, errs.NewInternalError("missing text in ctx")
	}
	api.Send(message.Text(text))

	return ctx, nil
}
