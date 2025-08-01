package handler

import (
	"github.com/sirupsen/logrus"
	"github.com/suzmii/ACMBot/internal/database/repo"
	"github.com/suzmii/ACMBot/internal/datasync"
	"github.com/suzmii/ACMBot/internal/model"
	"github.com/suzmii/ACMBot/internal/render"
	"github.com/suzmii/ACMBot/internal/render/rendermodel"
	"github.com/suzmii/ACMBot/internal/util/param"
	"time"
)

func RatingDetailHandler(c *model.Context) error {
	username, err := param.AsCodeforcesUsername(c.Params)
	if err != nil {
		return err
	}

	user, err := datasync.GetOrSyncUser(c.Ctx, username)
	if err != nil {
		return err
	}

	ratingRecords, err := repo.GetRatingRecords(c.Ctx, user.ID)
	if err != nil {
		return err
	}

	var records model.RatingRecords
	logrus.Debug(ratingRecords)

	if ratingRecords != nil && time.Now().Sub(ratingRecords.UpdatedAt) < 4*time.Hour {
		records = model.RatingRecordsFromRepo(*ratingRecords)
	} else {
		records, err = datasync.RatingRecords(c.Ctx, user)
		if err != nil {
			return err
		}
	}

	start := time.Now()
	image, err := render.RatingDetail(rendermodel.CodeforcesRatingRecordsFromCommon(username, records))
	if err != nil {
		return err
	}
	end := time.Now()
	logrus.Debug("render cost: ", end.Sub(start))

	c.Messager.Send(model.ImageMessage{Image: image})
	return nil
}

func ProfileHandler(c *model.Context) error {
	start := time.Now()
	username, err := param.AsCodeforcesUsername(c.Params)
	if err != nil {
		return err
	}

	user, err := datasync.GetOrSyncUser(c.Ctx, username)
	if err != nil {
		return err
	}

	if time.Now().Sub(user.UpdatedAt) > 4*time.Hour {
		user, err = datasync.User(c.Ctx, username)
		if err != nil {
			return err
		}
	}

	if time.Now().Sub(user.SubmissionUpdatedAt) > 24*time.Hour {
		err = datasync.Submissions(c.Ctx, user)
		if err != nil {
			return err
		}
	}

	detail, err := repo.CountSolvedDetail(c.Ctx, user.ID)
	if err != nil {
		return err
	}

	end := time.Now()
	logrus.Debug("query cost:", end.Sub(start))

	start = time.Now()
	image, err := render.ProfileV2(rendermodel.CodeforcesUserProfileFromCommon(user, *detail))
	if err != nil {
		return err
	}
	end = time.Now()
	logrus.Debug("render cost:", end.Sub(start))

	start = time.Now()
	c.Messager.Send(model.ImageMessage{Image: image})
	end = time.Now()
	logrus.Debug("send message cost:", end.Sub(start))
	return nil
}
