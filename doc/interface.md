# DocDB

## Documents

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
* DELETE /indexes/:name

# TODO

* PUT
* Sharding (first local to goroutines, later to separate machines)
* Compression
* Indexes on non-strings