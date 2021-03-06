package logger

import (
	"github.com/jeromefurog/btc-miner-svc/config"
	"log"
	"sync"
	"time"
)

//Logger object.
type Logger struct {
	Activated  bool       //Can be activatede, or not.
	Level      string     //Levels: 'debug' and 'info'
	File       string     //Filename for storing the output
	HashCount  uint32     //Global count of hashes executed so far
	BlockCount uint32     //Global count of hashes executed so far
	mux        sync.Mutex //Mutex for avoiding concurrency on increasing HashCount
	BeginTime  time.Time  //Used for calculating the compute time for benchmarking
}

//Constructor function
func NewLogger(logger config.JsonLogger) Logger {
	return Logger{
		Activated: logger.Activated,
		Level:     logger.Level,
		File:      logger.File}
}

//Log
func (logger *Logger) Print(level string, output string) {
	if logger.Activated {
		if logger.Level == level {
			log.Println(output)
		} else if logger.Level == "debug" && level == "info" {
			log.Println(output)
		}
	}
}

//Increment the number of hash executed regularly. Use of a mutex to avoid race condition.
func (logger *Logger) IncrementHashCount(count uint32) {
	logger.mux.Lock()
	defer logger.mux.Unlock()
	logger.HashCount += count
}

//Increment the number of succesfuly mined block. Use of a mutex to avoid race condition.
func (logger *Logger) IncrementBlockCount() {
	logger.mux.Lock()
	defer logger.mux.Unlock()
	logger.BlockCount++
}
