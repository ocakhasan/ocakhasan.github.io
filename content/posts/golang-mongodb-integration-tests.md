---
layout: post
title: Write Integration Tests For Mongo With Golang 
summary: let's write some integration tests for mongodb with golang. 
date: 2024-01-19
tags: [golang, mongodb, tests]
---

In most of today's services, there is almost always a data storage to store some data. This could be some relational databases such as `MySQL` or some document database such 
as `MongoDB`. In this blog I will show how to write some tests for MongoDB.

All of the code can be received in [ocakhasan/golang-mongo-integration-tests](https://github.com/ocakhasan/golang-mongo-integration-test)

If you would like to learn more about how to use `MongoDB` with `golang`, please check
1. [Golang & MongoDB Query Cheat Sheet](https://ocakhasan.github.io/golang-mongodb-query-examples/)
2. [Integration Testing for MongoDB-Backed REST APIs with Golang](https://ocakhasan.github.io/golang-mongo-db-rest-api-integration-tests/)

## PREQUISITES

You need to have `Docker` installed in your system since this project requires [testcontainers](https://testcontainers.com/).

We are going to use the [Official SDK](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo) provided by the MongoDB corporation.

We are going to use [testify/suite](https://pkg.go.dev/github.com/stretchr/testify/suite) to setup and teardown tests. This is not a need but it would be feasible to use
`suite.Suite` once there are more methods.

Here are the packages used in this project

```
github.com/stretchr/testify v1.8.4
github.com/testcontainers/testcontainers-go v0.27.0
go.mongodb.org/mongo-driver v1.13.1
```

## Code

I think first we should see the methods we are going to test to get the idea. 

### Model

First let's see the database model we are dealing with

```go
type Book struct {
	ID     primitive.ObjectID `bson:"_id"`
	Author string             `bson:"author"`
	Title  string             `bson:"title"`
	Likes  int                `bson:"likes"`
}
```

### Repository

We have an interface called `Repository` which has all the methods for our project needs.

To make the blog shorter, I added 2 simple methods just to show the idea.


Now let's see the `Repository` interface.

```go
type Repository interface {
	CreateBook(ctx context.Context, book Book) (Book, error)
	FindBook(ctx context.Context, id primitive.ObjectID) (*Book, error)
}
```

We have a struct named `mongoRepository` which implements the `Repository` interface methods.

```go
func NewRepository(db *mongo.Database) *mongoRepository {
    return &mongoRepository{db: db}
}

type mongoRepository struct {
	db *mongo.Database
}

func (m *mongoRepository) CreateBook(ctx context.Context, book Book) (Book, error) {
	if book.ID.IsZero() {
		book.ID = primitive.NewObjectID()
	}

	_, err := m.db.Collection("books").InsertOne(ctx, book)
	if err != nil {
		return Book{}, err
	}

	return book, nil
}

func (m *mongoRepository) FindBook(ctx context.Context, id primitive.ObjectID) (*Book, error) {
	var book Book
	filter := bson.M{
		"_id": id,
	}

	if err := m.db.Collection("books").FindOne(ctx, filter).Decode(&book); err != nil {
		return nil, err
	}

	return &book, nil
}
```

As we can see, `mongoRepository` only accepts a client which is [`mongo.Database`](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo#Database). 

The implementation is super straight forward.

### MONGO DATABASE

To be able to meet the needs of the `mongoRepository` we must create a client.

```go
package mongo

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func NewMongoDatabase(uri string, database string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, err
	}

	db := client.Database(database)

	return db, nil
}
```

### TEST CONTAINERS SETUP

First we need to create a container to be able to test our needs. Creating a container is quite easy with the followings.

Let's create a struct called `TestDatabase` which will have all of our needs to test the functionality.

```go
type TestDatabase struct {
	DbInstance *mongo.Database
	DbAddress  string
	container  testcontainers.Container
}
```

We need the `*mongo.Database` to create the `Repository` method. 

We need the `container` to terminate it when the testing is done.

Now let's create the rest.

```go
func SetupTestDatabase() *TestDatabase {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, dbAddr, err := createMongoContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	_ = tdb.container.Terminate(context.Background())
}
```

The `TearDown` method of the `TestDatabase` will be used after all of the tests run to terminate the container so we free the resources.

Now let's see how to create the container

```go
func createMongoContainer(ctx context.Context) (testcontainers.Container, *mongo.Database, string, error) {
	var env = map[string]string{
		"MONGO_INITDB_ROOT_USERNAME": "root",
		"MONGO_INITDB_ROOT_PASSWORD": "pass",
		"MONGO_INITDB_DATABASE":      "testdb",
	}
	var port = "27017/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{port},
			Env:          env,
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("mongo container ready and running at port: ", p.Port())

	uri := fmt.Sprintf("mongodb://root:pass@localhost:%s", p.Port())
	db, err := NewMongoDatabase(uri, "testdb")
	if err != nil {
		return container, db, uri, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, uri, nil
}
```

It can look complex but in reality.
1. First you create the container
2. Then you get the mapped port from mongo container to your localhost
3. Then create the mongo URI to connect the database.

All of the code can be seen below.

```go
package mongo

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
)

type TestDatabase struct {
	DbInstance *mongo.Database
	DbAddress  string
	container  testcontainers.Container
}

func SetupTestDatabase() *TestDatabase {
	ctx, _ := context.WithTimeout(context.Background(), time.Second*60)
	container, dbInstance, dbAddr, err := createMongoContainer(ctx)
	if err != nil {
		log.Fatal("failed to setup test", err)
	}

	return &TestDatabase{
		container:  container,
		DbInstance: dbInstance,
		DbAddress:  dbAddr,
	}
}

func (tdb *TestDatabase) TearDown() {
	_ = tdb.container.Terminate(context.Background())
}

func createMongoContainer(ctx context.Context) (testcontainers.Container, *mongo.Database, string, error) {
	var env = map[string]string{
		"MONGO_INITDB_ROOT_USERNAME": "root",
		"MONGO_INITDB_ROOT_PASSWORD": "pass",
		"MONGO_INITDB_DATABASE":      "testdb",
	}
	var port = "27017/tcp"

	req := testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "mongo",
			ExposedPorts: []string{port},
			Env:          env,
		},
		Started: true,
	}
	container, err := testcontainers.GenericContainer(ctx, req)
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to start container: %v", err)
	}

	p, err := container.MappedPort(ctx, "27017")
	if err != nil {
		return container, nil, "", fmt.Errorf("failed to get container external port: %v", err)
	}

	log.Println("mongo container ready and running at port: ", p.Port())

	uri := fmt.Sprintf("mongodb://root:pass@localhost:%s", p.Port())
	db, err := NewMongoDatabase(uri, "testdb")
	if err != nil {
		return container, db, uri, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, uri, nil
}
```

### TESTS

Now that we will have

1. Repository methods
2. How to create the container and connect to it with Mongo

We can pass to testing phase. As I mentioned in the beginning, we will use the `testify/suite' method to write the tests.

First, let's create a struct called `RepositorySuite` which has the `suite.Suite`, so it can call the helper functions of the `suite.Suite`. 

```go
type RepositorySuite struct {
	suite.Suite
	repository   Repository
	testDatabase *TestDatabase
}
```

Now we will implement some interfaces from the `testify/suite` package.

First let's implement `SetupAllSuite` interface. 

```go
type SetupAllSuite interface {
	SetupSuite()
}
```

> **SetupAllSuite** has a **SetupSuite** method, which will run **before the tests** in the suite are run.

We are going to implement this interface and we will create the mongo container and the repository.

Now, let's implement the `TearDownAllSuite` interface.

```go
type TearDownAllSuite interface {
	TearDownSuite()
}
```

> **TearDownAllSuite** has a **TearDownSuite** method, which will run **after all the tests** in the suite have been run.

We are going to implement this interface and we will terminate the mongo container.

```go
func (suite *RepositorySuite) SetupSuite() {
	suite.testDatabase = SetupTestDatabase()
	suite.repository = NewRepository(suite.testDatabase.DbInstance)
}

func (suite *RepositorySuite) TearDownSuite() {
	suite.testDatabase.container.Terminate(context.Background())
}
```

Now we can write our tests. Now let's write some tests for the `CreateBook` method. Let's recall the method.

```go
func (m *mongoRepository) CreateBook(ctx context.Context, book Book) (Book, error) {
	if book.ID.IsZero() {
		book.ID = primitive.NewObjectID()
	}

	_, err := m.db.Collection("books").InsertOne(ctx, book)
	if err != nil {
		return Book{}, err
	}

	return book, nil
}
```

The method is quite simple, it checks if the ID is provided by the function. If it is not provided, it generates a new unique id.

So what we can test is
1. Provide a book with no id
2. Provide a book with id

Then check whether they are created or not.

```go
// All methods that begin with "Test" are run as tests within a
// suite.
func (suite *RepositorySuite) TestCreateBook() {
    suite.Run("when id is not provided", func() {
        book := Book{
            Author: "Irvin D. Yalom",
            Title:  "Staring at the Sun: Overcoming the Terror of Death",
            Likes:  100,
        }

        createdBook, createBookErr := suite.repository.CreateBook(context.Background(), book)

        suite.Nil(createBookErr)
        suite.Equal(createdBook.Title, "Staring at the Sun: Overcoming the Terror of Death")
        suite.Equal(createdBook.Author, "Irvin D. Yalom")
        suite.False(createdBook.ID.IsZero())
    })

    suite.Run("when id is provided", func() {
        book := Book{
            ID:     primitive.NewObjectID(),
            Author: "Dostoyevksi",
            Title:  "Notes From the Underground",
            Likes:  100,
        }

        createdBook, createBookErr := suite.repository.CreateBook(context.Background(), book)

        suite.Nil(createBookErr)
        suite.Equal(createdBook, book)
    })
}
```

Now let's write some tests for `FindBook` method of the `mongoRepository`. Let's recall the method.

```go
func (m *mongoRepository) FindBook(ctx context.Context, id primitive.ObjectID) (*Book, error) {
	var book Book
	filter := bson.M{
		"_id": id,
	}

	if err := m.db.Collection("books").FindOne(ctx, filter).Decode(&book); err != nil {
		return nil, err
	}

	return &book, nil
}
```

It is quite simple, it tries to find the book with given id. So to test it,

1. First, we need to create a book, then try to fetch it and it should be successful.
2. Try to fetch a document which not exists, then it should not found it.

```go
func (suite *RepositorySuite) TestFindBook() {
	suite.Run("when there is no record", func() {
		id := primitive.NewObjectID()

		foundBook, findBookErr := suite.repository.FindBook(context.Background(), id)

		suite.Equal(findBookErr, mongo.ErrNoDocuments)
		suite.Nil(foundBook)
	})

	suite.Run("when there is record for given id", func() {
		book := Book{
			Author: "Dostoyevksi",
			Title:  "Notes From the Underground",
			Likes:  100,
		}

		createdBook, createBookErr := suite.repository.CreateBook(context.Background(), book)
		suite.Nil(createBookErr)

		id := createdBook.ID

		foundBook, findBookErr := suite.repository.FindBook(context.Background(), id)

		suite.Nil(findBookErr)
		suite.Equal(*foundBook, createdBook)
	})
}
```

In the second test, first we create a book to ensure that there is some data to fetch.

In the first test, random id is tried to be fetched and it returned an error `mongo.ErrNoDocuments` which states that there is no record for given filter in this collection.


To run the tests, you just need to run

```bash
go test -v ./..
```

You will see output similar to this.

```
2024/01/19 21:24:15 🐳 Creating container for image testcontainers/ryuk:0.6.0
2024/01/19 21:24:15 ✅ Container created: 6096e5f94047
2024/01/19 21:24:15 🐳 Starting container: 6096e5f94047
2024/01/19 21:24:15 ✅ Container started: 6096e5f94047
2024/01/19 21:24:15 🚧 Waiting for container id 6096e5f94047 image: testcontainers/ryuk:0.6.0. Waiting for: &{Port:8080/tcp timeout:<nil> PollInterval:100ms}
2024/01/19 21:24:15 🐳 Creating container for image mongo
2024/01/19 21:24:15 ✅ Container created: 6a6541ee0bcf
2024/01/19 21:24:15 🐳 Starting container: 6a6541ee0bcf
2024/01/19 21:24:15 ✅ Container started: 6a6541ee0bcf
2024/01/19 21:24:15 mongo container ready and running at port:  50354
=== RUN   TestExampleTestSuite/TestCreateBook
=== RUN   TestExampleTestSuite/TestCreateBook/when_id_is_not_provided
=== RUN   TestExampleTestSuite/TestCreateBook/when_id_is_provided
=== RUN   TestExampleTestSuite/TestFindBook
=== RUN   TestExampleTestSuite/TestFindBook/when_there_is_no_record
=== RUN   TestExampleTestSuite/TestFindBook/when_there_is_record_for_given_id
2024/01/19 21:24:20 🐳 Terminating container: 6a6541ee0bcf
2024/01/19 21:24:21 🚫 Container terminated: 6a6541ee0bcf
--- PASS: TestExampleTestSuite (6.04s)
    --- PASS: TestExampleTestSuite/TestCreateBook (5.04s)
        --- PASS: TestExampleTestSuite/TestCreateBook/when_id_is_not_provided (5.04s)
        --- PASS: TestExampleTestSuite/TestCreateBook/when_id_is_provided (0.00s)
    --- PASS: TestExampleTestSuite/TestFindBook (0.00s)
        --- PASS: TestExampleTestSuite/TestFindBook/when_there_is_no_record (0.00s)
        --- PASS: TestExampleTestSuite/TestFindBook/when_there_is_record_for_given_id (0.00s)
PASS

```

As we can see
1. First the container is created
2. The tests are running
3. Container is being terminated.

## REFERENCES

- [mongodb go-sdk](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo)
- [testify/suite](https://pkg.go.dev/github.com/stretchr/testify/suite)
- [testcontainers](https://testcontainers.com/)

