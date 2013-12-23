package riddick

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
