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
  "log"
  "bytes"
  "net/http"
  "encoding/json"
  "github.com/codegangsta/martini"
)

/*
* GET /
* GET /:id
* GET /:index/:value
* POST /
* PUT /:id
* DELETE /:id
* DELETE /:index/:value

## Indexes

* GET /indexes
* POST /indexes
* DELETE /index/:name
*/

func StartWebInterface(database *Database) {
  app := martini.Classic()

  app.Get("/", func() string {
    return "DocDB version " + VERSION
  })

  app.Get("/indexes", func() (int, string) {
    var output bytes.Buffer
    IOGetIndexes(database, &output)
    return 200, string(output.Bytes())
  })

  app.Post("/indexes", func(w http.ResponseWriter, r *http.Request) (int, string) {
    field := r.FormValue("field")
    err := database.AddIndex(field)

    if err != nil {
      log.Println("Error on creating index: ", err)
      return 500, "{}"
    }

    return 200, "{}"
  })

  app.Delete("/index/:name", func(params martini.Params) (int, string) {
    database.RemoveIndex(params["name"])
    return 200, "{}"
  })

  app.Get("/:id", func(params martini.Params) (int, string) {
    doc, exists := database.FindById(params["id"])
    if !exists {
      return 404, "{}"
    }

    output, err := json.Marshal(*doc)

    if err != nil {
      log.Println("Error on marshalling document: ", err)
      return 500, "{}"
    }

    return 200, string(output)
  })

  app.Get("/:index/:value", func(params martini.Params) (int, string) {
    var output bytes.Buffer
    IOGetByIndex(database, &output, params["index"], params["value"])
    return 200, string(output.Bytes())
  })

  app.Post("/", func(w http.ResponseWriter, r *http.Request) (int, string) {
    input := r.FormValue("data")

    id, err := database.Add(input)

    if err != nil {
      log.Println("Error on saving document: ", err)
      return 500, "{}"
    }
    return 200, "{\"id\": \"" + id + "\"}"
  })

  app.Put("/:id", func(w http.ResponseWriter, r *http.Request) (int, string) {
    id := r.FormValue("id")
    data := r.FormValue("data")

    database.Update(id, data)
    return 200, "{}"
  })

  app.Delete("/:id", func(params martini.Params) (int, string) {
    database.DeleteById(params["id"])
    return 200, "{}"
  })

  app.Delete("/:index/:value", func(params martini.Params) (int, string) {
    database.DeleteByIndex(params["index"], params["value"])
    return 200, "{}"
  })

  app.Run()
}

