# ⚡ Tormenta [![GoDoc](https://godoc.org/github.com/jpincas/tormenta?status.svg)](https://godoc.org/github.com/jpincas/tormenta)

Tormenta is a functionality layer over BadgerDB key/value store.  It provides simple, embedded object persistence for Go projects with some data querying capabilities and ORM-like features.  It uses date-based IDs so is particuarly good for data sets that are naturally chronological, like financial transactions, soical media posts etc. Greatly inspired by [Storm](https://github.com/asdine/storm) and powered by:

- [BadgerDB](https://github.com/dgraph-io/badger)
- ['V6' UUIDs](https://github.com/bradleypeabody/gouuidv6)
- [JSONIter](https://github.com/json-iterator/go)

### (Update Jan 2019)

- Refactored the repo and removed the REST API layer and GUI for now. The REST layer will make a reappearance some time in the future, along with a better GUI built on top of it (in Elm!).  
- Switch from MessagePack in favour of good old JSON.  Although this will worsen performance a bit, it means no more codegen, making Tormenta easier to get started with, which fits with the philosophy of the project. 
- Still a WIP, but being used in two projects, both in dev and going to production in the next couple of months, so the API will have to become at least somewhat stable soon! Any questions, hit me up.

## Why would you use this?

Becuase you want to simplify your data persistence and you don't forsee the need for a mult-server setup in the future.  Tormenta relies on the excellent, embedded key/value store 'Badger'.  It's fast and simple, but embedded, so you won't be able to go multi-server and talk to a central DB.  If you can live with that, and without the querying power of SQL, Tormenta gives you simplicty - there are no database servers to run, configure and maintain, no schemas, no SQL, no ORMs etc.  You just open a connection to the DB, feed in your Go structs and get normal Go functions with which to persist, retrieve and query your data.  If you've been burned by complex database setups, errors in SQL strings or overly complex ORMs, you might appreciate Tormenta's simplicity.
 
## Features

- JSON for serialisation of data using `jsoniter` for speed
- Date-stamped UUIDs mean no need to maintain an ID counter, and
- You get date range querying and 'created at' field baked in
- Simple basic API for saving and retrieving your objects
- Option to skip saving of any field, independently of `json` struct tags
- Automatic indexing on all fields (can be skipped)
- Option to index by individual words in strings (split index)
- More complex querying of indices including exact matches, text prefix, ranges, reverse, limit, offset and order by
- Combine queries with AND/OR to arbitrary depth
- Fast counts and sums using Badger's 'key only' iteration
- Business logic using 'triggers' on save and get, including the ability to pass a 'context' through a query
- URL parameter -> query builder in package `urltoquery`, for quick construction of queries from URL strings
- Helpers for loading relations (WIP)

## Quick How To

- Add import `"github.com/jpincas/tormenta"`
- Add `tormenta.Model` to structs you want to persist
- Add `tormenta:"-"` tag to fields you don't want to save
- Add `tormenta:"noindex"` tag to fields you want to exclude from secondary indexing
- Add `tormenta:"split"` tag to string fields where you'd like to index each word separately instead of the the whole sentence
- Open a DB connection with `db, err := tormenta.Open("mydatadirectory")` (dont forget to `defer db.Close()`). For auto-deleting test DB, use `tormenta.OpenTest`
- Save a single entity with `db.Save(&MyEntity)` or multiple entities with `db.Save(&MyEntity1, &MyEntity2)`.
- Get a single entity by ID with `db.Get(&MyEntity, entityID)`.
- Construct a query to find single or mutliple entities with `db.First(&MyEntity)` or `db.Find(&MyEntities)` respectively. 
- Build up the query by chaining methods: `From()/.To()` to add a date range, `Match("indexName", value)` to add an exact match index search, `Range("indexname", start, end)` to add a range search, `StartsWith("indexname", "prefix")` for a text prefix search, `.Reverse()` to reverse the fullStruct of searching/results and `.Limit()/.Offset()` to limit the number of results. `Order()` can be used to specify results order, but with caveats (see below).
- Kick off the query with `.Run()`, or `.Count()` if you just need the count.  `.QuickSum()` is also available for float/int index searches.
- Add business logic by specifying `.PreSave()`, `.PostSave()` and `.PostGet()` methods on your structs.
	
See [the example](https://github.com/jpincas/tormenta/blob/tojson/example_test.go) to get a better idea of how to use.

## Gotchas

- Be type-specific when specifying index searches; e.g. `Match("int16field", int(16)")` if you are searching on an `int16` field.  This is due to slight encoding differences between variable/fixed length ints, signed/unsigned ints and floats.  If you let the compiler infer the type and the type you are searching on isn't the default `int` (or `int32`) or `float64`, you'll get odd results.  I understand this is a pain - perhaps we should switch to a fixed indexing scheme in all cases?
- Due to the way Tormenta returns results by iterating ordered keys, ordering functionality is limited to non-index query searches.  Essentially, adding `Order("myField")` creates an index search on `myField` but without any range (i.e. returns all results).  If you are already using an index, e.g. with `Range("someOtherField", 1, 2)` then that index will take priority and results would be ordered by `someOtherField`.  If you are only filtering by date with `From()/.To()` you CAN independently order as that doesn't use indexes.  If you are doing complex AND/OR query combinations which rely on multiple indexes, then date/ID ordering is the only option - sorry!  In general, best practice would be to limit your results set as much as possible and order results in application code if you require a different order to what you get from Tormenta.

## To Do

- [x] Correct indexing and aggregation on defined fields
- [x] Byte-ordered floats (UREGENT)
- [x] Index deletion on reecord deletion
- [x] Index update on record edit
- [x] Date field indexing (URGENT)
- [ ] More tests for indexes: more fields, post deletion, interrupted save transactions
- [ ] Nuke/rebuild indices command
- [ ] Documentation / Examples
- [x] Delete
- [x] Logic triggers (preSave, postSave, postGet)
- [x] Relation loading helpers: load single relation by ID
- [x] Relation loading helpers: load multiple relations by single ID
- [x] Relation loading helpers: load nested relations
- [x] Relation loading helpers: load multiple relations by slice of IDs
- [ ] Relation loading helpers: load relations of embedded structs
- [ ] Relation loading helpers: load relations by query (e.g. all unpaid invoices) - using reference ID stored on relation
- [ ] Document all the relation loading stuff
- [ ] Better error reporting from query construction
- [ ] Better protection against unsupported types being passed around as interfaces
- [ ] Fully benchmarked simulation of a real-world use case
- [x] Slices as indexes -> multiple index entries
- [x] Stack multiple queries, execute as AND/OR, execute in parallel
- [x] Split-string indexing with 'split' tag
- [x] 'Starts with' index match
- [x] Indexes on by default
- [x] Multiple entity `Get()`
- [x] Bulk unmarshall rather than 1 at a time? Concurrent?

## Maybe

- [ ] JSON dump/ backup
- [ ] JSON 'pass through' functionality for where you don't need to do any processing and therefore can skip unmarshalling.
- [ ] Partial JSON return, combined with above, using https://github.com/buger/jsonparser