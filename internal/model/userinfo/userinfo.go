package userinfo

import "github.com/suzmii/ACMBot/internal/renderer"

type UserInfo struct {
	Profile Profile
	Rating  Rating
}

type Profile interface {
	ToProfile() renderer.RenderAble
}

type Rating interface {
	ToRating() renderer.RenderAble
}
