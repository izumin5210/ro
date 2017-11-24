package query

// BuildForCount implements the type.Query interface
func (q *Query) BuildForCount() (string, []interface{}) {
	return q.countCommand(), q.countArgs()
}

func (q *Query) countCommand() string {
	if q.isWithScore() {
		return zcount
	}
	return zcard
}

func (q *Query) countArgs() []interface{} {
	args := []interface{}{q.key}
	if !q.isWithScore() {
		return args
	}
	min, max := q.getMinAndMax()
	return append(args, min, max)
}
