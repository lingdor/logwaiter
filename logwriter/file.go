package logwriter

import (
	"bufio"
	"fmt"
	"github.com/lingdor/logwaiter/common"
	"io"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
)

const BufferSize = 1024 * 2
const WriteMaxRetry = 3

var writer unsafe.Pointer

var indexNumber inc
var lastPath string

type LogWriter interface {
	io.Writer
	io.Closer
	Flush() error
	UnWrap() LogWriter
}
type logWriter struct {
	writer *bufio.Writer
	file   *os.File
}

func (l *logWriter) Write(p []byte) (n int, err error) {
	return l.writer.Write(p)
}
func (l *logWriter) Close() error {
	return l.file.Close()
}
func (l *logWriter) Flush() error {
	return l.writer.Flush()
}
func (l *logWriter) UnWrap() LogWriter {
	return l
}

type wrapLogWriter struct {
	writer LogWriter
}

func (l *wrapLogWriter) Write(p []byte) (n int, err error) {
	write, err := l.writer.Write(p)
	l.writer.Flush()
	return write, err
}
func (l *wrapLogWriter) Close() error {
	return l.writer.Close()
}
func (l *wrapLogWriter) Flush() error {
	return l.writer.Flush()
}
func (l *wrapLogWriter) UnWrap() LogWriter {
	return l.writer
}

func init() {
}

func LoadWriter(writePath string, splitLength int64) {
	writePath = replaceDate(writePath)
	if writePath != lastPath {
		indexNumber.Reset()
	} else {
		return
	}
	path := writePath
	for ; ; indexNumber.Add() {
		if indexNumber.Get() != 0 {
			path = fmt.Sprintf("%s.%d", writePath, indexNumber.Get())
		}
		info, err := os.Stat(path)
		if err == nil && info.IsDir() {
			continue
		}
		if os.IsNotExist(err) {
			break
		}
		if info.Size() < splitLength {
			break
		}
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		//open file failed
		indexNumber.Add()
		LoadWriter(writePath, splitLength)
		return
	}

	newFileWriter := bufio.NewWriterSize(file, BufferSize)
	newFile := wrapLogWriter{
		writer: &logWriter{
			file:   file,
			writer: newFileWriter,
		},
	}
	log.Println(fmt.Sprintf("log path chanage: %s -> %s", lastPath, path))
	lastPath = path
	swapWriter(&newFile)
}

func swapWriter(w LogWriter) {
	old := getWriter()
	newPointer := (unsafe.Pointer)(&w)
	atomic.SwapPointer(&writer, newPointer)
	if old != nil && old.UnWrap() != w.UnWrap() {
		time.Sleep(time.Millisecond)
		old.Flush()
		old.Close()
	}
}

func getWriter() LogWriter {
	if writer == nil {
		return nil
	}
	pointer := atomic.LoadPointer(&writer)
	val := (*LogWriter)(pointer)
	return *val
}

func replaceDate(writePath string) string {
	now := time.Now()
	year := fmt.Sprintf("%04d", now.Year())
	writePath = strings.Replace(writePath, "#Y", year, -1)
	shortYear := year[2:]
	writePath = strings.Replace(writePath, "#y", shortYear, -1)
	writePath = strings.Replace(writePath, "#m", fmt.Sprintf("%02d", now.Month()), -1)
	writePath = strings.Replace(writePath, "#d", fmt.Sprintf("%02d", now.Day()), -1)
	writePath = strings.Replace(writePath, "#H", fmt.Sprintf("%02d", now.Hour()), -1)
	writePath = strings.Replace(writePath, "#i", fmt.Sprintf("%02d", now.Minute()), -1)
	//don't support  file name with second level.
	return writePath
}

func WriteLine(line []byte, retry int) {
	writer := getWriter()
	if writer == nil {
		return
	}
	_, err := writer.Write(append(line, byte('\n')))
	if err != nil {
		if retry < WriteMaxRetry {
			WriteLine(line, retry+1)
			return
		}
		common.CheckPanic(err)
	}
}

func Flush() {
	writer := getWriter()
	if writer == nil {
		return
	}
	writer.Flush()
}

func Close() {
	writer := getWriter()
	if writer == nil {
		return
	}
	writer.Flush()
	writer.Close()
}

func CheckWrap() {
	writer := getWriter()
	if writer == nil {
		return
	}
	unwrap := writer.UnWrap()
	if unwrap != writer {
		return
	}
	swapWriter(&wrapLogWriter{writer: unwrap})
}
func CheckUnWrap() {
	writer := getWriter()
	if writer == nil {
		return
	}
	unwrap := writer.UnWrap()
	if unwrap == writer {
		return
	}
	swapWriter(unwrap)
}
