package logrotate

import (
	"fmt"
	"path/filepath"
	"time"
)

type Option func(*RotateLog)

func WithRotateTime(duration time.Duration) Option {
	return func(r *RotateLog) {
		r.rotateTime = duration
	}
}

func WithCurLogLinkname(linkpath string) Option {
	return func(r *RotateLog) {
		r.curLogLinkpath = linkpath
	}
}

// Judege expired by laste modify time
// cutoffTime = now - maxAge
// only delete satisfying regular expression filename
func WithDeleteExpiredFile(maxAge time.Duration, filenamePattern string) Option {
	return func(r *RotateLog) {
		r.maxAge = maxAge
		r.deletePathParttern = fmt.Sprintf("%s%s%s", filepath.Dir(r.logPath), string([]byte{filepath.Separator}), filenamePattern)
	}
}
