package xio

import "io"

// WriteFrom 从指定位置开始将 r 中数据写入到 w 中
func WriteFrom(w io.Writer, r io.ReadSeeker, off int64) error {
	curr, err := r.Seek(0, io.SeekCurrent)
	if err != nil {
		return err
	}
	defer func() { _, _ = r.Seek(curr, io.SeekStart) }()

	if _, err = r.Seek(off, io.SeekStart); err != nil {
		return err
	}

	if _, err = io.Copy(w, r); err != nil {
		return err
	}
	return nil
}
