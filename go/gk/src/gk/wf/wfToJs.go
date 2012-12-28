package wf

import (
	"flag"
	"fmt"
)

// this function will parse a 3D wavefront obj
// and write out a javascript compatible data structures
func WfToJsStart() {
	var fileName *string

	fileName = flag.String("file", "", "file name")
	flag.Parse()

	if *fileName == "" {
		flag.PrintDefaults()
		return
	}

	var wavefrontObjList []wavefrontObjDef
	var err error

	wavefrontObjList, err = ParseWavefront(*fileName)
	if err != nil {
		fmt.Printf("error: %v on file: %s\n", err, *fileName)
		return
	}

	for _, wfObj := range wavefrontObjList {
		fmt.Printf("var %sVerticeList = [", wfObj.groupName)

		for k, vertice := range wfObj.verticeList {
			if k > 0 {
				fmt.Printf(",")
			}
			fmt.Printf("[%v,%v,%v]", vertice.x, vertice.y, vertice.z)
		}

		fmt.Printf("];\n")

		fmt.Printf("var %sNormalList = [", wfObj.groupName)

		for k, normal := range wfObj.normalList {
			if k > 0 {
				fmt.Printf(",")
			}
			fmt.Printf("[%v,%v,%v]", normal.x, normal.y, normal.z)
		}

		fmt.Printf("];\n")

		fmt.Printf("var %sFaceList = [", wfObj.groupName)

		for k, face := range wfObj.faceList {
			if k > 0 {
				fmt.Printf(",")
			}
			fmt.Printf("[")
			for l, faceIndex := range face.faceIndexList {
				if l > 0 {
					fmt.Printf(",")
				}
				fmt.Printf("[%d,%d]", faceIndex.vertexIndex, faceIndex.normalIndex)
			}
			fmt.Printf("]")
		}

		fmt.Printf("];\n")
	}
}
