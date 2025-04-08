package oshelp

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/yuansuan/ticp/iPaaS/standard-compute/internal/systemuser"
	"github.com/yuansuan/ticp/iPaaS/standard-compute/pkg/xio"
)

type config struct {
	username string
	fileMode os.FileMode
}

type Option interface {
	apply(c *config)
}

type optionFunc func(c *config)

func (f optionFunc) apply(c *config) {
	f(c)
}

func WithChown(username string) Option {
	return optionFunc(func(c *config) {
		c.username = username
	})
}

func WithChmod(fileMode os.FileMode) Option {
	return optionFunc(func(c *config) {
		c.fileMode = fileMode
	})
}

func Write(file *os.File, content []byte, opts ...Option) (int, error) {
	c := &config{}
	for _, opt := range opts {
		opt.apply(c)
	}

	var err error
	if err = preFileWrite(file, c); err != nil {
		return 0, err
	}

	n, err := file.Write(content)
	if err != nil {
		return 0, err
	}

	return n, nil
}

func CopyToFile(ctx context.Context, file *os.File, src io.Reader, opts ...Option) error {
	c := &config{}
	for _, opt := range opts {
		opt.apply(c)
	}

	var err error
	if err = preFileWrite(file, c); err != nil {
		return err
	}

	if _, err = xio.Copy(ctx, file, src); err != nil {
		return err
	}

	return nil
}

func preFileWrite(file *os.File, c *config) error {
	if c.username != "" {
		user, err := systemuser.Get(c.username)
		if err != nil {
			return fmt.Errorf("get user info from cache failed, %w", err)
		}

		if err = file.Chown(user.Uid, user.Gid); err != nil {
			return fmt.Errorf("chown %s failed, %w", file.Name(), err)
		}
	}

	if c.fileMode != os.FileMode(0) {
		if err := file.Chmod(c.fileMode); err != nil {
			return fmt.Errorf("chmod %s failed, %w", file.Name(), err)
		}
	}

	return nil
}

// Mkdir chmod or chown only do at nearest dir
func Mkdir(dir string, opts ...Option) error {
	c := &config{}
	for _, opt := range opts {
		opt.apply(c)
	}

	var err error
	if err = os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("mkdirAll %s failed, %w", dir, err)
	}

	if err = postMkdir(dir, c); err != nil {
		return err
	}

	return nil
}

func postMkdir(dir string, c *config) error {
	if c.username != "" {
		user, err := systemuser.Get(c.username)
		if err != nil {
			return fmt.Errorf("get user info from cache failed, %w", err)
		}

		if err = os.Chown(dir, user.Uid, user.Gid); err != nil {
			return fmt.Errorf("chown %s failed, %w", dir, err)
		}
	}

	if c.fileMode != os.FileMode(0) {
		if err := os.Chmod(dir, c.fileMode); err != nil {
			return fmt.Errorf("chmod %s failed, %w", dir, err)
		}
	}

	return nil
}
