package gkcommon

import (
	"io"
	"os"
)

import (
	"gk/gkerr"
)

func GetFileContents(fileName string) ([]byte, *gkerr.GkErrDef) {
	var file *os.File
	var err error
	var gkErr *gkerr.GkErrDef
	var data []byte = make([]byte, 0, 0)

	file, err = os.Open(fileName)
	if err != nil {
		gkErr = gkerr.GenGkErr("os.Open file: "+fileName, err, ERROR_ID_OPEN_FILE)
		return nil, gkErr
	}

	defer file.Close()

	buf := make([]byte, 1024, 1024)
	var readCount int

	for {
		readCount, err = file.Read(buf)
		if readCount > 0 {
			data = append(data, buf[0:readCount]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			gkErr = gkerr.GenGkErr("file.Read file: "+fileName, err, ERROR_ID_READ_FILE)
			return nil, gkErr
		}
	}

	return data, nil
}
