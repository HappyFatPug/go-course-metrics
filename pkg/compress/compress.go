package compress

import (
	"bytes"
	"compress/gzip"
	"errors"
	"io"
)

func Compress(input []byte) ([]byte, error) {
	var buf bytes.Buffer

	w := gzip.NewWriter(&buf)

	_, err := w.Write(input)
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func Decompress(input io.Reader) ([]byte, error) {
	r, err := gzip.NewReader(input)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

type CopmpressWriter struct {
	gw     *gzip.Writer
	closer io.Closer
}

func NewCompressWriter(w io.Writer) *CopmpressWriter {
	var closer io.Closer
	if c, ok := w.(io.Closer); ok {
		closer = c
	}
	return &CopmpressWriter{
		gw:     gzip.NewWriter(w),
		closer: closer,
	}
}

func (w *CopmpressWriter) Write(p []byte) (int, error) {
	if w.gw == nil {
		return 0, errors.New("write on closed gzip writer")
	}
	return w.gw.Write(p)
}

func (w *CopmpressWriter) Flush() error {
	if w.gw == nil {
		return errors.New("flush on closed gzip writer")
	}
	return w.gw.Flush()
}

func (w *CopmpressWriter) Close() error {
	if w.gw == nil {
		return nil
	}

	gzErr := w.gw.Close()
	w.gw = nil

	var err error
	if w.closer != nil {
		err = w.closer.Close()
	}

	return errors.Join(gzErr, err)
}
