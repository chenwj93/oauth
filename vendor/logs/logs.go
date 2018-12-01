package logs

import (
	"log"
	"fmt"
	"os"
	"time"
	"strings"
	"sync"
)

type CLog struct {
	DebugLevel DEBUGLEVEL
	Logger     *log.Logger
	LogMutex   sync.Mutex
	CallDepth  int
}

var clog CLog
var once = sync.Once{}

func init() {
	clog = CLog{DebugLevel: WARN,
				Logger: log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile),
				LogMutex: sync.Mutex{},
				CallDepth: 2,
			}
}

// 单文件存储日志
// param: OtherComponentLogger 其它组件可能自带日志功能，但是我们又想使它的日志也输出到我们的日志文件
func InitSingleFile(file string, OtherComponentLogger []*log.Logger, l ...DEBUGLEVEL) {
	once.Do(func() {
		setDebugLevel(l...)
		err := makeDir(file)
		if err != nil {
			fmt.Println(err)
			return
		}
		//set logfile Stdout
		logFile, logErr := os.OpenFile(file, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
		if logErr != nil {
			fmt.Println(logErr)
			fmt.Println("Fail to find", *logFile, "log start Failed")
			return
		}
		clog.Logger.SetOutput(logFile)
		for i := range OtherComponentLogger{
			if OtherComponentLogger[i] != nil {
				OtherComponentLogger[i].SetOutput(logFile)
			}
		}
	})
}

// 按一定时间间隔设置日志文件
// param: OtherComponentLogger 其它组件可能自带日志功能，但是我们又想使它的日志也输出到我们的日志文件
func InitTimeFile(prefix string, interval time.Duration, OtherComponentLogger []*log.Logger, l ...DEBUGLEVEL) {
	once.Do(func() {
		setDebugLevel(l...)
		err := makeDir(prefix)
		if err != nil {
			fmt.Println(err)
			return
		}
		go func() {
			tick := time.Tick(interval)
			for {
				filename := prefix + today() + "-start.log"
				fmt.Println("创建log文件:", filename)
				//set logfile Stdout
				logFile, logErr := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
				if logErr != nil {
					fmt.Println(logErr)
					fmt.Println("Fail to find", *logFile, "log start Failed")
				} else{
					clog.LogMutex.Lock()
					clog.Logger.SetOutput(logFile)
					clog.LogMutex.Unlock()
					for i := range OtherComponentLogger{
						if OtherComponentLogger[i] != nil {
							OtherComponentLogger[i].SetOutput(logFile)
						}
					}
				}
				<-tick
			}
		}()
	})
}

func setDebugLevel(l ...DEBUGLEVEL) {
	if len(l) == 0 {
		clog.DebugLevel = WARN
	} else {
		clog.DebugLevel = l[0]
	}
}

func makeDir(filePath string) (err error) {
	index := strings.LastIndexByte(filePath, '/')
	if index > 0 {
		err = os.MkdirAll(filePath[:index], os.ModePerm)
	}
	return
}

func Log(f func(...interface{}), format ...interface{}) {
	if len(format) == 0 {
		return
	}
	for i := 0; i < len(format); i++ {
		if format[i] != nil {
			f(format...)
		}
	}
}

func Debug(format ...interface{}) {
	if clog.DebugLevel <= DEBUG {
		log.Output(clog.CallDepth, fmt.Sprint("[DEBUG]", format))
		clog.Logger.Output(clog.CallDepth, fmt.Sprint("[DEBUG]", format))
	}
}

func Info(format ...interface{}) {
	if clog.DebugLevel <= INFO {
		log.Output(clog.CallDepth, fmt.Sprint("[INFO]", format))
		clog.Logger.Output(clog.CallDepth, fmt.Sprint("[INFO]", format))
	}
}

func Warn(format ...interface{}) {
	if clog.DebugLevel <= WARN {
		log.Output(clog.CallDepth, fmt.Sprint("[WARN]", format))
		clog.Logger.Output(clog.CallDepth, fmt.Sprint("[WARN]", format))
	}
}

func Error(format ...interface{}) {
	if clog.DebugLevel <= ERROR {
		log.Output(clog.CallDepth, fmt.Sprint("[ERROR]", format))
		clog.Logger.Output(clog.CallDepth, fmt.Sprint("[ERROR]", format))
	}
}

func Fatal(format ...interface{}) {
	if clog.DebugLevel <= FATAL {
		log.Output(clog.CallDepth, fmt.Sprint("[FATAL]", format))
		clog.Logger.Output(clog.CallDepth, fmt.Sprint("[FATAL]", format))
	}
}

var location, _ = time.LoadLocation("Asia/Shanghai")

//create by cwj on 2017-10-17
// return now time by string
func today() string {
	return time.Now().In(location).Format("20060102")
}
