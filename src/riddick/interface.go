package riddick

import (
  "encoding/json"
)

var VERSION = "0.1"

type CanWriteString interface {
  WriteString(str string) (int, error)
  Write(b []byte) (int, error)
}

func IOGetIndexes(db *Database, w CanWriteString) {
  indexes := db.Indexes()

  first := true

  w.WriteString("[")

  for name, _ := range indexes {
    if first {
      first = false
    } else {
      w.WriteString(",")
    }

    w.WriteString("\"" + name + "\"")
  }

  w.WriteString("]")
}

func IOGetByIndex(db *Database, w CanWriteString, field, value string) {
  docs := db.FindByIndex(field, value)

  w.WriteString("[")

  first := true

  if docs != nil {
    for doc := docs.Front(); doc != nil; doc = doc.Next() {
      byteOutput, err := json.Marshal(doc.Value)

      if first {
        first = false
      } else {
        w.WriteString(",")
      }

      if err == nil {
        w.Write(byteOutput)
      }
    }
  }

  w.WriteString("]")
}
