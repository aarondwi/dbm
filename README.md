# dbm
ORM-style database-schema migration tools, minus the ORM bloats

[![Build Status](https://travis-ci.org/aarondwi/dbm.svg?branch=master)](https://travis-ci.org/aarondwi/dbm)
[![Go Report Card](https://goreportcard.com/badge/github.com/aarondwi/dbm)](https://goreportcard.com/report/github.com/aarondwi/dbm)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Why?
--------------------------
I am personally not a fan of ORM tools, IMO that increase bloats without any benefits.
But the migration tool model is really nice. Easy to both track and version controlled.
So I made this. Provide usual goodies of ORM database-schema migration, without the ORM

Supported database
--------------------------
1. PostgreSQL

Building
--------------------------
Either build yourself, or download some precompiled version

To build yourself:
1. clone/download this repository
2. run `go build -o dbm .`

After that, put the binary on system/user path, e.g. `export PATH=$PATH:/path/to/dbm/executable`

So it can be used from anywhere in your system

Usage Flow
--------------------------
First, generate the project directory
```
dbm init exampleproject
```

Then, `cd` into the project. Then, you can create `srcfile` by using
```
dbm generate DummyName
```
It will add a file in `src` directory in the format `{UNIXSECONDS}-DummyName.yaml`. The file has `up` and `down` attributes, which you gonna set to whatever the need is.

Up is mostly for applying, such as *CREATE TABLE*, *ADD INDEX*, etc. Down is mostly for removing, such as *DROP TABLE*, *DROP INDEX*, etc. Of course, definition of applying and removing may vary, depending on the case

> Note that, a *srcfile* is only for one thing, for example *CREATE TABLE*. If you need another table, or wanna add > index to existing, you can create another *srcfile*

Before you start applying the schema, you need to create `dbm_logs` table in the database. It can be done using
```
dbm setup
```
This call reads *conf.yaml* file for connection detail, then connect to it, and create the table

Now all is set, you can start applying with
```
dbm up
```
It will read all *srcfile* in src directory, read already-applied file in *dbm_logs*, then apply those not in *dbm_logs*

> Note, you can also specify exact filename, to only apply that filename, with *dbm up --FILENAME SomeFileNameInSrcDirectory*

To check those already applied, you can use
```
dbm status
```
It will list all file in src directory, coupled with its status (*up* if already applied, *down* otherwise)

To un-apply, use
```
dbm down
```
This will **ONLY** drop the latest one applied to the database. We don't want to by mistake dropping all data

> Same as *up*, you can also specify filename to drop specific change, instead of latest one 

See `dbm help` for more detail and usage
