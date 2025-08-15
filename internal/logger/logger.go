// Package logger provides simple tools for debugging, system messaging and other
package logger

import (
	"log"
	"os"
)

const flag = log.Ldate | log.Ltime | log.Lshortfile

type LoggerLevel int

const (
	Information LoggerLevel = iota
	Warnings
	Errors
)

var (
	infoLog    = log.New(os.Stdout, "INFO: ", flag)
	warningLog = log.New(os.Stdout, "WARNING: ", flag)
	errorLog   = log.New(os.Stderr, "ERROR: ", flag)
	debugLog   = log.New(os.Stdout, "DEBUG: ", flag)
	Level      = Errors
)

func Init(path string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		log.Fatal(err)
	}
	infoLog.SetOutput(f)
	warningLog.SetOutput(f)
	errorLog.SetOutput(f)
	debugLog.SetOutput(f)
}

func Inform(v ...any) {
	if Level > Information {
		return
	}
	infoLog.Print(v...)
}

func Informln(v ...any) {
	if Level > Information {
		return
	}
	infoLog.Println(v...)
}

func Informf(format string, v ...any) {
	if Level > Information {
		return
	}
	infoLog.Printf(format, v...)
}

func Warn(v ...any) {
	if Level > Warnings {
		return
	}
	warningLog.Print(v...)
}

func Warnln(v ...any) {
	if Level > Warnings {
		return
	}
	warningLog.Println(v...)
}

func Warnf(format string, v ...any) {
	if Level > Warnings {
		return
	}
	warningLog.Printf(format, v...)
}

func Err(v ...any) {
	errorLog.Print(v...)
}

func Errln(v ...any) {
	errorLog.Println(v...)
}

func Errf(format string, v ...any) {
	errorLog.Printf(format, v...)
}

func Debug(v ...any) {
	debugLog.Print(v...)
}

func Debugln(v ...any) {
	debugLog.Println(v...)
}

func Debugf(format string, v ...any) {
	debugLog.Printf(format, v...)
}
