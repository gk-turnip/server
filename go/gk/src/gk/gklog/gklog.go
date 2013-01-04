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

// package for logging
package gklog

import (
	"fmt"
	"time"
	"os"
)

import (
	"gk/gkerr"
)

const _logSuffix = ".log"
var _logDir string = ""

const _error = 1
const _trace = 2

func LogInit(logDir string) {
	_logDir = logDir
}

func LogError(message string) {
	logAll(_error, message, nil, nil)
}

func LogTrace(message string) {
	logAll(_trace, message, nil, nil)
}

func LogErr(message string, err error) {
	logAll(_error, message, err, nil)
}

func LogGkErr(message string, gkErr *gkerr.GkErrDef) {
	logAll(_error, message, nil, gkErr)
}

func logAll(level int, message string, argErr error, argGkErr *gkerr.GkErrDef) {
	if _logDir == "" {
		fmt.Printf("missing call to gklog.LogInit()\n")
		return
	}

	var levelString string = "Unknown"

	switch level {
	case _trace:
		levelString = "Trace"
	case _error:
		levelString = "Error"
	}

	var fileName string
	var file *os.File
	var err error

	now := time.Now()
	dateName := fmt.Sprintf("%d_%d_%d",now.Year(), now.Month(), now.Day())
	timeStamp := fmt.Sprintf("%d %d %02d:%02d:%02d.%02d",now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second(), now.Nanosecond() / 10000000.0)
	fileName = _logDir + string(os.PathSeparator) + dateName + _logSuffix

	file, err = os.OpenFile(fileName, os.O_WRONLY | os.O_APPEND | os.O_CREATE, 0644)
	if err != nil {
		fmt.Printf("could not open output log file: %s\n",fileName)
		return
	}

	defer file.Close()

	totalMessage := timeStamp + " " + levelString + " " + message
	if err != nil {
		totalMessage = totalMessage + " " + fmt.Sprintf("[%v]",err)
	}
	if argGkErr != nil {
		totalMessage = totalMessage + " <" + argGkErr.String() + ">"
	}

	totalMessage = totalMessage + "\n"

	_, err = file.Write([]byte(totalMessage))
	if err != nil {
		fmt.Printf("could not write log file: %s\n",fileName)
		return
	}
}

