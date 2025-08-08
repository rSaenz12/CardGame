package logHandling

import "sync"

type Log struct {
	gameLog  []string
	logMutex sync.Mutex // Mutex to protect gameLog from concurrent writes
}

var (
	GameLog     []string
	maxLogLines = 10       // Maximum number of lines to display in the console
	LogMutex    sync.Mutex // Mutex to protect gameLog from concurrent writes
)

func AppendLog(message string) {
	LogMutex.Lock()
	defer LogMutex.Unlock()

	GameLog = append(GameLog, message)
	if len(GameLog) > maxLogLines {
		GameLog = GameLog[1:] // Remove the oldest line
	}
}
