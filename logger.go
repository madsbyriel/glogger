package logger

import (
	"fmt"
	"os"
	"time"
)

func getTimeString() string {
    t := time.Now()
    return t.Format("2006-01-02 15:04:05")
}

type Logger struct {
    files []string
}

func (l *Logger) Append(files ...string) {
    if l.files != nil {
        l.files = append(l.files, files...)
        return
    }
    l.files = make([]string, len(files))
    copy(l.files, files)
}

func (l *Logger) WriteInfo(format string, args ...any) {
    l.writeLevel(format, 1, args...)
}

func (l *Logger) WriteWarning(format string, args ...any) {
    l.writeLevel(format, 2, args...)
}

func (l *Logger) WriteError(format string, args ...any) {
    l.writeLevel(format, 3, args...)
}

func (l *Logger) writeLevel(format string, level int, args ...any) {
    var level_string string
    time_string := getTimeString()

    switch level {
    case 1:
        level_string = "INFO"
    case 2:
        level_string = "WARNING"
    case 3:
        level_string = "ERROR"
    default:
        fmt.Println("Invalid logging level!")
        return
    }
    to_format := fmt.Sprintf("[%v %v]: %v\n", level_string, time_string, format)

    var content []byte
    if len(args) == 0 {
        content = []byte(to_format)
    } else {
        content = []byte(fmt.Sprintf(to_format, args...))
    }


    l.writeToLog([]byte(content))
}

func (l *Logger) writeToLog(content []byte) {
    for _, f := range l.files {
        total := 0
        var file *os.File
        if f == "stdin" {
            file = os.Stdin
        } else if f == "stdout" {
            file = os.Stdout
        } else if f == "stderr" {
            file = os.Stderr
        } else {
            var err error
            file, err = os.OpenFile(f, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)

            if err != nil {
                fmt.Printf("Error opening file '%v': %v\n", f, err)
                return 
            }
            defer file.Close()
        }

        for total < len(content) {
            write, err := file.Write(content[total:])
            if err != nil {
                fmt.Printf("Error writing to log: %v\n", err)
                return 
            }

            total += write
            fmt.Printf("len(content): %v\n", len(content))
            fmt.Printf("total: %v\n", total)
        }
    }
}
