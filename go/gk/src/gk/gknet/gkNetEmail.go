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

// package to hande network
package gknet

import (
	"bytes"
	"fmt"
	"io"
	"net/smtp"
	"strings"
)

import (
	"gk/gkerr"
)

func SendEmail(address string, from string, toArray []string, subject string, message []byte) (string, *gkerr.GkErrDef) {

	var err error
	var conn *smtp.Client
	var dataWriter io.WriteCloser
	var gkErr *gkerr.GkErrDef
	var sendId string

	conn, err = smtp.Dial(address)
	if err != nil {
		gkErr = gkerr.GenGkErr("smtp.Dial", err, ERROR_ID_SMTP_DIAL)
		return "", gkErr
	}

	conn.Mail(from)
	for _, adr := range toArray {
		conn.Rcpt(adr)
	}

	dataWriter, err = conn.Data()
	if err != nil {
		gkErr = gkerr.GenGkErr("conn.Data", err, ERROR_ID_SMTP_DATA)
		return "", gkErr
	}

	defer dataWriter.Close()

	buf := bytes.NewBufferString("From: " + from + "\nTo: " + toArray[0] + "\nSubject: " + subject + "\n")
	_, err = buf.WriteTo(dataWriter)
	if err != nil {
		gkErr = gkerr.GenGkErr("buf.WriteTo", err, ERROR_ID_SMTP_WRITE)
		return "", gkErr
	}

	buf = bytes.NewBufferString(string(message))
	_, err = buf.WriteTo(dataWriter)
	if err != nil {
		gkErr = gkerr.GenGkErr("buf.WriteTo", err, ERROR_ID_SMTP_WRITE)
		return "", gkErr
	}

	err = conn.Quit()
	if err != nil {
		localError := fmt.Sprintf("%v", err)
		if strings.Index(localError, "Ok") != -1 &&
			strings.Index(localError, "queued as") != -1 {
			sendId = localError
		} else {
			gkErr = gkerr.GenGkErr("conn.Quit", err, ERROR_ID_SMTP_QUIT)
			//gklog.LogTrace(fmt.Sprintf("smtp quit %T [%v] %v",err, err, gkErr))
			return "", gkErr
		}
	}

	return sendId, nil
}
