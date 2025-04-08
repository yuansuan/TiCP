package compress

import (
	"fmt"
	"io"

	"github.com/miolini/datacounter"
)

type debug struct {
	prefix string
	Compressor
}

// NewDebug NewDebug
func NewDebug(prefix string, compressor Compressor) Compressor {
	return &debug{
		prefix,
		compressor,
	}
}

type debugWriter struct {
	prefix                      string
	counterBefore, counterAfter *datacounter.WriterCounter
	underlayer                  io.WriteCloser
}

func (p *debug) newDebugWriter(w io.Writer) (*debugWriter, error) {
	before := datacounter.NewWriterCounter(w)
	compress, err := p.Compressor.Compress(before)
	if err != nil {
		return nil, err
	}
	after := datacounter.NewWriterCounter(compress)

	return &debugWriter{
		p.prefix,
		before, after,
		compress,
	}, nil
}

func (p *debugWriter) Write(buf []byte) (n int, err error) {
	return p.counterAfter.Write(buf)
}

func (p *debugWriter) Close() error {
	err := p.underlayer.Close()
	fmt.Printf("%s CompressWriter: compress=%d uncompress=%d, ratio=%.2f",
		p.prefix, p.counterBefore.Count(), p.counterAfter.Count(),
		float64(p.counterBefore.Count())/float64(p.counterAfter.Count()))

	return err
}

func (p *debug) Compress(writer io.Writer) (io.WriteCloser, error) {
	return p.newDebugWriter(writer)
}

type debugReader struct {
	prefix                      string
	counterBefore, counterAfter *datacounter.ReaderCounter
	underlayer                  io.Reader
}

func (p *debug) newDebugReader(r io.Reader) (*debugReader, error) {
	before := datacounter.NewReaderCounter(r)
	uncompress, err := p.Compressor.Decompress(before)
	if err != nil {
		return nil, err
	}
	after := datacounter.NewReaderCounter(uncompress)

	return &debugReader{
		p.prefix,
		before, after,
		uncompress,
	}, nil
}

func (p *debugReader) Read(buf []byte) (n int, err error) {
	n, err = p.counterAfter.Read(buf)
	if err != nil {
		fmt.Printf("%s CompressReader: compress=%d uncompress=%d, ratio=%.2f",
			p.prefix, p.counterBefore.Count(), p.counterAfter.Count(),
			float64(p.counterBefore.Count())/float64(p.counterAfter.Count()))
	}
	return n, err
}

func (p *debug) Decompress(reader io.Reader) (io.Reader, error) {
	return p.newDebugReader(reader)
}
