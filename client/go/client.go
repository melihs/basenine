// Copyright 2021 UP9. All rights reserved.
// Use of this source code is governed by Apache License 2.0
// license that can be found in the LICENSE file.
//
// This is the client library for Basenine database server.
//
package basenine

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"regexp"
	"strings"
	"sync"
	"time"
)

// Metadata info that's streamed after each record
type Metadata struct {
	Current         uint64 `json:"current"`
	Total           uint64 `json:"total"`
	NumberOfWritten uint64 `json:"numberOfWritten"`
}

// Commands refers to TCP connection modes.
const (
	CMD_INSERT   string = "/insert"
	CMD_QUERY    string = "/query"
	CMD_SINGLE   string = "/single"
	CMD_VALIDATE string = "/validate"
	CMD_MACRO    string = "/macro"
	CMD_LIMIT    string = "/limit"
	CMD_METADATA string = "/metadata"
)

// Closing indicators
const (
	CloseChannel    = "%close%"
	CloseConnection = "%quit%"
)

// Connection is the struct that holds the TCP connection reference.
type Connection struct {
	net.Conn
}

// NewConnection establishes a new connection with the server in the address host:port.
// Returns a Connection reference and an error.
func NewConnection(host string, port string) (connection *Connection, err error) {
	dest := host + ":" + port
	var conn net.Conn
	conn, err = net.Dial("tcp", dest)
	connection = &Connection{conn}
	return
}

// Send sends given []byte to the server which the connection is established to.
func (c *Connection) Send(data []byte) {
	c.SetWriteDeadline(time.Now().Add(1 * time.Second))
	c.Write(data)
}

// SendText is the wrapper around Send method that allows user to send text directly.
func (c *Connection) SendText(text string) {
	c.Send([]byte(fmt.Sprintf("%s\n", text)))
}

// InsertMode turns the connection's mode into INSERT mode
func (c *Connection) InsertMode() {
	c.SendText(CMD_INSERT)
}

// Query is the method that user should use to stream the records from the database.
// It takes the filtering language (query) as the first parameter and
// a []byte channel which the records will be streamed into as the second parameter.
// Third parameter is the channel for streaming metadata, progress of the query.
func (c *Connection) Query(query string, data chan []byte, meta chan []byte) {
	var wg sync.WaitGroup
	go readConnection(&wg, c, data, meta)
	wg.Add(1)

	c.SendText(CMD_QUERY)
	c.SendText(query)
}

// Single returns a single record from the database server specified by the host:port pair
// and by given ID.
func Single(host string, port string, id int) (data []byte, err error) {
	var c *Connection
	c, err = NewConnection(host, port)
	if err != nil {
		return
	}

	ret := make(chan []byte)

	var wg sync.WaitGroup
	go readConnection(&wg, c, ret, nil)
	wg.Add(1)

	c.SendText(CMD_SINGLE)
	c.SendText(fmt.Sprintf("%d", id))

	data = <-ret
	c.Close()
	return
}

// Validate validates the given query against syntax errors by passing the query
// to the database server at host:port
func Validate(host string, port string, query string) (err error) {
	var c *Connection
	c, err = NewConnection(host, port)
	if err != nil {
		return
	}

	ret := make(chan []byte)

	var wg sync.WaitGroup
	go readConnection(&wg, c, ret, nil)
	wg.Add(1)

	c.SendText(CMD_VALIDATE)
	c.SendText(fmt.Sprintf("%s", query))

	data := <-ret
	text := string(data)
	if text != "OK" {
		err = errors.New(text)
	}
	c.Close()
	return
}

// Macro defines a macro that expands into expanded in the database server
// at host:port. Same macro can be overwritten by a second Macro call.
func Macro(host string, port string, macro string, expanded string) (err error) {
	var c *Connection
	c, err = NewConnection(host, port)
	if err != nil {
		return
	}

	ret := make(chan []byte)

	var wg sync.WaitGroup
	go readConnection(&wg, c, ret, nil)
	wg.Add(1)

	c.SendText(CMD_MACRO)
	c.SendText(fmt.Sprintf("%s~%s", macro, expanded))

	data := <-ret
	text := string(data)
	if text != "OK" {
		err = errors.New(text)
	}
	c.Close()
	return
}

// Limit sets a hard-limit for the database size in bytes like 1000000 (1MB)
// Such that actual disk usage ranges between half of the limit (soft-limit) and
// the limit's itself (hard-limit).
func Limit(host string, port string, limit int64) (err error) {
	var c *Connection
	c, err = NewConnection(host, port)
	if err != nil {
		return
	}

	ret := make(chan []byte)

	var wg sync.WaitGroup
	go readConnection(&wg, c, ret, nil)
	wg.Add(1)

	c.SendText(CMD_LIMIT)
	c.SendText(fmt.Sprintf("%d", limit))

	data := <-ret
	text := string(data)
	if text != "OK" {
		err = errors.New(text)
	}
	c.Close()
	return
}

// readConnection is a Goroutine that recieves messages from the TCP connection
// and sends them to a []byte channel provided by the data parameter.
func readConnection(wg *sync.WaitGroup, c *Connection, data chan []byte, meta chan []byte) {
	defer wg.Done()
	for {
		scanner := bufio.NewScanner(c)
		// Prevent buffer overflows
		buf := make([]byte, 0, 64*1024)
		scanner.Buffer(buf, 209715200)

		for {
			ok := scanner.Scan()

			if !ok {
				log.Println("Reached EOF on server connection.")
				return
			}

			bytes := scanner.Bytes()

			command := handleCommands(bytes, meta)
			switch command {
			case CMD_METADATA:
				b := make([]byte, len(bytes[len(CMD_METADATA)+1:]))
				copy(b, bytes[len(CMD_METADATA)+1:])

				meta <- b
				continue
			case CloseConnection:
				log.Println("Server is leaving. Hanging up.")
				break
			}

			b := make([]byte, len(bytes))
			copy(b, bytes)

			data <- b
		}
	}
}

// handleCommands is used by readConnection to make the server's orders
// in the client to take effect. Such that the server can hang up
// the connection.
func handleCommands(bytes []byte, meta chan []byte) string {
	r, err := regexp.Compile("^%.*%$")
	text := string(bytes)
	if err == nil {
		if strings.HasPrefix(text, CMD_METADATA) {
			return CMD_METADATA
		} else if r.MatchString(text) {
			return CloseConnection
		}
	}

	return ""
}
