/*
 * Copyright (C) 2019 LambdaCal Inc.
 */

package middleware

import (
	"context"
	"strconv"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util"
)

// UserKey userKey
const UserKey = "user_id"

// UserNameKey UserNameKey
const UserNameKey = "user_name"

// GetUserID GetUserID
// ISSUE: The int64 id is only used in onpremise, platform is using snowflake ID.
// Future refactoring is needed.
func (mw *Middleware) GetUserID(ctx context.Context) (int64, error) {
	user, err := mw.GetUserIDString(ctx)
	if err != nil {
		return 0, err
	}
	return strconv.ParseInt(user, 10, 64)
}

// GetUserIDString GetUserIDString
func (mw *Middleware) GetUserIDString(ctx context.Context) (string, error) {
	return util.GetInMetadata(ctx, UserKey)
}

// WithUserID WithUserID
func (mw *Middleware) WithUserID(ctx context.Context, userID int64) context.Context {
	return util.SetInMetadata(ctx, UserKey, strconv.FormatInt(userID, 10))
}

// GetUserName GetUserName
func (mw *Middleware) GetUserName(ctx context.Context) (string, error) {
	return util.GetInMetadata(ctx, UserNameKey)
}
