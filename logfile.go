package gobase

import (
	"fmt"
	"os"
	"path"
	"runtime"
	"syscall"
	"time"
	"sync"
)

const (
	FATAL int8 = 1
	ERROR int8 = 2
	WARN  int8 = 3
	INFO  int8 = 4
	DEBUG int8 = 5
)

type Logfile struct {
	filename  string
	loglevel  int8
	prefix    string
	out       *os.File
	checktime int64
	fileino   uint64
	locker	  sync.Mutex
}

func New(filename string, loglevel int8) *Logfile {
	log := &Logfile{filename: filename, loglevel: loglevel}
	log.Open()
	return log
}

func (l *Logfile) Open() {
	if l.filename == "stdout" {
		l.out = os.Stdout
	} else if l.filename == "stderr" {
		l.out = os.Stderr
	} else {
		if l.out != nil {
			l.out.Close()
		}
		f, err := os.OpenFile(l.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Printf("open logfile %s error: %s", l.filename, err)
			os.Exit(1)
		}
		l.checktime = time.Now().Unix()
		l.out = f

		fi, err := os.Stat(l.filename)
		if err == nil {
			st := fi.Sys().(*syscall.Stat_t)
			l.fileino = st.Ino
		}

	}

}

func (l *Logfile) SetPrefix(prefix string) {
	l.prefix = prefix
}

func file_inode(filename string) uint64 {
	fi, err := os.Stat(filename)
	if err == nil {
		st := fi.Sys().(*syscall.Stat_t)
		return st.Ino
	}
	return 0
}

func (l *Logfile) write(levelid int8, format string, v ...interface{}) {

	if levelid > l.loglevel {
		return
	}
	levelstr := [6]string{"", "F", "E", "W", "I", "D"}
	levelcolor := [6]string{"", "\033[31m", "\033[35m", "\033[33m", "\033[36m", ""}

	_, file, line, ok := runtime.Caller(2)
	if !ok {
		file = "???"
		line = 0
	}
	t := time.Now()
	s1 := fmt.Sprintf("%d%02d%02d %02d:%02d:%02d.%03d %s:%d [%s] ",
		t.Year()%2000, t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000,
		path.Base(file), line, levelstr[levelid])

	s2 := fmt.Sprintf(format, v...)
	if len(l.prefix) > 0 {
		s2 = l.prefix + s2
	}
	var s string
	if l.out == os.Stdout {
		s = levelcolor[levelid] + s1 + s2 + "\033[0m"
	} else {
		s = s1 + s2
	}
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	l.out.WriteString(s)

	// 每分钟检查日志文件是否有变动
	now := time.Now().Unix()
	if now-l.checktime > 60 {
		ino := file_inode(l.filename)
		if ino != l.fileino {
			l.locker.Lock()
			ino := file_inode(l.filename)
			if ino != l.fileino {
				l.Open()
			}
			l.locker.Unlock()
		}
		l.checktime = now
	}
}

func (l *Logfile) Debug(format string, v ...interface{}) {
	l.write(DEBUG, format, v...)
}

func (l *Logfile) Info(format string, v ...interface{}) {
	l.write(INFO, format, v...)
}

func (l *Logfile) Warn(format string, v ...interface{}) {
	l.write(WARN, format, v...)
}

func (l *Logfile) Error(format string, v ...interface{}) {
	l.write(ERROR, format, v...)
}

func (l *Logfile) Fatal(format string, v ...interface{}) {
	l.write(FATAL, format, v...)
	os.Exit(1)
}

func Warn(format string, v ...interface{}) {
	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "???"
		line = 0
	}
	t := time.Now()
	s1 := fmt.Sprintf("%d%02d%02d %02d:%02d:%02d.%03d %s:%d [W] ",
		t.Year()%2000, t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second(), t.Nanosecond()/1000000,
		path.Base(file), line)

	s2 := fmt.Sprintf(format, v...)
	s := s1 + s2
	if s[len(s)-1] != '\n' {
		s += "\n"
	}
	os.Stderr.WriteString(s)
}
