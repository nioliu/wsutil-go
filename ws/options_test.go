package ws

import (
	"testing"
)

func TestWithTags(t *testing.T) {
	conn := &SingleConn{}
	tags := WithTags("a", "b")
	conn.options = []Option{tags}
	apply(conn)
	println(conn.tags)
}
