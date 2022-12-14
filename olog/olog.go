package olog

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"sync/atomic"
	"time"
)

// These flags define which text to prefix to each log entry generated by the Logger.
// Bits are or'ed together to control what's printed.
// With the exception of the Lmsgprefix flag, there is no
// control over the order they appear (the order listed here)
// or the format they present (as described in the comments).
// The prefix is followed by a colon only when Llongfile or Lshortfile
// is specified.
// For example, flags Ldate | Ltime (or LstdFlags) produce,
//
//	2009/01/23 01:23:23 message
//
// while flags Ldate | Ltime | Lmicroseconds | Llongfile produce,
//
//	2009/01/23 01:23:23.123123 /a/b/c/d.go:23: message
const (
	Ldate         = 1 << iota     // the date in the local time zone: 2009/01/23
	Ltime                         // the time in the local time zone: 01:23:23
	Lmicroseconds                 // microsecond resolution: 01:23:23.123123.  assumes Ltime.
	Llongfile                     // full file name and line number: /a/b/c/d.go:23
	Lshortfile                    // final file name element and line number: d.go:23. overrides Llongfile
	LUTC                          // if Ldate or Ltime is set, use UTC rather than the local time zone
	Lmsgprefix                    // move the "prefix" from the beginning of the line to before the message
	LstdFlags     = Ldate | Ltime // initial values for the standard logger
)

// log level
const (
	LEVEL_OFF   = 0
	LEVEL_FATAL = 1
	LEVEL_ERROR = 2
	LEVEL_WARN  = 3
	LEVEL_INFO  = 4
	LEVEL_DEBUG = 5
	LEVEL_TRACE = 6
	LEVEL_ALL   = 7
)

type Olog struct {
	Level int
	Year  int
	Month int
	Day   int
	out   io.Writer // destination for output
}

type Logger struct {
	mu        sync.Mutex  // ensures atomic writes; protects the following fields
	prefix    string      // prefix on each line to identify the logger (but see Lmsgprefix)
	flag      int         // properties
	out       io.Writer   // destination for output
	buf       []byte      // for accumulating text to write
	isDiscard atomic.Bool // atomic boolean: whether out == io.Discard
}

func (o *Olog) Update() {
	Time := time.Now()
	var filePath string

	if Time.Local().Year() != o.Year || int(Time.Local().Month()) != o.Month || Time.Local().Day() != o.Day {
		o.Year = Time.Local().Year()
		o.Month = int(Time.Local().Month())
		o.Day = Time.Local().Day()
		filePath = fmt.Sprintf("./logs/%v_%v_%v.log", o.Year, o.Month, o.Day)
	}
	_, err := os.Stat(filePath)
	if err != nil || os.IsNotExist(err) {
		paths, _ := filepath.Split(filePath)
		if err := os.MkdirAll(paths, 0777); err != nil {
			log.Fatalln(err)
		}
	}

	logFile, err := os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		log.Fatalln("open log file failed, err:", err)
	}
	// defer logFile.Close()
	writers := []io.Writer{logFile, os.Stdout}
	o.out = io.MultiWriter(writers...)
}

// FATAL
func (o *Olog) FATAL(v ...any) {
	l := New(o.out, "[FATAL] ", Ldate|Ltime|Llongfile)
	if LEVEL_FATAL <= o.Level {
		l.Output(2, fmt.Sprintln(v...))
		os.Exit(1)
	}
}

// ERROR
func (o *Olog) ERROR(v ...any) {
	l := New(o.out, "[ERROR] ", Ldate|Ltime|Llongfile)
	if LEVEL_ERROR <= o.Level {
		if l.isDiscard.Load() {
			return
		}
		l.Output(2, fmt.Sprintln(v...))
	}
}

// WARN
func (o *Olog) WARN(v ...any) {
	l := New(o.out, "[WARN] ", Ldate|Ltime|Llongfile)
	if LEVEL_WARN <= o.Level {
		if l.isDiscard.Load() {
			return
		}
		l.Output(2, fmt.Sprintln(v...))
	}
}

// INFO
func (o *Olog) INFO(v ...any) {
	l := New(o.out, "[INFO] ", Ldate|Ltime|Llongfile)
	if LEVEL_INFO <= o.Level {
		if l.isDiscard.Load() {
			return
		}
		l.Output(2, fmt.Sprintln(v...))
	}
}

// DEBUG
func (o *Olog) DEBUG(v ...any) {
	l := New(o.out, "[DEBUG] ", Ldate|Ltime|Llongfile)
	if LEVEL_DEBUG <= o.Level {
		if l.isDiscard.Load() {
			return
		}
		l.Output(2, fmt.Sprintln(v...))
	}
}

// TRACE
func (o *Olog) TRACE(v ...any) {
	l := New(o.out, "[TRACE] ", Ldate|Ltime|Llongfile)
	if LEVEL_TRACE <= o.Level {
		if l.isDiscard.Load() {
			return
		}
		l.Output(2, fmt.Sprintln(v...))
	}
}

// New creates a new Logger. The out variable sets the
// destination to which log data will be written.
// The prefix appears at the beginning of each generated log line, or
// after the log header if the Lmsgprefix flag is provided.
// The flag argument defines the logging properties.
func New(out io.Writer, prefix string, flag int) *Logger {
	l := &Logger{out: out, prefix: prefix, flag: flag}
	if out == io.Discard {
		l.isDiscard.Store(true)
	}
	return l
}

// Output writes the output for a logging event. The string s contains
// the text to print after the prefix specified by the flags of the
// Logger. A newline is appended if the last character of s is not
// already a newline. Calldepth is used to recover the PC and is
// provided for generality, although at the moment on all pre-defined
// paths it will be 2.
func (l *Logger) Output(calldepth int, s string) error {
	now := time.Now() // get this early.
	var file string
	var line int
	l.mu.Lock()
	defer l.mu.Unlock()
	if l.flag&(Lshortfile|Llongfile) != 0 {
		// Release lock while getting caller info - it's expensive.
		l.mu.Unlock()
		var ok bool
		_, file, line, ok = runtime.Caller(calldepth)
		if !ok {
			file = "???"
			line = 0
		}
		l.mu.Lock()
	}
	l.buf = l.buf[:0]
	l.formatHeader(&l.buf, now, file, line)
	l.buf = append(l.buf, s...)
	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

// formatHeader writes log header to buf in following order:
//   - l.prefix (if it's not blank and Lmsgprefix is unset),
//   - date and/or time (if corresponding flags are provided),
//   - file and line number (if corresponding flags are provided),
//   - l.prefix (if it's not blank and Lmsgprefix is set).
func (l *Logger) formatHeader(buf *[]byte, t time.Time, file string, line int) {
	if l.flag&Lmsgprefix == 0 {
		*buf = append(*buf, l.prefix...)
	}
	if l.flag&(Ldate|Ltime|Lmicroseconds) != 0 {
		if l.flag&LUTC != 0 {
			t = t.UTC()
		}
		if l.flag&Ldate != 0 {
			year, month, day := t.Date()
			itoa(buf, year, 4)
			*buf = append(*buf, '/')
			itoa(buf, int(month), 2)
			*buf = append(*buf, '/')
			itoa(buf, day, 2)
			*buf = append(*buf, ' ')
		}
		if l.flag&(Ltime|Lmicroseconds) != 0 {
			hour, min, sec := t.Clock()
			itoa(buf, hour, 2)
			*buf = append(*buf, ':')
			itoa(buf, min, 2)
			*buf = append(*buf, ':')
			itoa(buf, sec, 2)
			if l.flag&Lmicroseconds != 0 {
				*buf = append(*buf, '.')
				itoa(buf, t.Nanosecond()/1e3, 6)
			}
			*buf = append(*buf, ' ')
		}
	}
	if l.flag&(Lshortfile|Llongfile) != 0 {
		if l.flag&Lshortfile != 0 {
			short := file
			for i := len(file) - 1; i > 0; i-- {
				if file[i] == '/' {
					short = file[i+1:]
					break
				}
			}
			file = short
		}
		*buf = append(*buf, file...)
		*buf = append(*buf, ':')
		itoa(buf, line, -1)
		*buf = append(*buf, ": "...)
	}
	if l.flag&Lmsgprefix != 0 {
		*buf = append(*buf, l.prefix...)
	}
}

// Cheap integer to fixed-width decimal ASCII. Give a negative width to avoid zero-padding.
func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
