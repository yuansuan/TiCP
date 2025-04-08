/*
 * // Copyright (C) 2018 LambdaCal Inc.
 *
 */

package validation_test

import (
	"log"
	"testing"

	"github.com/yuansuan/ticp/common/go-kit/gin-boot/util/validation"
)

func TestName(t *testing.T) {
	name := "hello)"
	log.Print(validation.Name(name))
}

func TestEmail(t *testing.T) {
	email := "name@yuansuan.cn"
	log.Print(validation.Email(email))
}
