---
layout: post
title: My Notes on Designing Data Intensive Applications 
summary: Here is my notes on Designing Data Intensive Applications 
date: 2023-08-27 
tags: [backend, database, books]
math: true
---

# MOTIVATION

I work as a Backend engineer almost 2 years now as of August 2023. My work mostly relies on database systems such as MySQL, Redis, Mongo etc. So it would be great to learn the internals or system designs related to those database systems.

Also it is stated in the book that, aynone working on backend side who processes data and the applications they developed uses internet should read this book, so I am quite a fit for the people who should read this book.

I will make a blog post on each of the chapter I read, mostly I will read after my working hours so it will probably take months for me to really finish this book.

##  CHAPTER 1 - Reliable, Scalable and Maintainable Applications

Most applications are data-intensive nowadays, the problems mostly related to amount of data etc.

Most of the tools developed are highly advanced nowadays but none of them can meet all of the needs of different data processing and storing requirements.

### Definitions

#### Reliability

The system should continue to work correctly even though a system error occurs.

- tolerate human errors
- prevents unauthorized access
- there could be some hardware problems such as hard disk crashs, ram becomes faulty etc.
- design systems in a way that human errors opportunity are minimized. 
- test your system, froom unit to integration tests.
- setup monitoring tools, perfomance metrics and error rates.


#### Scalability

The system should handle the load gracefully if the volume (data, network etc) grows.

- what happens to system resources when you increase the load to your system
- how much resource you need to increase when you increase the load.
- response time is what client sees, request sent and response is received from the client
- latency is the duration that a request is waiting to be handled
- in response times it is better to use percentile, not the average. because it does not tell you how many users are affected by a specific number of delay.


#### Maintainability

Project should be easily developed by many other engineers who work on the project.

- cost of software are mostly based on the ongoing mainteiance, not the initial software development.
- projects should be
    - evolvable: meaning making changes should be easy.
    - simple: a project should not be complex, should be easy to work with.


##  CHAPTER 2 - Data Models and Query Languages

### Relational Vs Document Model

Most famous data format is SQL. Goal of relational model was to hide the implementation detail behind a cleaner interface rather than forcing developers to think the internal representation of the data.

The driving forces for NoSQL (Document) Databases
- need for greater scalability
- specialized query operations not supported by SQL
- more dynamic and expressive data models. 

Chapter 2 will be continued.

## CHAPTER 3 - STORAGE AND RETRIEVAL

On the most basic model, a database needs to do 2 operations.

1. it should store the given data
2. when ask it again later, it should give the data back.


The questions needs to be asked as an application developer probably would not be
- how the database handles storage and retrieval internally?

But if you have to tune the program you use, it is better to know the internals of the tool.

### WORLD SIMPLEST DATABASE

Would be a key value store written into a file.

```bash
db_set () {
    echo "$1,$2" >> database
}
db_get () {
    grep "^$1," database | sed -e "s/^$1,//" | tail -n 1
}
```

Similarly to what `db_set` function does, the databasess also uses a **log** internally, append-only data file.


`db_get` function performance is terrible on large scale of data since it traverse the all of the file `O(N)`.

### Index

To retrieve the data efficiently, you need an **index**. Index is an additional data which can be derived from the original set of data. Creating indexes may create an overhead to write operations, since it cannot be more efficient than writing to end of file.


### HASH INDEXES