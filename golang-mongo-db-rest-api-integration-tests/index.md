# Integration Testing for MongoDB-Backed REST APIs with Golang



Building a REST API that plays nice with MongoDB is a common challenge in web development. But how do you make sure it all works seamlessly? That's where integration testing comes in. In this blog post, we're going to break down the process of writing integration tests for your REST API, specifically when MongoDB is in the mix.

You can get all of the code samples for this blog from [this repository](https://github.com/ocakhasan/golang-mongo-rest-api).

## Simple Design of the API

{{< figure src="/images/api-simple-system-design.png" title="simple design of api" >}}

As you can see, only component of our API is MongoDB, which is kind of not realistic for real life examples but you will get the idea 
on how to apply for it for multiple components for integration tests.

## Database Models For the API

{{< mermaid >}}
classDiagram
    Author *-- Book
    Book *-- Comment
    class Author{
        +String id
        +String name
    }
    class Book{
        +String title
        +Author author
        +Int likes
    }
    class Comment{
        +String postTitle
        +Int likes
        +String comment
    }
{{< /mermaid >}}

1. Each author can have many books
2. Each book can have many comments.

Please do not try to validate the design of the models. It is just designed in a way where I can write the code fast and have the tests ready in short period of time.

## API

Our api has 3 different endpoints.

1. `GET /api/books`: returns all of the books with their corresponding comments.
2. `GET /api/author/{id}/books`: returns the books of the author with given id.
3. `POST /api/book`: creates a new book.


You can check the example request and responses from the [project readme](https://github.com/ocakhasan/golang-mongo-rest-api).

### How to Design Integration Tests

Let's check our [PostsController](https://github.com/ocakhasan/golang-mongo-rest-api/blob/main/internal/controllers/controller.go) class which is basically handling all of the requests.

```go
type PostsController struct {
	repo repository.Repository
}

func New(repo repository.Repository) *PostsController {
	return &PostsController{repo: repo}
}
```

As we can see, the only dependency for the `PostsController` is the [Repository](https://github.com/ocakhasan/golang-mongo-rest-api/blob/main/internal/repository/repository.go). Let's check the `Repository` interface.

```go
type Repository interface {
	GetBooksWithComments(ctx context.Context, filter PostFilter) ([]models.BookWithComments, error)
	CreateBook(ctx context.Context, book models.Book) (models.Book, error)
	GetAuthorById(ctx context.Context, id string) (*models.Author, error)
}

func New(db *mongo.Database) Repository {
	return &mongoRepository{db: db}
}

type mongoRepository struct {
	db *mongo.Database
}
```

`mongoRepository` implements the `Repository` interface and, the only dependency for it is the [mongo.Database](https://pkg.go.dev/go.mongodb.org/mongo-driver/mongo@v1.12.1#Database).

In short terms, to be able to test our controller end2end, we need a `MongoDB` connection, but the real question is how to get a real MongoDB connection. 

### Test Containers

The answer is to use the Test-Containers. What is test-containers?

Testcontainers is an open source framework for providing throwaway, lightweight instances of databases, message brokers, web browsers, or just about anything that can run in a Docker container[^1].

So, here is our strategy for testing.

1. Run a MongoDB container with Test-Containers before doing the test.
2. Create the database connection with the MongoDB container.
3. Pass this connection to our API Controllers
4. Do the API Testing
5. Remove the MongoDB container after doing the testing.

### How to Implement With Golang

We can use the [testing.Main](https://pkg.go.dev/testing#M). 

M is a type passed to a TestMain function to run the actual tests [^2].

Let's implement the `TestingMain`

```go
var (
	testDbInstance *mongo.Database
)

func TestMain(m *testing.M) {
	log.Println("setup is running")
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	populateDB()
	exitVal := m.Run()
	log.Println("teardown is running")
	_ = testDB.container.Terminate(context.Background())
	os.Exit(exitVal)
}
```

`populateDB()` function inserts some data to the database so we can do our testing.

Let's check the `SetupTestDatabase()` which is basically creating the MongoDB container and creating the connection to that container.

```go
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
	db, err := database.NewMongoDatabase(uri)
	if err != nil {
		return container, db, uri, fmt.Errorf("failed to establish database connection: %v", err)
	}

	return container, db, uri, nil
}
```

Now that we have the `mongo.Database`, we can create the `Repository` and then we can create the `PostsController`. 

```go
import (
	"github.com/labstack/echo/v4"
	"github.com/ocakhasan/mongoapi/internal/controllers"
	"github.com/ocakhasan/mongoapi/internal/repository"
	"github.com/ocakhasan/mongoapi/pkg/router"
)

func InitializeTestRouter() *echo.Echo {
	postgreRepo := repository.New(testDbInstance)

	userController := controllers.New(postgreRepo)

	return router.Initialize(userController)
}
```

Let's also check the `router.Initialize()` to see which endpoints there are.

```go
func Initialize(controller *controllers.PostsController) *echo.Echo {
	e := echo.New()

	api := e.Group("/api")

	api.GET("/books", controller.GetBooksWithComments())
	api.POST("/book", controller.CreateBook())
	api.GET("/author/:id/books", controller.GetAuthorBooksWithComments())

	return e
}
```

Now we have the router and we can test the endpoints.

### apitest package

You can create the tests with `net/http` package but it will create a lot of boilerplate code. There is a package called [apitest](https://github.com/steinfletcher/apitest). 

It has a lot of easy features such as 

- reading body from a file
- easily check the response status code
- checking body from a file
- and so on...

One of the endpoints is to create books for given author. Let's see the controller code for context on what it is doing.

```go
func (u PostsController) CreateBook() echo.HandlerFunc {
	return func(c echo.Context) error {
		req := new(CreateBookRequest)

		if err := c.Bind(&req); err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"err": err.Error(),
			})
		}

		objId, err := primitive.ObjectIDFromHex(req.AuthorId)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]interface{}{
				"err": err.Error(),
			})
		}

		author, err := u.repo.GetAuthorById(c.Request().Context(), objId.Hex())
		if err != nil {
			if errors.Is(err, mongo.ErrNoDocuments) {
				return c.JSON(http.StatusNotFound, map[string]interface{}{
					"err": "author does not exist",
				})
			}
		}

		createdBook, err := u.repo.CreateBook(c.Request().Context(), models.Book{
			Title:  req.BookName,
			Author: *author,
			Likes:  0,
		})

		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]interface{}{
				"err": err.Error(),
			})
		}

		return c.JSON(http.StatusCreated, map[string]interface{}{
			"book": createdBook,
		})
	}
}
```

- it checks if the author exists
- if author exists, then create the book in the database.

Here is an example request and response from the server.

```bash
curl --location 'http://localhost:3030/api/book' \
--header 'Content-Type: application/json' \
--data '{
    "book_name": "The Idiot",
    "author_id": "654e619760034d917aa0ae64"
}'
```

Response

```
{
    "book": {
        "title": "The Idiot",
        "author": {
            "id": "654e619760034d917aa0ae64",
            "name": "Marcus Aurelius"
        },
        "likes": 0
    }
}
```

As we can see the book is created and returned from the response.

To test this endpoint end2end way you need to pass the correct body, expected response and expected response status code.

I already created the json files for you.

- request body: https://github.com/ocakhasan/golang-mongo-rest-api/blob/main/internal/controllers/integration_test/requests/create_book_success.json
- response body: https://github.com/ocakhasan/golang-mongo-rest-api/blob/main/internal/controllers/integration_test/responses/create_book_response.json

Let's write the test function

```go
package integrationtest

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/ocakhasan/mongoapi/internal/controllers"
	"github.com/ocakhasan/mongoapi/internal/repository"
	"github.com/ocakhasan/mongoapi/pkg/router"
	"github.com/steinfletcher/apitest"
	"github.com/steinfletcher/apitest-jsonpath"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	testDbInstance *mongo.Database
)

func TestMain(m *testing.M) {
	log.Println("setup is running")
	testDB := SetupTestDatabase()
	testDbInstance = testDB.DbInstance
	populateDB()
	exitVal := m.Run()
	log.Println("teardown is running")
	_ = testDB.container.Terminate(context.Background())
	os.Exit(exitVal)
}

func InitializeTestRouter() *echo.Echo {
	postgreRepo := repository.New(testDbInstance)

	userController := controllers.New(postgreRepo)

	return router.Initialize(userController)
}

func TestCreatePostSuccess(t *testing.T) {
	apitest.New().
		Handler(InitializeTestRouter()).
		Post("/api/book").
		Header("content-type", "application/json").
		BodyFromFile("requests/create_book_success.json").
		Expect(t).
		Status(http.StatusCreated).
		BodyFromFile("responses/create_book_response.json").
		End()
}
```

Let's analyze the commands step by step.
1. `apitest.New()`: New creates a new api test. The name is optional and will appear in test reports
2. `Handler(InitializeTestRouter())`: initializes the endpoints and their corresponding handlers.
3. `Post("/api/book").`: sends a `POST` request to `/api/book` endpoint.
4. `Header("content-type", "application/json").`: sets the content-type header.
5. `BodyFromFile("requests/create_book_success.json")`: reads the body from given file and sets the request body.
6. `Status(http.StatusCreated)`: expects the response status code to `http.StatusCreated`.
7. `BodyFromFile("responses/create_book_response.json")`: expects the body to be same as the given file content.

We send a request with given body and we expect the response to be in a certain format and certain data.

As we can see it is super easy to setup and test our endpoints.

Hope you enjoyed the blog. Once again, you may not grasp the whole concept by just looking at the code examples here, please check the [golang-mongo-rest-api](https://github.com/ocakhasan/golang-mongo-rest-api). 

You can check the other tests in the [controller_test.go](https://github.com/ocakhasan/golang-mongo-rest-api/blob/main/internal/controllers/integration_test/controller_test.go) file.


[^1]: https://testcontainers.com/
[^2]: https://pkg.go.dev/testing#M




