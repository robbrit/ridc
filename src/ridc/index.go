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
)

type IndexBuckets map[string]*list.List

type Index struct {
  Field string
  Documents IndexBuckets
}

// Create a new index
func CreateIndex(field string) (*Index, error) {
  // TODO: do validation on the field

  idx := new(Index)
  idx.Field = field
  idx.Documents = IndexBuckets{}

  return idx, nil
}

// Add an item to this index
func (idx *Index) Add(doc *Document) {
  bvalue, ok := (*doc)[idx.Field]

  if !ok {
    // document doesn't have my index
    return
  }

  value, ok := bvalue.(string)

  if !ok {
    // non-string type
    return
  }

  // see if this value already exists
  docs, exists := idx.Documents[value]

  if !exists || docs == nil {
    docs = list.New()
    idx.Documents[value] = docs
  }

  docs.PushBack(doc)
}

// Remove an item from this index
func (idx *Index) Remove(doc *Document) {
  bvalue, ok := (*doc)[idx.Field]

  if !ok {
    // this document doesn't have my field
    return
  }

  value, ok := bvalue.(string)

  if !ok {
    return
  }

  docs, exists := idx.Documents[value]

  if !exists {
    // this shouldn't happen
    return
  }

  for iter := docs.Front(); iter != nil; iter = iter.Next() {
    if iter.Value == doc {
      docs.Remove(iter)
      break
    }
  }
}

// Remove all documents with a specific value
func (idx *Index) RemoveAll(value string) {
  _, exists := idx.Documents[value]

  if exists {
    delete(idx.Documents, value)
  }
}

// Look up items in this index
func (idx *Index) Find(value string) (*list.List) {
  l, exists := idx.Documents[value]

  if exists {
    return l
  }
  return nil
}
