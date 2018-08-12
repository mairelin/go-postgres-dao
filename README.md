# go-postgres-dao

This package  represent a base CRUD dao using Go and Postgres. This is not a completed ORM, but if you are of the ones that does no like to use ORM with Go and you are working with Postgres maybe this tool is for you.

With this package you only need to add tags to the struct for generating database tables and work with the CRUD. Those the tags area useful if you want to generate the database table from the struct. But If don't want to do it in that way. You mush have at lest this fields on each model: ID::integer, created_at:timestamp and updated_at:timestamp.

