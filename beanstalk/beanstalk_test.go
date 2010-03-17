package beanstalk

import (
	"bytes"
	//"fmt"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

type ReadWriter struct {
	io.Reader
	io.Writer
}

func responder(reply string) (io.ReadWriter, *bytes.Buffer) {
	wr := new(bytes.Buffer)
	rd := strings.NewReader(reply)
	return &ReadWriter{rd, wr}, wr
}

func TestPutReplyEOF(t *testing.T) {
	rw, _ := responder("INSERTED 1") // no traling LF, so we hit EOF
	c := newConn("<fake>", rw)
	id, err := c.Tube("default").Put("a", 0, 0, 0)

	if id != 0 {
		t.Error("expected id 0, got", id)
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "put 0 0 0 1\r\na\r\n" {
		t.Errorf("expected put command, got %q", berr.Cmd)
	}

	if berr.Reply != "INSERTED 1" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != os.EOF {
		t.Errorf("expected os.EOF, got %v", berr.Error)
	}
}

func TestPutReplyUnknown(t *testing.T) {
	rw, _ := responder("FOO 1\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("default").Put("a", 0, 0, 0)

	if id != 0 {
		t.Error("expected id 0, got", id)
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "put 0 0 0 1\r\na\r\n" {
		t.Errorf("expected put command, got %q", berr.Cmd)
	}

	if berr.Reply != "FOO 1\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != BadReply {
		t.Errorf("expected beanstalk.BadReply, got %v", berr.Error)
	}
}

func TestPutReplyTooManyArgs(t *testing.T) {
	rw, _ := responder("INSERTED 1 2\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("default").Put("a", 0, 0, 0)

	if id != 0 {
		t.Error("expected id 0, got", id)
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "put 0 0 0 1\r\na\r\n" {
		t.Errorf("expected put command, got %q", berr.Cmd)
	}

	if berr.Reply != "INSERTED 1 2\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != BadReply {
		t.Fatalf("expected beanstalk.BadReply, got %v", berr.Error)
	}
}

func TestPutReplyNotEnoughArgs(t *testing.T) {
	rw, _ := responder("INSERTED\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("default").Put("a", 0, 0, 0)

	if id != 0 {
		t.Error("expected id 0, got", id)
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "put 0 0 0 1\r\na\r\n" {
		t.Errorf("expected put command, got %q", berr.Cmd)
	}

	if berr.Reply != "INSERTED\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != BadReply {
		t.Fatalf("expected beanstalk.BadReply, got %v", berr.Error)
	}
}

func TestPutReplyBadInteger(t *testing.T) {
	rw, _ := responder("INSERTED x\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("default").Put("a", 0, 0, 0)

	if id != 0 {
		t.Error("expected id 0, got", id)
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "put 0 0 0 1\r\na\r\n" {
		t.Errorf("expected put command, got %q", berr.Cmd)
	}

	if berr.Reply != "INSERTED x\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != BadReply {
		t.Fatalf("expected beanstalk.BadReply, got %v", berr.Error)
	}
}

func TestPutReplyInternalError(t *testing.T) {
	rw, _ := responder("INTERNAL_ERROR\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("default").Put("a", 0, 0, 0)

	if id != 0 {
		t.Error("expected id 0, got", id)
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "put 0 0 0 1\r\na\r\n" {
		t.Errorf("expected put command, got %q", berr.Cmd)
	}

	if berr.Reply != "INTERNAL_ERROR\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != InternalError {
		t.Fatalf("expected beanstalk.InternalError, got %v", berr.Error)
	}
}

func TestStripTab(t *testing.T) {
	rw, buf := responder("INSERTED 1\t\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("default").Put("a", 0, 0, 0)

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if id != 1 {
		t.Error("expected id 1, got", id)
	}

	if buf.String() != "put 0 0 0 1\r\na\r\n" {
		t.Errorf("expected put command, got %q", buf.String())
	}
}

func TestStripCR(t *testing.T) {
	rw, buf := responder("INSERTED 1\r\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("default").Put("a", 0, 0, 0)

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if id != 1 {
		t.Error("expected id 1, got", id)
	}

	if buf.String() != "put 0 0 0 1\r\na\r\n" {
		t.Errorf("expected put command, got %q", buf.String())
	}
}

func TestPut(t *testing.T) {
	rw, buf := responder("INSERTED 1\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("default").Put("a", 0, 0, 0)

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if id != 1 {
		t.Error("expected id 1, got", id)
	}

	if buf.String() != "put 0 0 0 1\r\na\r\n" {
		t.Errorf("expected put command, got %q", buf.String())
	}
}

func TestPut2(t *testing.T) {
	rw, buf := responder("INSERTED 2\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("default").Put("a", 0, 0, 0)

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if id != 2 {
		t.Error("expected id 2, got", id)
	}

	if buf.String() != "put 0 0 0 1\r\na\r\n" {
		t.Errorf("expected put command, got %q", buf.String())
	}
}

func TestPutOtherTube(t *testing.T) {
	rw, buf := responder("USING foo\nINSERTED 1\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("foo").Put("a", 0, 0, 0)

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if id != 1 {
		t.Error("expected id 1, got", id)
	}

	if buf.String() != "use foo\r\nput 0 0 0 1\r\na\r\n" {
		t.Errorf("expected use/put command, got %q", buf.String())
	}
}

func TestPutUseFail(t *testing.T) {
	rw, buf := responder("INTERNAL_ERROR\nINSERTED 1\n")
	c := newConn("<fake>", rw)
	id, err := c.Tube("foo").Put("a", 0, 0, 0)

	if buf.String() != "use foo\r\nput 0 0 0 1\r\na\r\n" {
		t.Errorf("expected use/put command, got %q", buf.String())
	}

	if id != 0 {
		t.Error("expected id 0, got", id)
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "use foo\r\n" {
		t.Errorf("expected use command, got %q", berr.Cmd)
	}

	if berr.Reply != "INTERNAL_ERROR\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != InternalError {
		t.Fatalf("expected beanstalk.InternalError, got %v", berr.Error)
	}
}

func TestDelete(t *testing.T) {
	rw, buf := responder("DELETED\n")
	err := Job{1, "a", newConn("<fake>", rw)}.Delete()

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if buf.String() != "delete 1\r\n" {
		t.Errorf("expected delete command, got %q", buf.String())
	}
}


func TestDeleteNotFound(t *testing.T) {
	rw, _ := responder("NOT_FOUND\n")
	err := Job{1, "a", newConn("<fake>", rw)}.Delete()

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "delete 1\r\n" {
		t.Errorf("expected delete command, got %q", berr.Cmd)
	}

	if berr.Reply != "NOT_FOUND\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != NotFound {
		t.Fatalf("expected beanstalk.NotFound, got %v", berr.Error)
	}
}

func TestTouch(t *testing.T) {
	rw, buf := responder("TOUCHED\n")
	err := Job{1, "a", newConn("<fake>", rw)}.Touch()

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if buf.String() != "touch 1\r\n" {
		t.Errorf("expected touch command, got %q", buf.String())
	}
}

func TestTouchNotFound(t *testing.T) {
	rw, _ := responder("NOT_FOUND\n")
	err := Job{1, "a", newConn("<fake>", rw)}.Touch()

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "touch 1\r\n" {
		t.Errorf("expected touch command, got %q", berr.Cmd)
	}

	if berr.Reply != "NOT_FOUND\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != NotFound {
		t.Fatalf("expected beanstalk.NotFound, got %v", berr.Error)
	}
}

func TestStats(t *testing.T) {
	rw, buf := responder("OK 14\n---\na: 1\nx: y\n\r\n")
	c := newConn("<fake>", rw)
	stats, err := c.Stats()

	if buf.String() != "stats\r\n" {
		t.Errorf("expected stats command, got %q", buf.String())
	}

	if err != nil {
		t.Error("unexpected error", err)
	}

	if stats == nil {
		t.Fatal("stats is nil")
	}

	exp := map[string]string{"a": "1", "x": "y"}
	if !reflect.DeepEqual(stats, exp) {
		t.Errorf("stats doesn't match, got %#v", stats)
	}
}

func TestStatsJob(t *testing.T) {
	rw, buf := responder("OK 14\n---\na: 1\nx: y\n\r\n")
	stats, err := Job{1, "a", newConn("<fake>", rw)}.Stats()

	if buf.String() != "stats-job 1\r\n" {
		t.Errorf("expected stats-job command, got %q", buf.String())
	}

	if err != nil {
		t.Error("unexpected error", err)
	}

	if stats == nil {
		t.Fatal("stats is nil")
	}

	exp := map[string]string{"a": "1", "x": "y"}
	if !reflect.DeepEqual(stats, exp) {
		t.Errorf("stats doesn't match, got %#v", stats)
	}
}

func TestStatsTube(t *testing.T) {
	rw, buf := responder("OK 14\n---\na: 1\nx: y\n\r\n")
	c := newConn("<fake>", rw)
	stats, err := c.Tube("foo").Stats()

	if buf.String() != "stats-tube foo\r\n" {
		t.Errorf("expected stats-tube command, got %q", buf.String())
	}

	if err != nil {
		t.Error("unexpected error", err)
	}

	if stats == nil {
		t.Fatal("stats is nil")
	}

	exp := map[string]string{"a": "1", "x": "y"}
	if !reflect.DeepEqual(stats, exp) {
		t.Errorf("stats doesn't match, got %#v", stats)
	}
}

func TestPeekNotFound(t *testing.T) {
	rw, _ := responder("NOT_FOUND\n")
	c := newConn("<fake>", rw)
	j, err := c.Peek(1)

	if j != nil {
		t.Error("expected nil job")
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "peek 1\r\n" {
		t.Errorf("expected peek command, got %q", berr.Cmd)
	}

	if berr.Reply != "NOT_FOUND\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != NotFound {
		t.Fatalf("expected beanstalk.NotFound, got %v", berr.Error)
	}
}

func TestPeek(t *testing.T) {
	rw, _ := responder("FOUND 1 1\na\n")
	c := newConn("<fake>", rw)
	j, err := c.Peek(1)

	if err != nil {
		t.Error("unexpected error", err)
	}

	if j == nil {
		t.Fatal("job is nil")
	}

	if j.Id != 1 {
		t.Error("expedted id 1, got", j.Id)
	}

	if j.Body != "a" {
		t.Errorf("expedted body \"a\", got %q", j.Body)
	}
}

func TestPeekReplyNotEnoughArgs(t *testing.T) {
	rw, buf := responder("FOUND\na\n")
	c := newConn("<fake>", rw)
	j, err := c.Peek(1)

	if buf.String() != "peek 1\r\n" {
		t.Errorf("expected peek command, got %q", buf.String())
	}

	if j != nil {
		t.Errorf("unexpected job %#v", j)
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "peek 1\r\n" {
		t.Errorf("expected peek command, got %q", berr.Cmd)
	}

	if berr.Reply != "FOUND\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != BadReply {
		t.Fatalf("expected beanstalk.BadReply, got %v", berr.Error)
	}
}

func TestPeekReadyOtherTube(t *testing.T) {
	rw, buf := responder("USING foo\nFOUND 1 1\na\r\n")
	c := newConn("<fake>", rw)
	j, err := c.Tube("foo").PeekReady()

	if buf.String() != "use foo\r\npeek-ready\r\n" {
		t.Errorf("expected use/peek-ready command, got %q", buf.String())
	}

	if err != nil {
		t.Error("unexpected error", err)
	}

	if j == nil {
		t.Fatal("job is nil")
	}

	if j.Id != 1 {
		t.Error("expedted id 1, got", j.Id)
	}

	if j.Body != "a" {
		t.Errorf("expedted body \"a\", got %q", j.Body)
	}

}

func TestPeekDelayedOtherTube(t *testing.T) {
	rw, buf := responder("USING foo\nFOUND 1 1\na\r\n")
	c := newConn("<fake>", rw)
	j, err := c.Tube("foo").PeekDelayed()

	if buf.String() != "use foo\r\npeek-delayed\r\n" {
		t.Errorf("expected use/peek-delayed command, got %q", buf.String())
	}

	if err != nil {
		t.Error("unexpected error", err)
	}

	if j == nil {
		t.Fatal("job is nil")
	}

	if j.Id != 1 {
		t.Error("expedted id 1, got", j.Id)
	}

	if j.Body != "a" {
		t.Errorf("expedted body \"a\", got %q", j.Body)
	}

}

func TestPeekBuriedOtherTube(t *testing.T) {
	rw, buf := responder("USING foo\nFOUND 1 1\na\r\n")
	c := newConn("<fake>", rw)
	j, err := c.Tube("foo").PeekBuried()

	if buf.String() != "use foo\r\npeek-buried\r\n" {
		t.Errorf("expected use/peek-buried command, got %q", buf.String())
	}

	if err != nil {
		t.Error("unexpected error", err)
	}

	if j == nil {
		t.Fatal("job is nil")
	}

	if j.Id != 1 {
		t.Error("expedted id 1, got", j.Id)
	}

	if j.Body != "a" {
		t.Errorf("expedted body \"a\", got %q", j.Body)
	}

}

func TestPeekReadyNotFound(t *testing.T) {
	rw, _ := responder("NOT_FOUND\n")
	c := newConn("<fake>", rw)
	j, err := c.Tube("default").PeekReady()

	if j != nil {
		t.Error("expected nil job")
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "peek-ready\r\n" {
		t.Errorf("expected peek-ready command, got %q", berr.Cmd)
	}

	if berr.Reply != "NOT_FOUND\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != NotFound {
		t.Fatalf("expected beanstalk.NotFound, got %v", berr.Error)
	}
}

func TestPeekDelayedNotFound(t *testing.T) {
	rw, _ := responder("NOT_FOUND\n")
	c := newConn("<fake>", rw)
	j, err := c.Tube("default").PeekDelayed()

	if j != nil {
		t.Error("expected nil job")
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "peek-delayed\r\n" {
		t.Errorf("expected peek-delayed command, got %q", berr.Cmd)
	}

	if berr.Reply != "NOT_FOUND\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != NotFound {
		t.Fatalf("expected beanstalk.NotFound, got %v", berr.Error)
	}
}

func TestPeekBuriedNotFound(t *testing.T) {
	rw, _ := responder("NOT_FOUND\n")
	c := newConn("<fake>", rw)
	j, err := c.Tube("default").PeekBuried()

	if j != nil {
		t.Error("expected nil job")
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "peek-buried\r\n" {
		t.Errorf("expected peek-buried command, got %q", berr.Cmd)
	}

	if berr.Reply != "NOT_FOUND\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != NotFound {
		t.Fatalf("expected beanstalk.NotFound, got %v", berr.Error)
	}
}

func TestReserve(t *testing.T) {
	rw, buf := responder("RESERVED 1 1\na\r\n")
	c := newConn("<fake>", rw)
	names := []string{"default"}
	j, err := c.Tubes(names).Reserve()

	if buf.String() != "reserve-with-timeout 4000000000\r\n" {
		t.Errorf("expected reserve command, got %q", buf.String())
	}

	if err != nil {
		t.Error("unexpected error", err)
	}

	if j == nil {
		t.Fatal("job is nil")
	}

	if j.Id != 1 {
		t.Error("expedted id 1, got", j.Id)
	}

	if j.Body != "a" {
		t.Errorf("expedted body \"a\", got %q", j.Body)
	}
}

func TestReserveExtraTube(t *testing.T) {
	rw, buf := responder("WATCHING 2\nRESERVED 1 1\na\r\n")
	c := newConn("<fake>", rw)
	names := []string{"default", "foo"}
	j, err := c.Tubes(names).Reserve()

	if buf.String() != "watch foo\r\nreserve-with-timeout 4000000000\r\n" {
		t.Errorf("expected watch/reserve command, got %q", buf.String())
	}

	if err != nil {
		t.Error("unexpected error", err)
	}

	if j == nil {
		t.Fatal("job is nil")
	}

	if j.Id != 1 {
		t.Error("expedted id 1, got", j.Id)
	}

	if j.Body != "a" {
		t.Errorf("expedted body \"a\", got %q", j.Body)
	}
}

func TestReserveAlternateTube(t *testing.T) {
	rw, buf := responder("WATCHING 2\nWATCHING 1\nRESERVED 1 1\na\r\n")
	c := newConn("<fake>", rw)
	names := []string{"foo"}
	j, err := c.Tubes(names).Reserve()

	if buf.String() != "watch foo\r\nignore default\r\nreserve-with-timeout 4000000000\r\n" {
		t.Errorf("expected watch/ignore/reserve command, got %q", buf.String())
	}

	if err != nil {
		t.Error("unexpected error", err)
	}

	if j == nil {
		t.Fatal("job is nil")
	}

	if j.Id != 1 {
		t.Error("expedted id 1, got", j.Id)
	}

	if j.Body != "a" {
		t.Errorf("expedted body \"a\", got %q", j.Body)
	}
}

func TestParseDict(t *testing.T) {
	in := "---\na: 1\nx: y\n"
	got := parseDict(in)
	exp := map[string]string{"a": "1", "x": "y"}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("map doesn't match, got %#v", got)
	}
}

func TestParseDictMissingDocSep(t *testing.T) {
	in := "a: 1\nx: y\n"
	got := parseDict(in)
	exp := map[string]string{"a": "1", "x": "y"}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("map doesn't match, got %#v", got)
	}
}

func TestParseDictMissingFinalNewline(t *testing.T) {
	in := "---\na: 1\nx: y"
	got := parseDict(in)
	exp := map[string]string{"a": "1", "x": "y"}
	if !reflect.DeepEqual(got, exp) {
		t.Errorf("map doesn't match, got %#v", got)
	}
}

func TestKick(t *testing.T) {
	rw, buf := responder("KICKED 3\n")
	c := newConn("<fake>", rw)
	n, err := c.Tube("default").Kick(3)

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if n != 3 {
		t.Error("expected n 3, got", n)
	}

	if buf.String() != "kick 3\r\n" {
		t.Errorf("expected kick command, got %q", buf.String())
	}
}

func TestKickFewer(t *testing.T) {
	rw, buf := responder("KICKED 2\n")
	c := newConn("<fake>", rw)
	n, err := c.Tube("default").Kick(3)

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if n != 2 {
		t.Error("expected n 2, got", n)
	}

	if buf.String() != "kick 3\r\n" {
		t.Errorf("expected kick command, got %q", buf.String())
	}
}

func TestKickOtherTube(t *testing.T) {
	rw, buf := responder("USING foo\nKICKED 3\n")
	c := newConn("<fake>", rw)
	n, err := c.Tube("foo").Kick(3)

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if n != 3 {
		t.Error("expected n 3, got", n)
	}

	if buf.String() != "use foo\r\nkick 3\r\n" {
		t.Errorf("expected use/kick command, got %q", buf.String())
	}
}

func TestTubePause(t *testing.T) {
	rw, buf := responder("PAUSED\n")
	c := newConn("<fake>", rw)
	err := c.Tube("foo").Pause(3)

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if buf.String() != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", buf.String())
	}
}

func TestTubePauseNotFound(t *testing.T) {
	rw, buf := responder("NOT_FOUND\n")
	c := newConn("<fake>", rw)
	err := c.Tube("foo").Pause(3)

	if buf.String() != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", buf.String())
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", berr.Cmd)
	}

	if berr.Reply != "NOT_FOUND\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != NotFound {
		t.Fatalf("expected beanstalk.NotFound, got %v", berr.Error)
	}
}

func TestTubePauseInternalError(t *testing.T) {
	rw, buf := responder("INTERNAL_ERROR\n")
	c := newConn("<fake>", rw)
	err := c.Tube("foo").Pause(3)

	if buf.String() != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", buf.String())
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", berr.Cmd)
	}

	if berr.Reply != "INTERNAL_ERROR\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != InternalError {
		t.Fatalf("expected beanstalk.InternalError, got %v", berr.Error)
	}
}

func TestTubePauseOutOfMemory(t *testing.T) {
	rw, buf := responder("OUT_OF_MEMORY\n")
	c := newConn("<fake>", rw)
	err := c.Tube("foo").Pause(3)

	if buf.String() != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", buf.String())
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", berr.Cmd)
	}

	if berr.Reply != "OUT_OF_MEMORY\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != OutOfMemory {
		t.Fatalf("expected beanstalk.OutOfMemory, got %v", berr.Error)
	}
}

func TestTubeBadFormat(t *testing.T) {
	rw, buf := responder("BAD_FORMAT\n")
	c := newConn("<fake>", rw)
	err := c.Tube("foo").Pause(3)

	if buf.String() != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", buf.String())
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", berr.Cmd)
	}

	if berr.Reply != "BAD_FORMAT\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != BadFormat {
		t.Fatalf("expected beanstalk.BadFormat, got %v", berr.Error)
	}
}

func TestTubeUnknownCommand(t *testing.T) {
	rw, buf := responder("UNKNOWN_COMMAND\n")
	c := newConn("<fake>", rw)
	err := c.Tube("foo").Pause(3)

	if buf.String() != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", buf.String())
	}

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "pause-tube foo 3\r\n" {
		t.Errorf("expected pause-tube command, got %q", berr.Cmd)
	}

	if berr.Reply != "UNKNOWN_COMMAND\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != UnknownCommand {
		t.Fatalf("expected beanstalk.UnknownCommand, got %v", berr.Error)
	}
}

func TestBury(t *testing.T) {
	rw, buf := responder("BURIED\n")
	err := Job{1, "a", newConn("<fake>", rw)}.Bury(8)

	if err != nil {
		t.Error("got unexpected error:\n  ", err)
	}

	if buf.String() != "bury 1 8\r\n" {
		t.Errorf("expected bury command, got %q", buf.String())
	}
}

func TestBuryNotFound(t *testing.T) {
	rw, _ := responder("NOT_FOUND\n")
	err := Job{1, "a", newConn("<fake>", rw)}.Bury(8)

	if err == nil {
		t.Fatal("expected error, got none")
	}

	berr, ok := err.(Error)

	if !ok {
		t.Fatalf("expected beanstalk.Error, got %T", err)
	}

	if berr.Cmd != "bury 1 8\r\n" {
		t.Errorf("expected bury command, got %q", berr.Cmd)
	}

	if berr.Reply != "NOT_FOUND\n" {
		t.Errorf("reply was %q", berr.Reply)
	}

	if berr.Error != NotFound {
		t.Fatalf("expected beanstalk.NotFound, got %v", berr.Error)
	}
}

