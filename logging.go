package logging

import (
	"encoding/json"
	"fmt"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
	"time"

	"github.com/fatih/color"
	"gopkg.in/natefinch/lumberjack.v2"
)

var logger = new(lumberjack.Logger)

// Rotate rotate the output log file
func Rotate() (err error) {
	if err = logger.Rotate(); err != nil {
		return
	}
	return
}

var logLevel Level

type Config struct {
	LogLevel   Level
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	LocalTime  bool
	Compress   bool
}

// Setting set the logger output method
func Setting(conf Config) {
	logger = &lumberjack.Logger{
		Filename:   conf.Filename,
		MaxSize:    conf.MaxSize,
		MaxBackups: conf.MaxBackups,
		MaxAge:     conf.MaxAge,
		LocalTime:  conf.LocalTime,
		Compress:   conf.Compress,
	}

	logLevel = conf.LogLevel
}

// Fields ..
type Fields map[string]interface{}

// inst ..
type Inst struct {
	fields Fields
	trace  uint
	msg    []interface{}
	time   string
	level  Level
}

// Entry ..
type Entry interface {
	Panic(...interface{})
	Panicf(string, ...interface{})
	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Error(...interface{})
	Errorf(string, ...interface{})
	Warn(...interface{})
	Warnf(string, ...interface{})
	Warningf(string, ...interface{})
	Info(...interface{})
	Infof(string, ...interface{})
	Debug(...interface{})
	Debugf(string, ...interface{})
	withFields() *Inst
}

// WithFields add more field
func WithFields(fields Fields) (entry Entry) {
	entry = &Inst{fields: fields, trace: 2}
	entry.withFields()
	return
}

// WithFields ..
func (i *Inst) withFields() *Inst {
	for k, v := range i.fields {
		i.fields[k] = v
	}
	return i
}

func Panicf(format string, str ...interface{}) {
	i := &Inst{}
	i.Panicf(format, str...)
}

// Panicf ..
func (i *Inst) Panicf(format string, str ...interface{}) {
	if len(str) != 0 {
		i.msg = []interface{}{fmt.Sprintf(format, str...)}
	} else {
		i.msg = []interface{}{format}
	}
	i.level = PanicLevel
	i.output()
}

// Panic ..
func Panic(str ...interface{}) {
	i := &Inst{}
	i.Panic(str)
}

// Panic ..
func (i *Inst) Panic(str ...interface{}) {
	i.msg = str
	i.level = PanicLevel
	i.output()
}

// Fatalf ..
func Fatalf(format string, str ...interface{}) {
	i := &Inst{}
	i.Fatalf(format, str...)
}

// Fatalf ..
func (i *Inst) Fatalf(format string, str ...interface{}) {
	if len(str) != 0 {
		i.msg = []interface{}{fmt.Sprintf(format, str...)}
	} else {
		i.msg = []interface{}{format}
	}
	i.level = FatalLevel
	i.output()
}

// Fatal ..
func Fatal(str ...interface{}) {
	i := &Inst{}
	i.Fatal(str)
}

// Fatal ..
func (i *Inst) Fatal(str ...interface{}) {
	i.msg = str
	i.level = FatalLevel
	i.output()
}

// Errorf ..
func Errorf(format string, str ...interface{}) {
	i := &Inst{}
	i.Errorf(format, str...)
}

// Errorf ..
func (i *Inst) Errorf(format string, str ...interface{}) {
	if len(str) != 0 {
		i.msg = []interface{}{fmt.Sprintf(format, str...)}
	} else {
		i.msg = []interface{}{format}
	}
	i.level = ErrorLevel
	i.output()
}

// Error ..
func Error(str ...interface{}) {
	i := &Inst{}
	i.Error(str)
}

// Error ..
func (i *Inst) Error(str ...interface{}) {
	i.msg = str
	i.level = ErrorLevel
	i.output()
}

// Warnf ..
func Warnf(format string, str ...interface{}) {
	i := &Inst{}
	i.Warnf(format, str...)
}

// Warningf ..
func Warningf(format string, str ...interface{}) {
	Warnf(format, str...)
}

// Warnf ..
func (i *Inst) Warnf(format string, str ...interface{}) {
	if len(str) != 0 {
		i.msg = []interface{}{fmt.Sprintf(format, str...)}
	} else {
		i.msg = []interface{}{format}
	}
	i.level = WarnLevel
	i.output()
}

func (i *Inst) Warningf(format string, str ...interface{}) {
	i.Warnf(format, str...)
}

// Warn ..
func Warn(str ...interface{}) {
	i := &Inst{}
	i.Warn(str)
}

// Warn ..
func (i *Inst) Warn(str ...interface{}) {
	i.msg = str
	i.level = WarnLevel
	i.output()
}

// Infof ..
func Infof(format string, str ...interface{}) {
	i := &Inst{}
	i.Infof(format, str...)
}

// Info ..
func (i *Inst) Infof(format string, str ...interface{}) {
	if len(str) != 0 {
		i.msg = []interface{}{fmt.Sprintf(format, str...)}
	} else {
		i.msg = []interface{}{format}
	}
	i.level = InfoLevel
	i.output()
}

// Info ..
func Info(str ...interface{}) {
	i := &Inst{}
	i.Info(str)
}

// Info ..
func (i *Inst) Info(str ...interface{}) {
	i.msg = str
	i.level = InfoLevel
	i.output()
}

// Debugf ..
func Debugf(format string, str ...interface{}) {
	i := &Inst{}
	i.Debugf(format, str...)
}

// Debugf ..
func (i *Inst) Debugf(format string, str ...interface{}) {
	if len(str) != 0 {
		i.msg = []interface{}{fmt.Sprintf(format, str...)}
	} else {
		i.msg = []interface{}{format}
	}
	i.level = DebugLevel
	i.output()
}

// Debug ..
func Debug(str ...interface{}) {
	i := &Inst{}
	i.Debug(str)
}

// Debug ..
func (i *Inst) Debug(str ...interface{}) {
	i.msg = str
	i.level = DebugLevel
	i.output()
}

func (i *Inst) output() {
	var colorFun func(...interface{}) string
	var waitWrite []byte
	if i.level < logLevel {
		return
	}
	switch i.level {
	case DebugLevel:
		colorFun = color.New(color.FgHiWhite).SprintFunc()
	case WarnLevel:
		colorFun = color.New(color.FgHiYellow).SprintFunc()
	case ErrorLevel, FatalLevel, PanicLevel:
		colorFun = color.New(color.FgHiRed).SprintFunc()
	default:
		colorFun = color.New(color.FgHiBlue).SprintFunc()
	}
	levelText := strings.ToUpper(i.level.String())[0:4]
	var output string
	if i.fields == nil {
		i.fields = Fields{}
	}
	var trace = 3
	if i.trace == 2 {
		trace = 2
	}
	if _, file, line, ok := runtime.Caller(trace); ok {
		i.fields["_file"] = filepath.Base(file)
		i.fields["_line"] = line
	}

	t := time.Now()
	i.time = t.Format("15:04:05.000")
	i.fields["__time"] = t.Format("01-02T15:04:05.000")

	var keys []string
	for key := range i.fields {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	for _, key := range keys {
		yellow := color.New(color.FgYellow).SprintFunc()
		if key == "_line" || key == "_file" {
			output += fmt.Sprintf(" %s=%+v", yellow(key[1:]), i.fields[key])
		}
	}

	for _, key := range keys {
		green := color.New(color.FgGreen).SprintFunc()
		if key != "_line" && key != "_file" && key != "__time" {
			switch i.fields[key].(type) {
			case string:
				output += fmt.Sprintf(" %s=%+v", green(key), strings.TrimSpace(i.fields[key].(string)))
			default:
				output += fmt.Sprintf(" %s=%+v", green(key), i.fields[key])
			}
		}
	}

	msg := strings.TrimSuffix(strings.TrimPrefix(fmt.Sprint(i.msg...), "["), "]")

	if len(msg) > 1024*10 {
		msg = "msg is too long and cannot be display"
	}

	fmt.Printf("%s[%s] %-40v %s\n", colorFun(levelText), i.time, strings.TrimSpace(msg), output)

	i.fields["level"] = i.level.String()
	i.fields["msg"] = strings.TrimSuffix(strings.TrimPrefix(fmt.Sprint(i.msg...), "["), "]")
	waitWrite, _ = json.Marshal(i.fields)
	waitWrite = append(waitWrite, '\n')

	if logger != nil {
		if _, err := logger.Write(waitWrite); err != nil {
			logger = nil
			Error("Cannot write log to file.")
		}
	}

	if PanicLevel == i.level || FatalLevel == i.level {
		panic(fmt.Sprintf("Something serious event occured."))
	}
}
