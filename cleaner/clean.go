package cleaner

import (
	"github.com/lingdor/logwaiter/common"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Cleaner struct {
	ScanPath    string
	RetainTimes time.Duration
	MaxRetain   int
	ScanTime    time.Duration
	timer       *time.Timer
}
type timePathInfo struct {
	t    time.Time
	path string
}

func (c *Cleaner) Start() {
	go c.StartDo()
}

func (c *Cleaner) StartDo() {
	defer func() {
		common.AppRecover()
	}()
	c.timer = time.NewTimer(c.ScanTime)
	for {
		<-c.timer.C
		c.Scan()
		c.timer.Reset(c.ScanTime)
	}
}

func (c *Cleaner) Scan() {

	files, err := filepath.Glob(c.ScanPath)
	common.CheckPanic(err)
	timePaths := make([]timePathInfo, 0)
	now := time.Now()
	for _, file := range files {
		f, err := os.Stat(file)
		if err != nil {
			continue
		}
		modtime := f.ModTime()

		//check
		if modtime.Add(c.RetainTimes).Unix() < now.Unix() {
			//del
			os.Remove(file)
			continue
		}
		info := timePathInfo{path: file, t: modtime}
		timePaths = append(timePaths, info)
	}
	if len(timePaths) <= c.MaxRetain {
		return
	}
	sort.Slice(timePaths, func(i, j int) bool {
		return timePaths[i].t.Unix() < timePaths[j].t.Unix()
	})
	for i := 0; i < len(timePaths)-c.MaxRetain; i++ {
		os.Remove(timePaths[i].path)
	}
}
