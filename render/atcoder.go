package render

import (
	"context"
	"time"

	"github.com/suzmii/ACMBot/render/internal"
)

type AtcoderUserSolvedData struct {
	Range   string
	Percent float32
}

type AtcoderUserProfile struct {
	Avatar           string
	Handle           string
	MaxRating        int
	PromotionMessage string
	Rating           int
	Level            string
	Solved           int
	SolvedData       []AtcoderUserSolvedData
	Time             string
}

func AtcoderRating2Level(rating int, rank string) string {
	if rating <= 0 {
		if rank != "" {
			return "Unrated"
		}
		return "Unrated"
	}

	switch {
	case rating < 400:
		return "Gray"
	case rating < 800:
		return "Brown"
	case rating < 1200:
		return "Green"
	case rating < 1600:
		return "Cyan"
	case rating < 2000:
		return "Blue"
	case rating < 2400:
		return "Yellow"
	case rating < 2800:
		return "Orange"
	case rating < 3200:
		return "Red"
	default:
		return "Legend"
	}
}

func (r *Render) AtcoderProfile(ctx context.Context, user AtcoderUserProfile) ([]byte, error) {
	timestamp := user.Time
	if timestamp == "" {
		timestamp = time.Now().Format("2006-01-02 15:04:05")
	}

	return r.executeTemplate(
		ctx,
		internal.TemplateAtcoderProfile,
		&AtcoderUserProfile{
			Avatar:           user.Avatar,
			Handle:           user.Handle,
			MaxRating:        user.MaxRating,
			PromotionMessage: user.PromotionMessage,
			Rating:           user.Rating,
			Level:            user.Level,
			Solved:           user.Solved,
			SolvedData:       user.SolvedData,
			Time:             timestamp,
		},
	)
}
