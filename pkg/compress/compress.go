package compress

import (
	"bytes"
	"compress/gzip"
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
	w io.Writer
}

func NewCompressWriter(w io.Writer) *CopmpressWriter {
	return &CopmpressWriter{w: w}
}

func (c *CopmpressWriter) Write(input []byte) (int, error) {
	compressed, err := Compress(input)
	if err != nil {
		return 0, err
	}

	return c.w.Write(compressed)
}
