# go-postgres-dao

This package represent a base CRUD dao using Go and Postgres. This is not a completed ORM, but if you are of the ones that does no like to use ORM with Go and you are working with Postgres maybe this tool is for you.

With this package you only need to add tags to the struct for generating database tables and work with the CRUD:

Supported tags:

- mandatory : to non null column
- unique: to unique column
- constraint: to create a constraint
- reference: to create foreing keys to other tables
- model: with the name of the column
- type: data type of the column

This is an example of a struct tagged:

type Examplestruct struct {
	ID int64   `model:"id" type:"bigserial" constraint:"note_pk PRIMARY KEY (id)" `
	Description string `model:"description" type:"text" mandatory:"true"`
	Title string `model:"title" type:"varchar(300)" mandatory:"true"`
}

Those tags are useful if you want to generate the database table from the struct. But If don't want to do it in that way. You mush have at least this fields on each model: ID::integer, created_at:timestamp and updated_at:timestamp and the strutc must have the same table name.
