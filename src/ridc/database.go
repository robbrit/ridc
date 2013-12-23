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
  "container/list"
  "encoding/json"
  "github.com/nu7hatch/gouuid"
)

type Database struct {
  Documents map[string]*Document
  indexes map[string]*Index
}

func CreateDatabase() *Database {
  db := new(Database)
  db.Documents = make(map[string]*Document)
  db.indexes = make(map[string]*Index)
  return db
}

func (db *Database) Add(unparsed string) (string, error) {
  // first parse the JSON
  var doc map[string]interface{}

  bytes := []byte(unparsed)
  err := json.Unmarshal(bytes, &doc)

  if err != nil {
    return "", err
  }

  realDoc := make(Document)

  // give me an ID
  uuid, err := uuid.NewV4()

  if err != nil {
    return "", err
  }

  id := uuid.String()

  realDoc["id"] = id

  // copy everything over (is this necessary?)
  for key, value := range doc {
    // should always happen
    realDoc[key] = value
  }

  // add this to myself and each index
  db.Documents[id] = &realDoc

  for _, idx := range db.indexes {
    idx.Add(&realDoc)
  }

  return id, nil
}

// Look up an element by id
func (db *Database) FindById(id string) (*Document, bool) {
  doc, exists := db.Documents[id]
  return doc, exists
}

// Look up elements by index
func (db *Database) FindByIndex(field, value string) (*list.List) {
  index, exists := db.indexes[field]

  if !exists {
    // no such index
    return nil
  }

  return index.Find(value)
}

// Delete element by id
func (db *Database) DeleteById(id string) (int, error) {
  doc, exists := db.Documents[id]

  if !exists {
    return 0, nil
  }

  // need to delete the document
  for _, idx := range db.indexes {
    idx.Remove(doc)
  }

  delete(db.Documents, id)

  return 1, nil
}

// Delete elements by index
func (db *Database) DeleteByIndex(field, value string) (int, error) {
  index, exists := db.indexes[field]

  if !exists {
    // the index does not exist, ignore
    return 0, nil
  }

  // get all the documents that match this value
  docs := index.Find(value)

  if docs == nil || docs.Len() == 0 {
    // nothing to delete
    return 0, nil
  }

  // TODO: this could probably be sped up
  for doc := docs.Front(); doc != nil; doc = doc.Next() {
    // first remove them from every index
    for _, idx := range db.indexes {
      if idx != index {
        idx.Remove(doc.Value.(*Document))
      }
    }

    // now remove it from the database
    pdoc := doc.Value.(*Document)
    id := (*pdoc)["id"].(string)

    delete(db.Documents, id)
  }

  // need to do this one separately or it messes up the linked list
  index.RemoveAll(value)

  return docs.Len(), nil
}

// Get the list of indexes
func (db *Database) Indexes() map[string]*Index {
  return db.indexes
}

// add an index to the database
func (db *Database) AddIndex(field string) error {
  // Create an index on field `field`
  _, exists := db.indexes[field]

  // if it already exists, don't create it
  if exists {
    return nil
  }

  // doesn't exist, so create it
  idx, err := CreateIndex(field)

  if err != nil {
    return err
  }

  db.indexes[field] = idx

  // Add this document to the new index
  for _, doc := range db.Documents {
    idx.Add(doc)
  }

  return nil
}

// remove a certain index
func (db *Database) RemoveIndex(field string) error {
  _, exists := db.indexes[field]

  if !exists {
    // index doesn't exist, ignore
    return nil
  }

  delete(db.indexes, field)

  return nil
}
