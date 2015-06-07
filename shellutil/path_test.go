package shellutil

import (
	"io/ioutil"
	"testing"
)

func TestPath(t *testing.T) {
	t.Log("testing path functionality")
	tmpdir, _ := ioutil.TempDir("/tmp", "test_shellutil_")
	for _, dir := range []string{"/tmp", "/usr/local", tmpdir} {
		if IsFile(dir) || !IsDirectory(dir) {
			t.Error(dir, "is a directory")
		}
	}
	tmpfile, _ := ioutil.TempFile(tmpdir, "a_file_")
	for _, file := range []string{"/etc/hosts", "/etc/passwd", tmpfile.Name()} {
		if !IsFile(file) || IsDirectory(file) {
			t.Error(file, "is a file")
		}
	}
}

func TestInputOutput(t *testing.T) {
	t.Log("testing input/output functionality")
	R := [][4]string{
		{"/home/apefind/test.txt", "/tmp", ".csv", "/tmp/test.csv"},
		{"/home/apefind/test.txt", "/tmp/_not_a_dir_", ".csv", "/tmp/_not_a_dir_.csv"},
		{"/home/apefind/test.txt", "/tmp/test.xls", "", "/tmp/test.xls"},
		{"/home/apefind/test", "", ".csv", "/home/apefind/test.csv"},
		{"/home/apefind/test", "", ".mp3", "/home/apefind/test.mp3"},
	}
	for _, r := range R {
		if GetOutputFilename(r[0], r[1], r[2]) != r[3] {
			t.Error("expected", r[3], "got", GetOutputFilename(r[0], r[1], r[2]))
		}
	}
}
