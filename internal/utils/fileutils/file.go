package fileutils

import (
	"io"
	"os"

	"github.com/h2non/filetype"
)

func IsSupportedFileType(filePath string) bool {
	file, _ := os.Open(filePath)
	head := make([]byte, 261)
	file.Read(head)
	fileType, _ := filetype.Match(head)
	return filetype.IsImage(head) || fileType.MIME.Value == "application/pdf"
}

func Copy(src, dest string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func CloseAndDelete(file *os.File) {
	file.Close()
	os.Remove(file.Name())
}
