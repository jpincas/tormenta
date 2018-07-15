# Tormenta

A small project to help take the pain and boilerplate out of building CRUDy web apps in Go.  Tormenta is currently comprised of the following sub-projects:

- TormentaDB (WIP): a thin functionality layer over BadgerDB key/value store, providing simple object persistence with some querying capabilities.
- TormentaREST (WIP): generic REST API generator - feed in your structs, get a REST API with persistence by TormentaDB. Use triggers to implement logic.
- TormentaGUI (to come): web-based GUI to easily inspect your data in TormentaDB.