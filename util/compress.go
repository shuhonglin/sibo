package util

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"compress/zlib"
)

type Compressor interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
}

type None struct {}

func (n None) Decompress(data []byte) ([]byte, error) {
	return data, nil
}

func (n None) Compress(data []byte) ([]byte, error) {
	return data, nil
}

type Gzip struct{}

// Unzip unzips data.
func (g Gzip) Decompress(data []byte) ([]byte, error) {
	gr, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer gr.Close()
	data, err = ioutil.ReadAll(gr)
	if err != nil {
		return nil, err
	}
	return data, err
}

// Zip zips data.
func (g Gzip) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type Zip struct{}

// Unzip unzips data.
func (z Zip) Decompress(data []byte) ([]byte, error) {
	zr, err := zlib.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	data, err = ioutil.ReadAll(zr)
	if err != nil {
		return nil, err
	}
	return data, err
}

// Zip zips data.
func (z Zip) Compress(data []byte) ([]byte, error) {
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}
	err = w.Flush()
	if err != nil {
		return nil, err
	}
	err = w.Close()
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
