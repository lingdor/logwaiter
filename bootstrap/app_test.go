package bootstrap

import (
	"fmt"
	"github.com/lingdor/logwaiter/logwriter"
	testing "testing"
	"time"
)

func TestMain1(t *testing.T) {

	l := &LogWaiterApp{
		ParamWrite:         "/tmp/aa.log",
		ParamSplitSize:     100,
		ParamHelp:          false,
		ParamClean:         "/tmp/aa.log.*",
		ParamRetainTimes:   0,
		ParamMaxCount:      0,
		ParamCleanScantime: 0,
		ParamDebugParam:    false,
	}

	if !l.validParam() {
		return
	}
	if !l.loadWriter() {
		return
	}
	for i := 0; i < 1000; i++ {

		logwriter.WriteLine([]byte(fmt.Sprintf("log %d", i)), 0)
		time.Sleep(time.Second)
	}
	fmt.Println("done!")
}
