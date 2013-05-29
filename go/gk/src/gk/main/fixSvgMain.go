package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
)

import (
	"gk/gkerr"
	"gk/gksvg"
)

func main() {

	var inputFileName *string = flag.String("in", "", "svg input filename")
	var outputFileName *string = flag.String("out", "", "svg output filename")

	flag.Parse()

	if (*inputFileName == "") || (*outputFileName == "") {
		flag.PrintDefaults()
		return
	}

	var inputData []byte
	var err error

	inputData, err = readSvgData(*inputFileName)
	if err != nil {
		fmt.Printf("error reading file: %v\n", err)
		return
	}

	var gkErr *gkerr.GkErrDef

	var index int
	index = strings.LastIndex(*inputFileName, "/")
	var prefix string
	prefix = (*inputFileName)[index+1:]
	prefix = prefix[:len(prefix)-4]

	var outputData []byte

	outputData, gkErr = gksvg.FixSvgData(inputData, prefix)
	if gkErr != nil {
		fmt.Printf("error fixing svg file: %s\n", gkErr.String())
		return
	}

	err = writeSvgData(*outputFileName, outputData)
	if err != nil {
		fmt.Printf("error writing file: %v\n", err)
		return
	}
}

func readSvgData(fileName string) ([]byte, error) {
	var result []byte = make([]byte, 0, 0)
	var file *os.File
	var err error

	buf := make([]byte, 128, 128)

	file, err = os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	for {
		var readCount int

		readCount, err = file.Read(buf)
		if readCount > 0 {
			result = append(result, buf[:readCount]...)
		}
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}
	}

	return result, nil
}

func writeSvgData(fileName string, data []byte) error {
	var file *os.File
	var err error
	var writeCount int

	file, err = os.OpenFile(fileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer file.Close()

	writeCount, err = file.Write(data)
	if err != nil {
		return err
	}
	if writeCount != len(data) {
		return errors.New(fmt.Sprintf("short write did: %d should: %d\n", writeCount, len(data)))
	}

	return nil
}
