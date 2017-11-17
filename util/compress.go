package util

import (
	"bytes"
	"compress/gzip"
	"io/ioutil"
	"compress/zlib"
)

// Unzip unzips data.
func UnGzip(data []byte) ([]byte, error) {
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
func Gzip(data []byte) ([]byte, error) {
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

// Unzip unzips data.
func UnzipByZlib(data []byte) ([]byte, error) {
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
func ZipByZlib(data []byte) ([]byte, error) {
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