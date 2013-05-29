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

// package to capture and handle errors
package gkerr

import (
	"fmt"
	"runtime"
)

const MAX_STACK = 15

type GkErrDef struct {
	message string
	err     error
	errorId uint32
	GkStack []GkStackEntryDef
}

type GkStackEntryDef struct {
	program string
	line    int
}

// capture the current error stack
func GenGkErr(message string, err error, errorId uint32) *GkErrDef {
	var gkErr *GkErrDef = new(GkErrDef)
	var i int

	gkErr.message = message
	gkErr.err = err
	gkErr.errorId = errorId

	for i < MAX_STACK {
		var gkStackEntry GkStackEntryDef
		var ok bool

		_, gkStackEntry.program, gkStackEntry.line, ok = runtime.Caller(i)
		if !ok {
			break
		}

		gkErr.GkStack = append(gkErr.GkStack, gkStackEntry)
		i += 1
	}

	return gkErr
}

func (gkErr *GkErrDef) GetErrorId() uint32 {
	return gkErr.errorId
}

func (gkErr *GkErrDef) String() string {
	var stack string

	stack = gkErr.message
	if gkErr.err != nil {
		stack = stack + fmt.Sprintf(" [%v] %x", gkErr.err, gkErr.errorId)
	}
	stack = stack + "\n"

	for i := 0; i < len(gkErr.GkStack); i++ {
		stack = stack + fmt.Sprintf("\t%s %d\n", gkErr.GkStack[i].program, gkErr.GkStack[i].line)
	}

	return stack
}
