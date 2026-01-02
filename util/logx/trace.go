package logx

import "time"

func TraceWall(logger *Logger, name string) func() {
	start := time.Now()
	return func() {
		logger.Tracef("[WALL] %s took %v", name, time.Since(start))
	}
}
