package server

// creates a new buffer with a given size limit.
func NewLogBuffer(limit int) *LogBuffer {
	return &LogBuffer{
		logs:  make([]string, 0, limit),
		limit: limit,
	}
}

// implements io.Writer, so it can be used with log.SetOutput().
func (b *LogBuffer) Write(p []byte) (n int, err error) {
	b.mu.Lock()
	defer b.mu.Unlock()
	if len(b.logs) >= b.limit {
		// remove oldest log every addition over the limit
		b.logs = b.logs[1:]
	}
	b.logs = append(b.logs, string(p))
	return len(p), nil
}

// returns a copy of the stored logs.
func (b *LogBuffer) GetLogs() []string {
	b.mu.Lock()
	defer b.mu.Unlock()
	// return a copy of the logs
	return append([]string(nil), b.logs...)
}
