// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"path"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/wuhan005/NekoBox/internal/context"
	"github.com/wuhan005/NekoBox/internal/db"
	"github.com/wuhan005/NekoBox/internal/form"
)

func Login(ctx context.Context) {
	ctx.Success("auth/login")
}

func LoginAction(ctx context.Context, f form.Login) {
	if ctx.HasError() {
		ctx.Success("auth/login")
		return
	}

	uri := ctx.Request().Request.RequestURI // Keep the query when redirecting.

	user, err := db.Users.Authenticate(ctx.Request().Context(), f.Email, f.Password)
	if err != nil {
		if errors.Is(err, db.ErrBadCredential) {
			ctx.SetErrorFlash(errors.Cause(err).Error())
		} else {
			logrus.WithContext(ctx.Request().Context()).WithError(err).Error("Failed to authenticate user")
			ctx.SetInternalErrorFlash()
		}
		ctx.Redirect(uri)
		return
	}

	to := ctx.Query("to")
	to = path.Clean("/" + to)
	if to == "" {
		to = "/_/" + user.Domain
	}

	ctx.Session.Set("uid", user.ID)
	ctx.Redirect(to)
}
