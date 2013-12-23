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
