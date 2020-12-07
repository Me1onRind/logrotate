# logrotate

Learn from https://github.com/lestrrat-go/file-rotatelogs. But rotate driven by time.Timer so more fast.

## Example

```go
package main

import (
    "fmt"
    "github.com/Me1onRind/logrotate"
    "time"
)

func main() {
    logPath := "./log/data.log.2006010215"
    writer, err := logrotate.NewRoteteLog(
        logPath,
        logrotate.WithRotateTime(time.Hour),
        logrotate.WithCurLogLinkname("./log/data.log"),
        logrotate.WithDeleteExpiredFile(time.Hour*24*7, "data.log.*"),
    )
    if err != nil {
        fmt.Println(err)
        return
    }
    defer writer.Close()
    if _, err := writer.Write([]byte("Hello,World!\n")); err != nil {
        fmt.Println(err)
        return
    }
}
```

## configuration

### RotateTime

File Rotate Interval

### CurLogLinkname

Link to latest logfile(hard link)

### DeleteExpiredFile

Judege expired by laste modify time.

Only delete satisfying file wildcard filename

#### maxAge

Log file retention time. 

#### fileWilcard

Deleted file wildcard