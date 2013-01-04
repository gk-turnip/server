/*
	Copyright 2012-2013 1620469 Ontario Limited.

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

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
