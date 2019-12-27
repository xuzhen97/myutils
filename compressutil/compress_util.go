package compressutil

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

//压缩成.gz,类似于linux 下 gzip filename
func CompressGz(srcFile, dest string) error {
	_, filename := filepath.Split(srcFile)
	d, _ := os.Create(dest + filename + ".gz")
	defer d.Close()
	gw := gzip.NewWriter(d)
	defer gw.Close()

	file, err := os.Open(srcFile)
	if err != nil {
		return err
	}
	info, err := file.Stat()

	header, err := tar.FileInfoHeader(info, "")
	if err != nil {
		return err
	}
	gw.Name = header.Name
	if err != nil {
		return err
	}
	_, err = io.Copy(gw, file)
	defer file.Close()
	if err != nil {
		return err
	}
	return nil
}

//解压类似于Linux gzip filename压缩后的.gz文件
func DeCompressGz(gzFile, dest string) error {
	srcFile, err := os.Open(gzFile)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	gr, err := gzip.NewReader(srcFile)
	if err != nil {
		return err
	}
	defer gr.Close()
	filename := dest + gr.Header.Name
	file, err := createFile(filename)
	if err != nil {
		return err
	}
	io.Copy(file, gr)
	return nil
}

func createFile(name string) (*os.File, error) {
	dir, _ := filepath.Split(name)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return nil, err
	}
	return os.Create(name)
}
