package rq

// Error represents errors caused by query building.
type Error interface {
	error
	Query() *Query
}

type queryError struct {
	query *Query
	msg   string
}

func newQueryError(q *Query, msg string) error {
	return &queryError{
		query: q,
		msg:   msg,
	}
}

func (e *queryError) Error() string {
	return e.msg
}

func (e *queryError) Query() *Query {
	return e.query
}
