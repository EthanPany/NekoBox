// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/wuhan005/NekoBox/internal/conf"
	"github.com/wuhan005/NekoBox/internal/context"
	"github.com/wuhan005/NekoBox/internal/db"
	"github.com/wuhan005/NekoBox/internal/form"
)

func Register(ctx context.Context) {
	ctx.Success("auth/register")
}

func RegisterAction(ctx context.Context, f form.Register) {
	if ctx.HasError() {
		ctx.Success("auth/register")
		return
	}

	if err := db.Users.Create(ctx.Request().Context(), db.CreateUserOptions{
		Name:       f.Name,
		Password:   f.Password,
		Email:      f.Email,
		Avatar:     conf.Upload.DefaultAvatarURL,
		Domain:     f.Domain,
		Background: conf.Upload.DefaultBackground,
		Intro:      "问你想问的",
	}); err != nil {
		switch {
		case errors.Is(err, db.ErrUserNotExists),
			errors.Is(err, db.ErrBadCredential),
			errors.Is(err, db.ErrDuplicateEmail),
			errors.Is(err, db.ErrDuplicateDomain):
			ctx.SetError(errors.Cause(err))

		default:
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to create new user")
			ctx.SetInternalError()
		}

		ctx.Success("auth/register")
		return
	}

	ctx.SetSuccessFlash("注册成功，欢迎来到 NekoBox！")
	ctx.Redirect("/login")
}
