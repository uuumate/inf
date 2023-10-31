package rolling

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type RollingFormat string

const (
	MonthlyRolling  RollingFormat = "200601"
	DailyRolling                  = "20060102"
	HourlyRolling                 = "2006010215"
	MinutelyRolling               = "200601021504"
	SecondlyRolling               = "20060102150405"
)

type File struct {
	fileMu     sync.RWMutex
	file       *os.File
	filePrefix string
	fileName   string
	filePath   string
	ctx        context.Context

	rollingFormat RollingFormat
	valueCh       chan string
	closed        int64
	writeCount    int64
}

func (f *File) SetRollingFormat(rollingFormat RollingFormat) {
	f.rollingFormat = rollingFormat
}

func NewRollingFile(filePath, filePrefix string) *File {
	if len(filePath) != 0 {
		_ = os.Mkdir(filePath, os.ModePerm)
	}

	rf := &File{
		filePrefix:    filePrefix,
		filePath:      filePath,
		rollingFormat: DailyRolling,
		valueCh:       make(chan string, 1024),
	}

	go rf.flush()
	return rf
}

func (f *File) Close() {
	atomic.SwapInt64(&f.closed, 1)
	close(f.valueCh)
	_ = f.Sync()
}

func (f *File) Sync() error {
	if atomic.LoadInt64(&f.writeCount) == 0 {
		return nil
	}
	_ = f.getWriter().Sync()
	return nil
}

func (f *File) makeCurrentFileName() string {
	return fmt.Sprintf("%s/%s-%s.log", f.filePath, f.filePrefix, time.Now().Format(string(f.rollingFormat)))
}

func (f *File) getWriter() *os.File {
	f.fileMu.RLock()

	fileName := f.makeCurrentFileName()

	if fileName == f.fileName {
		f.fileMu.RUnlock()
		return f.file
	}

	f.fileMu.RUnlock()
	_ = f.file.Sync()
	_ = f.file.Close()

	f.fileMu.Lock()
	f.fileName = fileName
	f.file, _ = os.OpenFile(fileName, os.O_APPEND|os.O_RDWR|os.O_CREATE, os.ModePerm)
	f.fileMu.Unlock()

	return f.file
}

func (f *File) Write(bytes []byte) (n int, err error) {
	if atomic.LoadInt64(&f.closed) == 1 {
		return 0, ErrRollingFileIsAlreadyClosed
	}
	f.valueCh <- string(bytes)
	return len(bytes), nil
}

func (f *File) flush() {
	for {
		v, ok := <-f.valueCh
		if !ok {
			_ = f.Sync()
			break
		}
		_, _ = f.getWriter().Write([]byte(v))
		atomic.AddInt64(&f.writeCount, 1)
	}
}
