package main

import (
	"container/list"
	"io/ioutil"
	"os"
	"strings"
	"testing"
)

func ok(err error) {
	if err != nil {
		panic(err)
	}
}

// TestList test
func TestList(t *testing.T) {
	temp, err := ioutil.TempFile("", "napa")
	ok(err)
	defer temp.Close()
	defer os.Remove(temp.Name())

	ls := list.New()
	ls.PushBack("abcd-efgh")
	ls.PushBack("wxyz-lmnop")
	ls.PushBack("asdf-qwop")

	var buffer strings.Builder
	for e := ls.Front(); e != nil; e = e.Next() {
		v := e.Value.(string)
		buffer.WriteString(v)
		buffer.WriteString("\n")
	}
	writeBytes(temp.Name(), []byte(buffer.String()))

	read := readList(temp.Name())
	expectString(t, read[0], "abcd-efgh")
	expectString(t, read[1], "wxyz-lmnop")
	expectString(t, read[2], "asdf-qwop")
}

// TestFileIO test
func TestFileIO(t *testing.T) {
	tempA, err := ioutil.TempFile("", "napa")
	ok(err)
	tempA.Close()
	defer os.Remove(tempA.Name())

	tempB, err := ioutil.TempFile("", "napa")
	ok(err)
	tempB.Close()
	defer os.Remove(tempB.Name())

	tempSwapA, err := ioutil.TempFile("", "napa")
	ok(err)
	tempSwapA.Close()
	defer os.Remove(tempSwapA.Name())

	tempSwapB, err := ioutil.TempFile("", "napa")
	ok(err)
	tempSwapB.Close()
	defer os.Remove(tempSwapB.Name())

	err = writeBytes(tempA.Name(), []byte("contents of temp a\n"))
	ok(err)
	err = writeBytes(tempB.Name(), []byte("temp b contents\n"))
	ok(err)
	err = copyFile(tempA.Name(), tempSwapA.Name())
	ok(err)
	err = copyFile(tempB.Name(), tempSwapB.Name())
	ok(err)
	err = renameFile(tempB.Name(), tempA.Name())
	ok(err)

	var read []string

	read = readList(tempA.Name())
	expectString(t, read[0], "temp b contents")

	_, err = ioutil.ReadFile(tempB.Name())
	if os.IsExist(err) {
		t.Error("file should not exist", tempB.Name())
	}

	read = readList(tempSwapA.Name())
	expectString(t, read[0], "contents of temp a")

	read = readList(tempSwapB.Name())
	expectString(t, read[0], "temp b contents")
}
