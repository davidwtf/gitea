// Copyright 2018 The Gitea Authors. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package setting

import (
	"net/http"
	"testing"

	"code.gitea.io/gitea/models"
	"code.gitea.io/gitea/modules/auth"
	"code.gitea.io/gitea/modules/setting"
	"code.gitea.io/gitea/modules/test"

	"github.com/stretchr/testify/assert"
)

func TestChangePassword(t *testing.T) {
	oldPassword := "password"
	setting.MinPasswordLength = 6
	setting.PasswordComplexity = map[string]string{
		"lower": "[a-z]+",
		"upper": "[A-Z]+",
		"digit": "[0-9]+",
		"spec":  "[-_]+",
	}
	var pcLUN = map[string]string{
		"lower": "[a-z]+",
		"upper": "[A-Z]+",
		"digit": "[0-9]+",
	}
	var pcLU = map[string]string{
		"lower": "[a-z]+",
		"upper": "[A-Z]+",
	}

	for _, req := range []struct {
		OldPassword        string
		NewPassword        string
		Retype             string
		Message            string
		PasswordComplexity map[string]string
	}{
		{
			OldPassword:        oldPassword,
			NewPassword:        "Qwerty123456-",
			Retype:             "Qwerty123456-",
			Message:            "",
			PasswordComplexity: setting.PasswordComplexity,
		},
		{
			OldPassword:        oldPassword,
			NewPassword:        "12345",
			Retype:             "12345",
			Message:            "auth.password_too_short",
			PasswordComplexity: setting.PasswordComplexity,
		},
		{
			OldPassword:        "12334",
			NewPassword:        "123456",
			Retype:             "123456",
			Message:            "settings.password_incorrect",
			PasswordComplexity: setting.PasswordComplexity,
		},
		{
			OldPassword:        oldPassword,
			NewPassword:        "123456",
			Retype:             "12345",
			Message:            "form.password_not_match",
			PasswordComplexity: setting.PasswordComplexity,
		},
		{
			OldPassword:        oldPassword,
			NewPassword:        "Qwerty",
			Retype:             "Qwerty",
			Message:            "settings.password_complexity",
			PasswordComplexity: setting.PasswordComplexity,
		},
		{
			OldPassword:        oldPassword,
			NewPassword:        "Qwerty",
			Retype:             "Qwerty",
			Message:            "settings.password_complexity",
			PasswordComplexity: pcLUN,
		},
		{
			OldPassword:        oldPassword,
			NewPassword:        "QWERTY",
			Retype:             "QWERTY",
			Message:            "settings.password_complexity",
			PasswordComplexity: pcLU,
		},
	} {
		models.PrepareTestEnv(t)
		ctx := test.MockContext(t, "user/settings/security")
		test.LoadUser(t, ctx, 2)
		test.LoadRepo(t, ctx, 1)

		AccountPost(ctx, auth.ChangePasswordForm{
			OldPassword: req.OldPassword,
			Password:    req.NewPassword,
			Retype:      req.Retype,
		})

		assert.EqualValues(t, req.Message, ctx.Flash.ErrorMsg)
		assert.EqualValues(t, http.StatusFound, ctx.Resp.Status())
	}
}
