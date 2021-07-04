package utils

import (
	"bytes"
	"github.com/sirupsen/logrus"
	"os"
	"path"
)

func FileAppendToFile(f *os.File, path string) error {
	bytesFromFile, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	if !bytes.HasSuffix(bytesFromFile, []byte{'\n'}) {
		bytesFromFile = append(bytesFromFile, '\n')
	}

	_, err = f.Write(bytesFromFile)
	if err != nil {
		return err
	}

	return nil
}

func AllFilesFromDirAppendToFile(f *os.File, pathToDir string, filesList []string) error {
	for i := range filesList {
		err := FileAppendToFile(f, path.Join(pathToDir, filesList[i]))
		if err != nil {
			return err
		}
	}

	return nil
}

func CloseFile(f *os.File) {
	err := f.Close()
	if err != nil {
		logrus.Panic(err)
	}
}
