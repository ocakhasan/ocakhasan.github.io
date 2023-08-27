# My Notes on Designing Data Intensive Applications


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
