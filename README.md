# Tormenta

A small project to help take the pain and boilerplate out of building CRUDy web apps in Go.  Currently a WIP - but new features and stability enhancements are constantly being added as I explore use cases. Tormenta is comprised of the following sub-projects:

- [TormentaDB](https://github.com/jpincas/tormenta/tree/master/tormentadb): a functionality layer over BadgerDB key/value store, providing simple object persistence with some querying capabilities.
- [TormentaREST](https://github.com/jpincas/tormenta/tree/master/tormentarest): generic REST API generator - feed in your structs, get a REST API with persistence by TormentaDB. Use triggers to implement logic.
- [TormentaGUI](https://github.com/jpincas/tormenta/tree/master/tormentagui): web-based GUI to easily inspect your data in TormentaDB - feed in your structs, get a database GUI which supports listing items, editing, creating new ones, deleting and advanced queries.

## Why would you use this?

Becuase you want to simplify your data persistence and you don't forsee the need for a mult-server setup in the future.  Tormenta relies on the excellent, embedded key/value store 'Badger'.  It's fast and simple, but embedded, so you won't be able to go multi-server and talk to a central DB.  If you can live with that, and without the querying power of SQL, Tormenta gives you simplicty - there are no database servers to run, configure and maintain, no schemas, no SQL, no ORMs etc.  You just open a connection to the DB, feed in your Go structs and get normal Go functions with which to persist, retrieve and query your data.  If you've been burned by complex database setups, errors in SQL strings or overly complex ORMs, you might appreciate Tormenta's simplicity.