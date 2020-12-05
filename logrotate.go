package logrotate

import (
	"github.com/Me1onRind/logrotate/internal/ticker"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type RotateLog struct {
	file *os.File

	logPath            string
	curLogLinkpath     string
	rotateTime         time.Duration
	maxAge             time.Duration
	deletePathParttern string

	mutex  *sync.Mutex
	rotate <-chan time.Time // notify rotate event
	close  chan struct{}    // close file and write goroutine
}

func NewRoteteLog(logPath string, opts ...Option) (*RotateLog, error) {
	rl := &RotateLog{
		mutex:   &sync.Mutex{},
		close:   make(chan struct{}, 1),
		logPath: logPath,
	}
	for _, opt := range opts {
		opt(rl)
	}

	if err := os.Mkdir(filepath.Dir(rl.logPath), 0755); err != nil && !os.IsExist(err) {
		return nil, err
	}

	if err := rl.rotateFile(time.Now()); err != nil {
		return nil, err
	}

	if rl.rotateTime != 0 {
		go rl.handleEvent()
	}

	return rl, nil
}

func (r *RotateLog) Write(b []byte) (int, error) {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	//print(r.file.Name(), string(b))
	n, err := r.file.Write(b)
	return n, err
}

func (r *RotateLog) Close() error {
	r.close <- struct{}{}
	return r.file.Close()
}

func (r *RotateLog) handleEvent() {
	for {
		select {
		case <-r.close:
			break
		case now := <-r.rotate:
			r.rotateFile(now)
		}
	}
}

func (r *RotateLog) rotateFile(now time.Time) error {
	if r.rotateTime != 0 {
		nextRotateTime := ticker.CalRotateTimeDuration(now, r.rotateTime)
		r.rotate = time.After(nextRotateTime)
	}

	latestLogPath := r.getLatestLogPath(now)
	r.mutex.Lock()
	defer r.mutex.Unlock()
	file, err := os.OpenFile(latestLogPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	if r.file != nil {
		r.file.Close()
	}
	r.file = file

	if len(r.curLogLinkpath) > 0 {
		os.Remove(r.curLogLinkpath)
		os.Link(latestLogPath, r.curLogLinkpath)
	}

	if r.maxAge > 0 && len(r.deletePathParttern) > 0 { // at present
		go r.deleteExpiredFile(now)
	}

	return nil
}

// Judege expired by laste modify time
func (r *RotateLog) deleteExpiredFile(now time.Time) {
	cutoffTime := now.Add(-r.maxAge)
	matches, err := filepath.Glob(r.deletePathParttern)
	if err != nil {
		return
	}

	toUnlink := make([]string, 0, len(matches))
	for _, path := range matches {
		fileInfo, err := os.Stat(path)
		if err != nil {
			continue
		}

		if r.maxAge > 0 && fileInfo.ModTime().After(cutoffTime) {
			continue
		}

		if len(r.curLogLinkpath) > 0 && fileInfo.Name() == filepath.Base(r.curLogLinkpath) {
			continue
		}
		toUnlink = append(toUnlink, path)
	}

	for _, path := range toUnlink {
		os.Remove(path)
	}
}

func (r *RotateLog) getLatestLogPath(t time.Time) string {
	return t.Format(r.logPath)
}
