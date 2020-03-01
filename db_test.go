package main

// A package scope test DB instance isn't always a good idea,
// but it's sufficient in simple cases.
var testDB DB

// One time setup before any tests
func init() {
	testDB = NewDB(":memory:?cache=shared", true)
}

// TODO
// We want "table-driven" unit tests for our DB functions, but current
// functionality only stores data; there's no retrieval. This either limits the
// usefulness of the tests, or requires additional query logic just for the
// tests. Neither of these are good options. In cases like this, where the existing
// functionality is at least sanity checked by integration tests, I think it's pragmatic
// to defer DB level unit tests until query functionality is added.
