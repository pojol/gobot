package utils

import (
	"bytes"
	"errors"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

type UploadFile struct {
	File     multipart.File
	Header   *multipart.FileHeader
	fileExt  string //file extensions
	fileName string
}

func NewUploadFile(file multipart.File, header *multipart.FileHeader) *UploadFile {
	return &UploadFile{
		File:     file,
		Header:   header,
		fileName: header.Filename,
		fileExt:  filepath.Ext(header.Filename), //update for issue #99
	}
}

// FileName get upload file client-local name
func (f *UploadFile) FileName() string {
	return f.fileName
}

// Size get upload file size
func (f *UploadFile) Size() int64 {
	return f.Header.Size
}

// SaveFile save file in server-local with filename
// special:
// if you SaveFile, it's will cause empty data when use ReadBytes
func (f *UploadFile) SaveFile(fileName string) (size int64, err error) {
	size = 0
	if fileName == "" {
		return size, errors.New("filename not allow empty")
	}

	fileWriter, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return size, err
	}
	defer fileWriter.Close()
	size, err = io.Copy(fileWriter, f.File)
	return size, err
}

// GetFileExt get upload file extensions
func (f *UploadFile) GetFileExt() string {
	return f.fileExt
}

// ReadBytes Bytes returns a slice of byte hoding the UploadFile.File
// special:
// if you read bytes, it's will cause empty data in UploadFile.File, so you use SaveFile will no any data to save
func (f *UploadFile) ReadBytes() []byte {
	buf := new(bytes.Buffer)
	buf.ReadFrom(f.File)
	return buf.Bytes()
}
