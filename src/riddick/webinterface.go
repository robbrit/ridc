package riddick

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

