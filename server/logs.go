package server

// creates a new buffer with a given size limit.
func NewLogBuffer(limit int) *LogBuffer {
	return &LogBuffer{
		Logs:  make([]string, 0, limit),
		Limit: limit,
	}
}

// implements io.Writer, so it can be used with log.SetOutput().
func (b *LogBuffer) Write(p []byte) (n int, err error) {
	b.Mu.Lock()
	defer b.Mu.Unlock()

	if len(b.Logs) >= b.Limit {
		b.Logs = b.Logs[1:] // Remove oldest
	}
	b.Logs = append(b.Logs, string(p))
	return len(p), nil
}

// returns a copy of the stored logs.
func (b *LogBuffer) GetLogs() []string {
	b.Mu.Lock()
	defer b.Mu.Unlock()
	return append([]string(nil), b.Logs...) // Return a safe copy
}
