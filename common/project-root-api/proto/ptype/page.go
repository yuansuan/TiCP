package ptype

import (
	"context"
)

// LimitOffset return (limit, offset)
func (p *Page) LimitOffset() (int64, int64) {
	return p.Size, (p.Index - 1) * p.Size
}

// LimitOffsetInt return (limit, offset)
func (p *Page) LimitOffsetInt() (int, int) {
	return int(p.Size), int((p.Index - 1) * p.Size)
}

func NewPageCtx(p *Page, total int64) *PageCtx {
	return &PageCtx{
		Index: p.Index,
		Size:  p.Size,
		Total: total,
	}
}

// FnPageRequest ...
type FnPageRequest func(ctx context.Context, page *Page) (*PageCtx, error)

// ForeachPageRequest ...
func ForeachPageRequest(ctx context.Context, fn FnPageRequest) error {
	var (
		step   = int64(100)
		total  = int64(100)
		offset = int64(0)
		index  = int64(1)
	)

	for offset < total {
		pageCtx, err := fn(ctx, &Page{
			Index: index,
			Size:  step,
		})
		if err != nil {
			return err
		}

		index++
		offset += step
		total = pageCtx.Total
	}
	return nil
}
