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

package gksvg

import (
	"fmt"
	"io"
	"bytes"
	"strings"
	"encoding/xml"
)

import (
	"gk/gkerr"
)

type nodeDef struct {
	nameSpace string
	nameLocal string
	charData []byte
	childList []*nodeDef
	attributeList []*attributeDef
	parentNode *nodeDef
}

type attributeDef struct {
	nameSpace string
	nameLocal string
	value string
}

// remove all whitespaces between tags
// remove <?xml version="1.0" ?>
// remove all comments <!-- -->
// add single <g> around everything within svg
// rename any "id" to "prefix_id"
func FixSvgData(svgData []byte, prefix string) ([]byte, *gkerr.GkErrDef) {
	var rootNode *nodeDef
	var gkErr *gkerr.GkErrDef

	rootNode, gkErr = parseSvg(svgData)
	if gkErr != nil {
		return nil, gkErr
	}

	fixNamespace(rootNode)

	addGNode(rootNode)
	fixId(rootNode, prefix)

	var buf *bytes.Buffer = bytes.NewBuffer(make([]byte, 0, 512))

	rebuildSvg(rootNode, buf)

	return buf.Bytes(), nil
}

func parseSvg(svgData []byte) (*nodeDef, *gkerr.GkErrDef) {
	var currentNode *nodeDef = new(nodeDef)
	var gkErr *gkerr.GkErrDef

	var reader *bytes.Buffer = bytes.NewBuffer(svgData)

	var decoder *xml.Decoder

	decoder = xml.NewDecoder(reader)

	var currentCharData = make([]byte, 0, 0)
	for {
		var token xml.Token
		var err error

		token, err = decoder.Token()
		if err != nil {
			if err == io.EOF {
				break
			}
			gkErr = gkerr.GenGkErr("decoder.Token", err, ERROR_ID_SVG_PARSE)
			return nil, gkErr
		}

		switch token := token.(type) {
		case xml.StartElement:
			var startElement = xml.StartElement(token)

			if currentNode == nil {
				currentNode.childList = make([]*nodeDef,0,0)
			}

			var childNode *nodeDef = new(nodeDef)
			childNode.nameSpace = startElement.Name.Space
			childNode.nameLocal = startElement.Name.Local
			childNode.attributeList = make([]*attributeDef,0,0)
			childNode.charData = make([]byte,0,0)

			for _, attr := range startElement.Attr {
				var attribute *attributeDef = new(attributeDef)
				attribute.nameSpace = attr.Name.Space
				attribute.nameLocal = attr.Name.Local
				attribute.value = attr.Value
				childNode.attributeList = append(childNode.attributeList,attribute)
			}

			currentNode.childList = append(currentNode.childList,childNode)
			childNode.parentNode = currentNode
			currentNode = childNode
		case xml.EndElement:
			currentNode.charData = currentCharData
			currentNode = currentNode.parentNode
			currentCharData = make([]byte,0,0)
		case xml.CharData:
			var charData = xml.CharData(token)
			currentCharData = append(currentCharData, charData...)
		}
	}

	if (len(currentNode.childList) == 1) && (currentNode.childList[0].nameLocal == "svg") {
		return currentNode.childList[0], nil
	}
	return nil, gkerr.GenGkErr("invalid svg, no starting svg tag or multipe root nodes", nil, ERROR_ID_INVALID_SVG)
}

func addGNode(node *nodeDef) {
	var gNode *nodeDef = new(nodeDef)

	gNode.nameSpace = ""
	gNode.nameLocal = "g"
	gNode.charData = make([]byte, 0, 0)
	gNode.attributeList = make([]*attributeDef, 0, 0)
	gNode.childList = node.childList
	node.childList = []*nodeDef {gNode}
}

func fixNamespace(node *nodeDef) {
	var namespaceMap map[string]string

	namespaceMap = getNamespaceMap(node)

	fixNamespaceRecursive(node, namespaceMap)
}

func fixNamespaceRecursive(node *nodeDef, namespaceMap map[string]string) {
	if node.nameSpace != "" {
		var shortName string
		shortName = namespaceMap[node.nameSpace]
		node.nameSpace = shortName
	}

	if node.nameSpace == "svg" {
		node.nameSpace = ""
	}

	for _, attribute := range node.attributeList {
		if attribute.nameSpace != "" {
			var shortName string
			shortName = namespaceMap[attribute.nameSpace]
			if shortName != "" {
				attribute.nameSpace = shortName
			}
		}
	}

	for _, childNode := range node.childList {
		fixNamespaceRecursive(childNode, namespaceMap)
	}
}

func getNamespaceMap(node *nodeDef) map[string]string {
	namespaceMap := make(map[string]string)

	for _, attribute := range node.attributeList {
		if attribute.nameSpace == "xmlns" {
			namespaceMap[attribute.value] = attribute.nameLocal
		}
	}

	return namespaceMap
}

func fixId(node *nodeDef, prefix string) *gkerr.GkErrDef {
	var idMap map[string]string = make(map[string]string)
	var gkErr *gkerr.GkErrDef

	populateIdMap(node, prefix, idMap)

	gkErr = substituteIdMap(node, idMap)
	if gkErr != nil {
		return gkErr
	}

	return nil
}

func populateIdMap(node *nodeDef, prefix string, idMap map[string]string) {

	for _, attribute := range node.attributeList {
		if attribute.nameLocal == "id" {
			var newId string = prefix + "_" + attribute.value
			idMap[attribute.value] = newId
			attribute.value = newId
		}
	}

	for _, childNode := range node.childList {
		populateIdMap(childNode, prefix, idMap)
	}
}

func substituteIdMap(node *nodeDef, idMap map[string]string) *gkerr.GkErrDef {
	var gkErr *gkerr.GkErrDef

	for _, attribute := range node.attributeList {
		attribute.value, gkErr = substituteOneAttributeId(
			idMap, attribute.nameSpace, attribute.nameLocal, attribute.value)
		if gkErr != nil {
			return gkErr
		}
	}

	for _, childNode := range node.childList {
		gkErr = substituteIdMap(childNode, idMap)
		if gkErr != nil {
			return gkErr
		}
	}
	return nil
}

// this assumes that the space is called xlink for
// xmlns:xlink="http://www.w3.org/1999/xlink"
func substituteOneAttributeId(idMap map[string]string, space string, name string, value string) (string, *gkerr.GkErrDef) {
	var returnValue string = value
	var ok bool
	var gkErr *gkerr.GkErrDef

	switch name {
	case "style":
		if space == "" {
			var id string

			id, gkErr = getIdOutOfStyle(value)
			if gkErr != nil {
				return "", gkErr
			}
			if id != "" {
				var newId string

				newId, ok = idMap[id]
				if !ok {
					gkErr = gkerr.GenGkErr("could not find id in idMap " + id, nil, ERROR_ID_INTERNAL_ID_MAP)
					return "", gkErr
				}
				returnValue = strings.Replace(value, "fill:url(#" + id + ")", "fill:url(#" + newId + ")", 1)
			}
		}
	case "href":
		if space == "xlink" {
			var id string
			var newId string

			id, gkErr = getIdOutOfHref(value)
			if gkErr != nil {
				return "", gkErr
			}
			newId, ok = idMap[id]
			if !ok {
				gkErr = gkerr.GenGkErr("could not find id in idMap " + id, nil, ERROR_ID_INTERNAL_ID_MAP)
				return "", gkErr
			}
			returnValue = strings.Replace(value, "#" + id, "#" + newId, 1)
		}
	}

	return returnValue, nil
}

func getIdOutOfHref(input string) (string, *gkerr.GkErrDef) {
	if len(input) > 2 {
		if input[0] == '#' {
			var id string

			id = input[1:]
			return id, nil
		}
	}
	return "", gkerr.GenGkErr("href id parse: " + input, nil, ERROR_ID_HREF_PARSE)
}

func getIdOutOfStyle(input string) (string, *gkerr.GkErrDef) {
	var index1, index2 int

	index1 = strings.Index(input, "fill:url(#")
	if index1 == -1 {
		return "", nil
	}

	index2 = strings.Index(input[index1 + 10:],")")
	if index2 == -1 {
		return "", gkerr.GenGkErr("style id parse: " + input, nil, ERROR_ID_STYLE_PARSE)
	}

	return input[index1 + 10:index1 + index2 + 10], nil
}


func rebuildSvg(node *nodeDef, buf io.Writer) *gkerr.GkErrDef {

	var gkErr *gkerr.GkErrDef
	var result []byte = make([]byte, 0, 128)
	var err error

	result = append(result,'<')
	if node.nameSpace != "" {
		result = append(result, []byte(node.nameSpace)...)
		result = append(result, ':')
	}
	result = append(result, []byte(node.nameLocal)...)

	for _, attribute := range node.attributeList {
		result = append(result, ' ')
		if attribute.nameSpace != "" {
			result = append(result, []byte(attribute.nameSpace)...)
			result = append(result, ':')
		}
		result = append(result, []byte(fmt.Sprintf("%s=\"%s\"", attribute.nameLocal, escapeXML([]byte(attribute.value))))...)
	}
	result = append(result, '>')

	_, err = buf.Write(result)
	if err != nil {
		gkErr = gkerr.GenGkErr("write fix svg results", err, ERROR_ID_WRITE_SVG)
		return gkErr
	}
	for _, childNode := range node.childList {
		gkErr = rebuildSvg(childNode, buf)

	}
	result = make([]byte, 0, 16)
	result = append(result, []byte("</")...)

	if node.nameSpace != "" {
		result = append(result, []byte(node.nameSpace)...)
		result = append(result, ':')
	}
	result = append(result, []byte(node.nameLocal)...)
	result = append(result, '>')

	_, err = buf.Write(result)
	if err != nil {
		gkErr = gkerr.GenGkErr("write fix svg results", err, ERROR_ID_WRITE_SVG)
		return gkErr
	}

	return nil
}

/*
func DumpSvg(node *nodeDef, tab int) {
	for i := 0;i < tab;i++ {
		fmt.Print(" ")
	}

	fmt.Printf("<")
	if node.nameSpace != "" {
		fmt.Printf("%s:", node.nameSpace)
	}
	fmt.Printf("%s",node.nameLocal)
	for _, attribute := range node.attributeList {
		fmt.Printf(" ",)
		if attribute.nameSpace != "" {
			fmt.Printf("%s:", attribute.nameSpace)
		}
		fmt.Printf("%s=\"%s\"", attribute.nameLocal, escapeXML([]byte(attribute.value)))
	}
	fmt.Printf(">\n")
	for _, childNode := range node.childList {
		DumpSvg(childNode, tab + 1)
	}
	for i := 0;i < tab;i++ {
		fmt.Print(" ")
	}
	fmt.Printf("</")
	if node.nameSpace != "" {
		fmt.Printf("%s:", node.nameSpace)
	}
	fmt.Printf("%s>\n",node.nameLocal)
}
*/

func escapeXML(input []byte) []byte {
	var buf *bytes.Buffer = bytes.NewBuffer(make([]byte,0,len(input) + 10))

	xml.Escape(buf, input)

	return buf.Bytes()
}

