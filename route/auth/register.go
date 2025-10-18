// Copyright 2022 E99p1ant. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package auth

import (
	"fmt"
	"net/url"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/wuhan005/NekoBox/internal/conf"
	"github.com/wuhan005/NekoBox/internal/context"
	"github.com/wuhan005/NekoBox/internal/db"
	"github.com/wuhan005/NekoBox/internal/form"
	"github.com/wuhan005/NekoBox/internal/storage"
)

func Register(ctx context.Context) {
	ctx.Success("auth/register")
}

func RegisterAction(ctx context.Context, f form.Register) {
	if ctx.HasError() {
		ctx.Success("auth/register")
		return
	}

	// Check recaptcha code if enabled.
	// Note: When recaptcha is disabled, we skip this check entirely.
	// When enabled, the recaptcha middleware must be registered to inject the service.
	if conf.Security.EnableRecaptcha && f.Recaptcha == "" {
		ctx.SetErrorFlash("验证码错误")
		ctx.Redirect("/register")
		return
	}

	// Generate unique avatar by downloading from DiceBear and uploading to R2
	avatarURL := conf.Upload.DefaultAvatarURL
	if avatarURL != "" {
		// Generate DiceBear URL with username + "nekobox" as seed
		seed := url.QueryEscape(f.Name + "nekobox")
		separator := "?"
		// Check if URL already has query parameters
		for i := 0; i < len(avatarURL); i++ {
			if avatarURL[i] == '?' {
				separator = "&"
				break
			}
		}
		diceBearURL := fmt.Sprintf("%s%sseed=%s", avatarURL, separator, seed)

		// Download the generated avatar and upload to R2
		r2URL, err := storage.DownloadAndUploadToS3(diceBearURL, ".svg")
		if err != nil {
			// If download/upload fails, log error but continue with DiceBear URL as fallback
			logrus.WithContext(ctx.Request().Context()).WithError(err).Warn("Failed to download and upload avatar to R2, using DiceBear URL as fallback")
			avatarURL = diceBearURL
		} else {
			// Successfully uploaded to R2, use the R2 URL
			avatarURL = r2URL
		}
	}

	if err := db.Users.Create(ctx.Request().Context(), db.CreateUserOptions{
		Name:       f.Name,
		Password:   f.Password,
		Email:      f.Email,
		Avatar:     avatarURL,
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
