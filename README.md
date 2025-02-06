## Free Code Camp's Intro To Go Course With Lane Wagner RSS Aggregator Project

These files are simply a following along of [Free Code Camp's Intro To Go Course](https://www.youtube.com/watch?v=un6ZyFkqFKo&pp=ygUaZnJlZSBjb2RlIGNhbXAgaW50cm8gdG8gZ28%3D) with [Lane Wager](https://github.com/wagslane).

You can find the project files [here](https://github.com/bootdotdev/fcc-learn-golang-assets/tree/main/project).

### Setting Up Environment Variables

```bash
go get github.com/joho/godotenv
```

Which just installs the dependency godotenv, which allows for the reading of the
.env file.

Additionally, after installation, he invokes:

```bash
go mod vendor
```

Which creates a `vendor` directory which keeps track of the installed dependencies.

```bash
go mod tidy
```

Cleans up imports so that they are seen by our editor/project. You may need to
run vendor again after running tidy:

```bash
go mod vendor
```

### Setting Up The Router

Install a commonly used go router package called chi:

```bash
go get github.com/go-chi/chi
```

And also the cors package from the same repo maintainer:

```bash
go get github.com/go-chi/cors
```

### Setting Up The Postgres Database

You'll need to install sqlc as well as goose. 

**Goose**

Goose is sets up migrations using .sql files. You'll need to create a database
in your local postgres instance called `rssagg`:

```Postgres
CREATE DATABASE rssagg;
```


Within your env sample file insert your credentials for your local postgres
database (see env.sample).

Using this DB_URL, you can now migrate your database up from within the
`sql/schema` directory run:

```bash
goose postgres $DB_URL up
```

Within psql or pgcli clients, you can now view your created table from within
the rssagg database like so:

```sql
use rssagg;
\d;
\d users;
```

You can also run the following to clean up your database.

```bash
goose postgres $DB_URL down
```

**SQLC**

Now that you have goose and sqlc installed, you can use sqlc to generate your
queries. sqlc requires a configuration yaml file in the root of your project
. Here's what we have in ours:

```yaml
version: "2"
sql:
  - schema: "sql/schema"
    queries: "sql/queries"
    engine: "postgresql"
    gen:
      go:
        out: "internal/database"
```

We now create our Queries within the sql/queries directory, creating .sql files,
here is our users.sql file:

```sql
-- name: CreateUser :one

INSERT INTO users (id, created_at, updated_at, name)
VALUES ($1, $2, $3, $4) RETURNING *;
```

Note the comments at the top, these are required and basically create an sql
function called CreateUser that is called once. Ther est of the sql statement
should be somewhat familiar to you.

You'll need to run the sqlc command from the root of your project (where your
sqlc.yaml file is) like so:

```bash
sqlc generate
```

It then utilizes the sqlc yaml file to generate the sql queries you needed by
reading the passed files (in this case the schema/001_users.sql file and the
queries/user.sql file). It will then generate a series of go files in the
internal/database directory.

We'll now adjust our DB_URL env variable to disable sslmode (and remove the
rssagg db specification) giving us something like this:

```
DB_URL=postgres://postgres:password@localhost:5432/rssagg?sslmode=disable
```

You'll also need to import a driver called `pq` from github, even though it is
never called, this allows the line:

```go
type apiConfig struct {
	DB *database.Queries
}
```

To work. Install it using go get:

```bash
go get github.com/lib/pq
```

And in your main.go file, import it:

```go
import (
    _ "github.com/lib/pq"
)
```

You'll also need to point the database to the sqlc generated files by importing
them into your main.go file:

```go
	"github.com/tomit4/rkabman/internal/database"
```

## Building And Running The Server

```bash
go build && ./rssagg
```

## Quering The DB From pgcli:

```pgcli
use rssagg;
select * from users;
```

## Using Postman to Test The API

From Postman, you can test the API by sending a GET request to the `users`
route. Navigate to the `users` route in Postman and send a GET/POST request.
For POST, you can navigate to `localhost:$PORT/v1/users` and within the Body,
select 'raw' and paste standard JSON such as:

```json
{ "name": "Lane" }
```

### Testing API Key endpoint

We set up a GET request for /users that requires Auth Headers to be sent. After
creating a new user in Postman using the previous section's method, grab the
ApiKey field's value and use it to authenticate against the GET /users route. In
the Authorization tab, under Key, set "Authorization" and in value, set it to
"ApiKey $COPIED_API_KEY". Send it to test.
