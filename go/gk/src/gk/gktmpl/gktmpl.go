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

package gktmpl

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"
	"os"
)

import (
	"gk/gkerr"
	"gk/gklog"
)

type TemplateDef struct {
	tmpl       *template.Template
	dataBuffer *bytes.Buffer
}

func NewTemplate(templateDir string, fileNames []string) (*TemplateDef, *gkerr.GkErrDef) {
	var gkTemplate *TemplateDef = new(TemplateDef)

	gkTemplate.tmpl = template.New("name")

	var localFileNames []string
	localFileNames = make([]string, len(fileNames), len(fileNames))
	for i := 0; i < len(fileNames); i++ {
		localFileNames[i] = templateDir + string(os.PathSeparator) + fileNames[i] + ".html"
	}

	var err error

	_, err = gkTemplate.tmpl.ParseFiles(localFileNames...)
	if err != nil {
		return nil, gkerr.GenGkErr("tmpl.ParseFiles", err, ERROR_ID_PARSE_FILES)
	}

	return gkTemplate, nil
}

func (gkTemplate *TemplateDef) Build(buildData interface{}) *gkerr.GkErrDef {
	gkTemplate.dataBuffer = bytes.NewBuffer(make([]byte, 0, 0))
	var err error

	gklog.LogError(fmt.Sprintf("trace buildData: %+v", buildData))

	err = gkTemplate.tmpl.ExecuteTemplate(gkTemplate.dataBuffer, "main", buildData)
	if err != nil {
		return gkerr.GenGkErr("tmpl.ExecuteTemplate", err, ERROR_ID_EXECUTE_TEMPLATE)
	}

	return nil
}

func (gkTemplate *TemplateDef) Send(res http.ResponseWriter, req *http.Request) *gkerr.GkErrDef {

	var writeCount int
	var err error

	if gkTemplate.dataBuffer == nil {
		return gkerr.GenGkErr("missing call to Build", nil, ERROR_ID_MISSING_BUILD)
	}

	writeCount, err = res.Write(gkTemplate.dataBuffer.Bytes())
	if err != nil {
		return gkerr.GenGkErr("res.Write", err, ERROR_ID_TEMPLATE_WRITE)
	}
	if writeCount != gkTemplate.dataBuffer.Len() {
		return gkerr.GenGkErr("write count short", nil, ERROR_ID_SHORT_WRITE_COUNT)
	}

	gkTemplate.dataBuffer = nil

	return nil
}
