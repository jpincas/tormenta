# ⚡ Tormenta [![GoDoc](https://godoc.org/github.com/jpincas/tormenta?status.svg)](https://godoc.org/github.com/jpincas/tormenta)

## WIP: Master branch is under active development.  API still in flux. Not ready for serious use yet.

Tormenta is a functionality layer over [BadgerDB](https://github.com/dgraph-io/badger) key/value store.  It provides simple, embedded-object persistence for Go projects with indexing, data querying capabilities and ORM-like features, including loading of relations.  It uses date-based IDs so is particuarly good for data sets that are naturally chronological, like financial transactions, soical media posts etc. Greatly inspired by [Storm](https://github.com/asdine/storm).

## Why would you use this?

Becuase you want to simplify your data persistence and you don't forsee the need for a mult-server setup in the future.  Tormenta relies on an embedded key/value store.  It's fast and simple, but embedded, so you won't be able to go multi-server and talk to a central DB.  If you can live with that, and without the querying power of SQL, Tormenta gives you simplicty - there are no database servers to run, configure and maintain, no schemas, no SQL, no ORMs etc.  You just open a connection to the DB, feed in your Go structs and get normal Go functions with which to persist, retrieve and query your data.  If you've been burned by complex database setups, errors in SQL strings or overly complex ORMs, you might appreciate Tormenta's simplicity.
 
## Features

- JSON for serialisation of data. Uses std lib by default, but you can specify custom serialise/unserialise functions, making it a snip to use [JSONIter](https://github.com/json-iterator/go) or [ffjson](https://github.com/pquerna/ffjson) for speed
- Date-stamped UUIDs mean no need to maintain an ID counter, and
- You get date range querying and 'created at' field baked in
- Simple basic API for saving and retrieving your objects
- Automatic indexing on all fields (can be skipped)
- Option to index by individual words in strings (split index)
- More complex querying of indices including exact matches, text prefix, ranges, reverse, limit, offset and order by
- Combine many index queries with AND/OR logic (but no complex nesting/bracketing of ANDs/ORs)
- Fast counts and sums using Badger's 'key only' iteration
- Business logic using 'triggers' on save and get, including the ability to pass a 'context' through a query
- String / URL parameter -> query builder, for quick construction of queries from URL strings
- Helpers for loading relations

## Quick How To (in place of better docs to come)

- Add import `"github.com/jpincas/tormenta"`
- Add `tormenta.Model` to structs you want to persist
- Add `tormenta:"-"` tag to fields you want to exclude from saving
- Add `tormenta:"noindex"` tag to fields you want to exclude from secondary indexing
- Add `tormenta:"split"` tag to string fields where you'd like to index each word separately instead of the the whole sentence
- Add `tormenta:"nested"` tag to struct fields where you'd like to index each member (using the index syntax "toplevelfield.nextlevelfield")
- Open a DB connection with standard options with `db, err := tormenta.Open("mydatadirectory")` (dont forget to `defer db.Close()`). For auto-deleting test DB, use `tormenta.OpenTest`
- If you want faster serialisation, I suggest [JSONIter](https://github.com/json-iterator/go)
- Save a single entity with `db.Save(&MyEntity)` or multiple (possibly different type) entities in a transaction with `db.Save(&MyEntity1, &MyEntity2)`.
- Get a single entity by ID with `db.Get(&MyEntity, entityID)`.
- Construct a query to find single or mutliple entities with `db.First(&MyEntity)` or `db.Find(&MyEntities)` respectively. 
- Build up the query by chaining methods.
- Add `From()/.To()` to restrict result to a date range (both are optional). 
- Add index-based filters: `Match("indexName", value)`, `Range("indexname", start, end)` and `StartsWith("indexname", "prefix")` for a text prefix search. 
- Chain multiple index filters together.  Default combination is AND - switch to OR with `Or()`.
- Shape results with `.Reverse()`, `.Limit()/.Offset()` and `Order()`.
- Execute the query with `.Run()`, `.Count()` or `.Sum()`.
- Add business logic by specifying `.PreSave()`, `.PostSave()` and `.PostGet()` methods on your structs.
	
See [the example](https://github.com/jpincas/tormenta/blob/tojson/example_test.go) to get a better idea of how to use.

## Gotchas

- Be type-specific when specifying index searches; e.g. `Match("int16field", int(16)")` if you are searching on an `int16` field.  This is due to slight encoding differences between variable/fixed length ints, signed/unsigned ints and floats.  If you let the compiler infer the type and the type you are searching on isn't the default `int` (or `int32`) or `float64`, you'll get odd results.  I understand this is a pain - perhaps we should switch to a fixed indexing scheme in all cases?
- 'Defined' `time.Time` fields e.g. `myTime time.Time` won't serialise properly as the fields on the underlying struct are unexported and you lose the marshal/unmarshal methods specified by `time.Time`.  If you must use defined time fields, specify custom marshalling functions.


## Help Needed / Contributing

- I don't have a lot of low level Go experience, so I reckon the reflect and/or concurrency code could be significantly improved
- I could really do with some help setting up some proper benchmarks
- Load testing or anything similar
- A performant command-line backup utility that could read raw JSON from keys and write to files in a folder structure, without even going through Tormenta (i.e. just hitting the Badger KV store and writing each key to a json file)

## To Do


- [ ] More tests for indexes: more fields, post deletion, interrupted save transactions
- [ ] Nuke/rebuild indices command
- [ ] Documentation / Examples
- [ ] Better protection against unsupported types being passed around as interfaces
- [ ] Fully benchmarked simulation of a real-world use case


## Maybe

- [ ] JSON dump/ backup
- [ ] JSON 'pass through' functionality for where you don't need to do any processing and therefore can skip unmarshalling.
- [ ] Partial JSON return, combined with above, using https://github.com/buger/jsonparser
