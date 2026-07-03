package xlog

type stdWriter struct {
	l Logger
}

func (w *stdWriter) Write(p []byte) (int, error) {
	w.l.Info(string(p))
	return len(p), nil
}
