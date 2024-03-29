package dtos

import "github.com/labstack/echo"

type MemberSignIn struct {
	Id       string `json:"id" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (m MemberSignIn) Validate(ctx echo.Context) error {
	return ctx.Validate(m)
}

type DoorayMember struct {
	Id                   string `json:"id"`
	UserCode             string `json:"userCode"`
	Name                 string `json:"name"`
	ExternalEmailAddress string `json:"externalEmailAddress"`
}
