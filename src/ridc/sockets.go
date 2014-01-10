/*
Copyright (C) 2013 Rob Britton

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package ridc

import (
  "io"
  "net"
  "log"
  "bufio"
  "regexp"
  "encoding/json"
)

var getByIndex = regexp.MustCompile("^GET /([a-z0-9]+)/([^\n]+)")
var getById = regexp.MustCompile("^GET /([a-z0-9-]+)")
var postIndex = regexp.MustCompile("^POST /indexes/([a-z0-9]+)")
var deleteIndex = regexp.MustCompile("^DELETE /indexes/([a-z0-9]+)")
var post = regexp.MustCompile("^POST ([^\n]+)")
var deleteByIndex = regexp.MustCompile("^DELETE /([a-z0-9]+)/([^\n]+)")
var del = regexp.MustCompile("^DELETE /([a-z0-9-]+)")
var put = regexp.MustCompile("^PUT /([a-z0-9-]+) ([^\n]+)")

func HandleConnection(conn net.Conn, database *Database) {
  reader := bufio.NewReader(conn)
  writer := bufio.NewWriter(conn)

  for {
    input, err := reader.ReadString('\n')

    if err != nil {
      if err == io.EOF {
        // connection was closed
        break
      } else {
        log.Println("Error on socket read:", err)
      }
    }

    // parse input
    if input == "GET /\n" {
      // Just output version
      log.Println("Incoming GET /")
      writer.WriteString("Ridc version " + VERSION + "\n")
    } else if input == "GET /indexes\n" {
      // Get all the indexes
      log.Println("Incoming GET /indexes")
      IOGetIndexes(database, writer)
      writer.WriteString("\n")
    } else if match := postIndex.FindStringSubmatch(input); match != nil {
      // Create an index
      log.Println("Incoming POST /indexes")
      err := database.AddIndex(match[1])
      if err != nil {
        log.Println("Error on creating index: ", err)
        writer.WriteString("{}\n")
      } else {
        writer.WriteString("{}\n")
      }
    } else if match := deleteIndex.FindStringSubmatch(input); match != nil {
      // Delete an index
      log.Println("Incoming delete index")
      err := database.RemoveIndex(match[1])
      if err != nil {
        log.Println("Error on removing index: ", err)
        writer.WriteString("{}\n")
      } else {
        writer.WriteString("{}\n")
      }
    } else if match := getByIndex.FindStringSubmatch(input); match != nil {
      // Get all the documents with field = value
      log.Println("Incoming GET by index")
      IOGetByIndex(database, writer, match[1], match[2])
      writer.WriteString("\n")
    } else if match := getById.FindStringSubmatch(input); match != nil {
      // Get document by ID
      log.Println("Incoming GET by ID")
      doc, exists := database.FindById(match[1])
      if !exists {
        writer.WriteString("{}\n")
      } else {
        output, err := json.Marshal(*doc)

        if err != nil {
          log.Println("Error on marshalling document: ", err)
          writer.WriteString("{}\n")
        } else {
          writer.Write(output)
          writer.WriteString("\n")
        }
      }
    } else if match := post.FindStringSubmatch(input); match != nil {
      log.Println("Incoming POST")
      input := match[1]

      id, err := database.Add(input)

      if err != nil {
        log.Println("Error on saving document: ", err)
        writer.WriteString("{}\n")
      } else {
        writer.WriteString("{\"id\": \"" + id + "\"}\n")
      }
    } else if match := deleteByIndex.FindStringSubmatch(input); match != nil {
      log.Println("Incoming DELETE by index")
      database.DeleteByIndex(match[1], match[2])
      writer.WriteString("{}\n")
    } else if match := del.FindStringSubmatch(input); match != nil {
      log.Println("Incoming DELETE by ID")
      database.DeleteById(match[1])
      writer.WriteString("{}\n")
    } else if match := put.FindStringSubmatch(input); match != nil {
      log.Println("Incoming PUT")
      id := match[1]
      data := match[2]
      database.Update(id, data)
      writer.WriteString("{}\n")
    } else {
      // unrecognized input
      log.Println("Unknown action")
      writer.WriteString("{}\n")
    }
    writer.Flush()
  }
}

func StartSocketInterface(database *Database, port string) {
  listener, err := net.Listen("tcp", ":" + port)

  if err != nil {
    log.Fatalf("Could not start socket listener: ", err)
  }

  for {
    conn, err := listener.Accept()

    if err != nil {
      log.Println("Error on accept: ", err)
    } else {
      go HandleConnection(conn, database)
    }
  }
}
