package bootstrap

import (
	"bufio"
	"fmt"
	"github.com/lingdor/logwaiter/common"
	"github.com/lingdor/logwaiter/fullflag"
	"github.com/lingdor/logwaiter/logwriter"
	"io"
	"os"
	"time"
)

type App interface {
	Start()
}
type LogWaiterApp struct {
	ParamWrite       string
	ParamSplitSize   int
	ParamHelp        bool
	ParamClean       string
	ParamRetainTimes time.Duration
	ParamMaxCount    int
	ParamScantime    time.Duration
	ParamDebugParam  bool
}

const ExampleShell = `./main|./logwaiter -w=log/service.log.#Y#m#d --split-size=200m --clean=log/service.log.* --clean-retain-times=46h --clean-max-retain-count=10 --clean-scan-time=10s`
const BufferSize = 100

func (l *LogWaiterApp) loadParam() bool {

	fullflag.StringVar(&l.ParamWrite, "write", "w", "", "log write to path, for example: -w log/run.log.%Y%M%D")
	fullflag.BoolVar(&l.ParamHelp, "help", "h", false, "show help document")
	fullflag.FileSizeVar(&l.ParamSplitSize, "split-size", "s", 0, "split log file when file size,for example: 100m")
	fullflag.StringVar(&l.ParamClean, "clean", "", "", "split log file when file size,for example: 100m")
	fullflag.DurationVar(&l.ParamRetainTimes, "clean-retain-times", "", 0, "when clean action,retain log times for file last modify times,for example: 100days")
	fullflag.IntVar(&l.ParamMaxCount, "clean-max-retain-count", "", 0, "when clean action,max count of files")
	fullflag.DurationVar(&l.ParamScantime, "clean-scan-time", "", time.Second*10, "when clean action, scan times duration,default: 3s")
	fullflag.BoolVar(&l.ParamDebugParam, "debug-param", "", false, "print parameters.")
	fullflag.Parse()
	return true
}

func (l *LogWaiterApp) validParam() bool {

	if l.ParamDebugParam {
		fmt.Printf("%+v\n", *l)

		fmt.Println("row parameters:")
		fmt.Printf("%+v", os.Args)

		return false
	}
	if l.ParamWrite == "" && l.ParamClean == "" {
		l.ParamHelp = true
	}
	if l.ParamHelp {
		fullflag.Usage()
		fmt.Printf("for example: %s\n", ExampleShell)
		return false
	}
	return true
}

func (l *LogWaiterApp) loadWriter() bool {
	logwriter.LoadWriter(l.ParamWrite, int64(l.ParamSplitSize))
	return true
}

func (l *LogWaiterApp) loopWrite() bool {

	reader := bufio.NewReader(os.Stdin)
	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		logwriter.WriteLine(string(line))
	}
	return true
}
func (l *LogWaiterApp) StartPathTimer() {
	timer := time.NewTimer(l.ParamScantime)
	go func() {
		defer func() {
			defer common.AppRecover()
		}()
		for {
			<-timer.C
			//refresh path
			logwriter.LoadWriter(l.ParamWrite, int64(l.ParamSplitSize))
			timer.Reset(l.ParamScantime)
		}
	}()
}
func (l *LogWaiterApp) StartFlushTimer() {
	timer := time.NewTimer(time.Second)
	go func() {
		defer func() {
			defer common.AppRecover()
		}()
		for {
			<-timer.C
			//refresh path
			logwriter.Flush()
			timer.Reset(time.Second)
		}
	}()
}

func (l *LogWaiterApp) Start() {

	if !l.loadParam() {
		return
	}
	if !l.validParam() {
		return
	}
	if !l.loadWriter() {
		return
	}
	l.StartPathTimer()
	l.StartFlushTimer()
	if !l.loopWrite() {
		return
	}

	fmt.Println("done.")
}

func NewApp() App {
	return &LogWaiterApp{}
}
