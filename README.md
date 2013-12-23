# Description

Ridc (pronounced Riddick) is an indexed document cache (IDC) - a system that
allows you to store JSON documents in memory and index those documents based on
certain fields. It is like a hybrid of in-memory caching systems like memcached
or Redis, and more sophisticated document storage systems like MongoDB or
CouchDB.

# Compilation

Ensure that Go is installed and is running at least version 1.1.

Install dependencies:

    make depends

Compile:

    make

Run:

    ./ridc

This will start up the server, then you need some sort of client to interact
with it.

# Clients

Here is the list of official clients:

* [Python](https://github.com/robbrit/ridc-python)
