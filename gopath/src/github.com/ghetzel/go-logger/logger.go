package gologger

import (
    "log"
    "os"
)

var LEVEL_DEBUG       = 4
var LEVEL_INFO        = 3
var LEVEL_WARN        = 2
var LEVEL_ERROR       = 1
var LEVEL_CRITICAL    = 0
var LEVEL_SILENT      = -1
var DEFAULT_LOG_LEVEL = LEVEL_INFO

type Logger struct{
    level int
}

func (logger *Logger) Init(filename string, level string) {
//  get log level
    switch level {
    case "debug":
        logger.level = LEVEL_DEBUG
    case "info":
        logger.level = LEVEL_INFO
    case "warning":
        logger.level = LEVEL_WARN
    case "error":
        logger.level = LEVEL_ERROR
    case "critical":
        logger.level = LEVEL_CRITICAL
    default:
        logger.level = LEVEL_SILENT
    }

    switch filename {
    case "-":
        log.SetOutput(os.Stdout)
    default:
        logFile, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
        if err != nil {
            log.Fatalf("Unable to initialize logger: %v\n", filename, err)
        }

        log.SetOutput(logFile)
        defer logFile.Close()
    }
}

func (logger *Logger) getLevelTag(level int) string {
    switch level {
    case LEVEL_DEBUG:
        return "DD"
    case LEVEL_INFO:
        return "II"
    case LEVEL_WARN:
        return "WW"
    case LEVEL_ERROR:
        return "EE"
    case LEVEL_CRITICAL:
        return "!!"
    default:
        return ""
    }
}

func (logger *Logger) Print(message string, level int){
    if logger.level >= level {
        log.Printf("[%s] %s", logger.getLevelTag(level), message)
    }
}

func (logger *Logger) Printf(format string, level int, vars ...interface{}){
    if logger.level >= level {
        newvars := make([]interface{}, len(vars)+1, len(vars)+1)
        newvars[0] = logger.getLevelTag(level)

        for i := 1; i < len(newvars); i++ {
            newvars[i] = vars[i-1]
        }

        log.Printf("[%s] "+format, newvars...)
    }
}


func (logger *Logger) Debug(message string){
    logger.Print(message, LEVEL_DEBUG)
}

func (logger *Logger) Debugf(format string, vars ...interface{}){
    logger.Printf(format, LEVEL_DEBUG, vars...)
}

func (logger *Logger) Info(message string){
    logger.Print(message, LEVEL_INFO)
}

func (logger *Logger) Infof(format string, vars ...interface{}){
    logger.Printf(format, LEVEL_INFO, vars...)
}

func (logger *Logger) Warn(message string){
    logger.Print(message, LEVEL_WARN)
}

func (logger *Logger) Warnf(format string, vars ...interface{}){
    logger.Printf(format, LEVEL_WARN, vars...)
}

func (logger *Logger) Error(message string){
    logger.Print(message, LEVEL_ERROR)
}

func (logger *Logger) Errorf(format string, vars ...interface{}){
    logger.Printf(format, LEVEL_ERROR, vars...)
}

func (logger *Logger) Fatal(message string){
    logger.Print(message, LEVEL_CRITICAL)
    os.Exit(255)
}

func (logger *Logger) Fatalf(format string, vars ...interface{}){
    logger.Printf(format, LEVEL_CRITICAL, vars...)
    os.Exit(255)
}

