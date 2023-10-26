# Database migrations

DB migrations are managed using golang migrate https://github.com/golang-migrate/ .
If on a mac then you can install this with `brew install golang-migrate`.

There's little point in duplicating the documentation for golang-migrate here,
but the two commands you'll likely need are

migrate -path /Users/mattp/repos/ncc/db/migrations -database $DATABASE_URL up 
and
migrate -path /Users/mattp/repos/ncc/db/migrations -database $DATABASE_URL down

# Testing

Be sure to migrate the test database if you make any schema changes. Use the same commands as above but specify the test location instead.

Run `DATABASE_URL=<<TEST DATBASE>> go test`

Don't be tempted to run the tests without overriding the db location as
otherwise you will likely lose data. You *are* developing locally and not on
the same as your production system, right?!
