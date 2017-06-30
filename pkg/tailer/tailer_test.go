package tailer

import (
	"testing"
	"time"
)

func Test_buildLogFileName(t *testing.T) {
	aTime := time.Date(2017, 6, 12, 11, 0, 0, 0, time.Local)
	result := buildLogFileName(aTime)
	expect := "error/postgresql.log.2017-06-12-11"
	if result != expect {
		t.Fatalf("result %s != expect %s", result, expect)
	}
}
