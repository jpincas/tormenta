# Tormenta

A small project to help take the pain and boilerplate out of building CRUDy web apps in Go.  Tormenta is currently comprised of the following sub-projects:

- [TormentaDB](https://github.com/jpincas/tormenta/tree/master/tormentadb): a functionality layer over BadgerDB key/value store, providing simple object persistence with some querying capabilities.
- [TormentaREST](https://github.com/jpincas/tormenta/tree/master/tormentarest): generic REST API generator - feed in your structs, get a REST API with persistence by TormentaDB. Use triggers to implement logic.
- [TormentaGUI](https://github.com/jpincas/tormenta/tree/master/tormentagui): web-based GUI to easily inspect your data in TormentaDB - feed in your structs, get a database GUI which supports listing items, editing, creating new ones, deleting and advanced queries.