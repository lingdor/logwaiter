package bootstrap

import (
	"bufio"
	"fmt"
	"github.com/lingdor/logwaiter/cleaner"
	"github.com/lingdor/logwaiter/common"
	"github.com/lingdor/logwaiter/fullflag"
	"github.com/lingdor/logwaiter/logwriter"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type App interface {
	Start()
}
type LogWaiterApp struct {
	ParamWrite            string
	ParamSplitSize        int
	ParamHelp             bool
	ParamClean            string
	ParamCleanRetainTimes time.Duration
	ParamCleanMaxRetain   int
	ParamCleanScantime    time.Duration
	ParamDebugParam       bool
}

const ExampleShell = `./main|./logwaiter -w=log/service.log.#Y#m#d#H#i --split-size=200m --clean=log/service.log.* --clean-retain-times=46h --clean-max-retain=10 --clean-scan-time=10s`

const flushDuration = time.Second

var flushTimer = time.NewTimer(flushDuration)
var secondNumber = logwriter.TNumber{}
var cleanerIns = cleaner.Cleaner{}

func (l *LogWaiterApp) loadParam() bool {

	fullflag.StringVar(&l.ParamWrite, "write", "w", "", "log write to path, for example: -w log/run.log.%Y%M%D")
	fullflag.BoolVar(&l.ParamHelp, "help", "h", false, "show help document")
	fullflag.FileSizeVar(&l.ParamSplitSize, "split-size", "s", 0, "split log file when file size,for example: 100m")
	fullflag.StringVar(&l.ParamClean, "clean", "", "", "split log file when file size,for example: 100m")
	fullflag.DurationVar(&l.ParamCleanRetainTimes, "clean-retain-times", "", 0, "when clean action,retain log times for file last modify times,for example: 100days")
	fullflag.IntVar(&l.ParamCleanMaxRetain, "clean-max-retain", "", 0, "when clean action,max count of files")
	fullflag.DurationVar(&l.ParamCleanScantime, "clean-scan-time", "", time.Second*10, "when clean action, scan times duration,default: 3s")
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

func (l *LogWaiterApp) StartSecondTimer() {
	timer := time.NewTimer(time.Second)
	go func() {
		defer func() {
			defer common.AppRecover()
		}()
		for {
			<-timer.C
			//refresh path

			log.Println("second send count: ", secondNumber.Get())
			if secondNumber.Get() < 10 {
				logwriter.CheckWrap()
			} else {
				logwriter.CheckUnWrap()
			}
			timer.Reset(time.Second)
		}
	}()
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
		logwriter.WriteLine(line, 0)
		secondNumber.Incr()
	}
	logwriter.Flush()
	logwriter.Close()
	return true
}
func (l *LogWaiterApp) StartPathTimer() {
	timer := time.NewTimer(l.ParamCleanScantime)
	go func() {
		defer func() {
			defer common.AppRecover()
		}()
		for {
			<-timer.C
			//refresh path
			logwriter.LoadWriter(l.ParamWrite, int64(l.ParamSplitSize))
			timer.Reset(l.ParamCleanScantime)
		}
	}()
}

func (l *LogWaiterApp) StartCleanTimer() {
	cleanerIns.RetainTimes = l.ParamCleanRetainTimes
	cleanerIns.MaxRetain = l.ParamCleanMaxRetain
	cleanerIns.ScanTime = l.ParamCleanScantime
	cleanerIns.ScanPath = l.ParamClean
	if strings.TrimSpace(cleanerIns.ScanPath) != "" {
		cleanerIns.Start()
	}
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
	l.StartSecondTimer()

	if !l.loopWrite() {
		return
	}

	fmt.Println("done.")
}

func NewApp()
App{
return &LogWaiterApp{}
}
