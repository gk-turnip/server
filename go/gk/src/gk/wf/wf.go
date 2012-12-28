package wf

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

// single wavefront object
type wavefrontObjDef struct {
	verticeList []vertexDef
	normalList  []normalDef
	faceList    []faceDef
	groupName   string // g - start of line
}

// single vertex
// v - start of line
// w defaults to 1
type vertexDef struct {
	x, y, z, w float64
}

// single normal
// vn - start of line
type normalDef struct {
	x, y, z, w float64
}

// single face
// f - start of line
type faceDef struct {
	faceIndexList []faceIndexDef
}

// single face index entry
type faceIndexDef struct {
	vertexIndex    int
	textCoordIndex int
	normalIndex    int
}

// parse the 3D wavefront obj file returning
// a slice of single wavefront objects
func ParseWavefront(fileName string) ([]wavefrontObjDef, error) {
	var results []wavefrontObjDef

	var file *os.File
	var err error

	file, err = os.Open(fileName)
	if err != nil {
		return nil, err
	}

	defer file.Close()

	var reader *bufio.Reader = bufio.NewReader(file)

	var line string
	var wavefrontObj wavefrontObjDef

	for {
		line, err = reader.ReadString('\n')

		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}

		if len(line) > 0 {
			if line[0] != '#' {
				var words []string
				words = strings.Split(strings.TrimRight(line, " \r\n"), " ")

				switch words[0] {
				case "mtllib":
					// skip mtllib lines
				case "usemtl":
					// skip usemtl lines
				case "g":
					if len(words) != 2 {
						err = errors.New(fmt.Sprintf("did not understand g line: %s", line))
						return nil, err
					}
					if !wavefrontObj.isEmpty() {
						results = append(results, wavefrontObj)
					}
					wavefrontObj = *new(wavefrontObjDef)
					wavefrontObj.groupName = words[1]
				case "v":
					if len(words) != 4 {
						err = errors.New(fmt.Sprintf("did not understand v line: %s", line))
						return nil, err
					}
					err = wavefrontObj.handleV(words)
					if err != nil {
						return nil, err
					}
				case "vn":
					if len(words) != 4 {
						err = errors.New(fmt.Sprintf("did not understand vn line: %s", line))
						return nil, err
					}
					err = wavefrontObj.handleVN(words)
					if err != nil {
						return nil, err
					}
				case "f":
					if len(words) < 4 {
						err = errors.New(fmt.Sprintf("did not understand f line: %s", line))
						return nil, err
					}
					err = wavefrontObj.handleF(words)
					if err != nil {
						return nil, err
					}
				default:
					err = errors.New(fmt.Sprintf("did not understand line: %s", line))
					return nil, err
				}
			}
		}
	}

	if !wavefrontObj.isEmpty() {
		results = append(results, wavefrontObj)
	}

	return results, nil
}

func (wavefrontObj *wavefrontObjDef) isEmpty() bool {
	return (len(wavefrontObj.verticeList) == 0) &&
		(len(wavefrontObj.normalList) == 0) &&
		wavefrontObj.groupName == ""
}

func (wavefrontObj *wavefrontObjDef) handleV(words []string) error {
	var vertex vertexDef
	var err error

	vertex.x, err = strconv.ParseFloat(words[1], 64)
	if err != nil {
		return err
	}
	vertex.y, err = strconv.ParseFloat(words[2], 64)
	if err != nil {
		return err
	}
	vertex.z, err = strconv.ParseFloat(words[3], 64)
	if err != nil {
		return err
	}

	wavefrontObj.verticeList = append(wavefrontObj.verticeList, vertex)

	return nil
}

func (wavefrontObj *wavefrontObjDef) handleVN(words []string) error {
	var normal normalDef
	var err error

	normal.x, err = strconv.ParseFloat(words[1], 64)
	if err != nil {
		return err
	}
	normal.y, err = strconv.ParseFloat(words[2], 64)
	if err != nil {
		return err
	}
	normal.z, err = strconv.ParseFloat(words[3], 64)
	if err != nil {
		return err
	}

	wavefrontObj.normalList = append(wavefrontObj.normalList, normal)

	return nil
}

func (wavefrontObj *wavefrontObjDef) handleF(words []string) error {

	var face faceDef
	var err error

	for i := 1; i < len(words); i++ {
		var indexList []string

		indexList = strings.Split(words[i], "/")
		if len(indexList) != 3 {
			err = errors.New(fmt.Sprintf("invalid face index: %s\n", words[i]))
			return err
		}
		if len(indexList[0]) < 1 {
			err = errors.New(fmt.Sprintf("invalid face index: %s\n", words[i]))
			return err
		}
		var faceIndex faceIndexDef
		faceIndex.vertexIndex, err = strconv.Atoi(indexList[0])
		if err != nil {
			return err
		}
		if len(indexList[2]) < 1 {
			err = errors.New(fmt.Sprintf("invalid face index: %s\n", words[i]))
			return err
		}
		faceIndex.normalIndex, err = strconv.Atoi(indexList[2])
		if err != nil {
			return err
		}

		face.faceIndexList = append(face.faceIndexList, faceIndex)
	}

	wavefrontObj.faceList = append(wavefrontObj.faceList, face)
	return nil
}
