package logwriter

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync/atomic"
	"time"
	"unsafe"
)

const BufferSize = 100

var writer unsafe.Pointer

var indexNumber inc
var lastPath string

type logWriter struct {
	writer *bufio.Writer
	file   *os.File
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
		LoadWriter(writePath, splitLength)
		return
	}

	newFileWriter := bufio.NewWriterSize(file, BufferSize)
	newFile := logWriter{
		file:   file,
		writer: newFileWriter,
	}
	swapWriter(&newFile)
}

func swapWriter(w *logWriter) {
	old := getWriter()
	newPointer := (unsafe.Pointer)(w)
	atomic.SwapPointer(&writer, newPointer)
	if old != nil {
		old.writer.Flush()
		old.file.Close()
	}
}

func getWriter() *logWriter {
	if writer == nil {
		return nil
	}
	pointer := atomic.LoadPointer(&writer)
	return (*logWriter)(pointer)
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

func WriteLine(line string) {
	writer := getWriter()
	if writer == nil {
		return
	}
	writer.writer.Write([]byte(line))
	writer.writer.Write([]byte{'\n'})
}

func Flush() {
	writer := getWriter()
	if writer == nil {
		return
	}
	writer.writer.Flush()
}

func Close() {
	writer := getWriter()
	if writer == nil {
		return
	}
	writer.writer.Flush()
	writer.file.Close()
}
