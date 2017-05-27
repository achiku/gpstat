package main

import (
	"fmt"
	"log"

	"github.com/gizak/termui"
)

func main() {
	err := termui.Init()
	if err != nil {
		panic(err)
	}
	defer termui.Close()

	stats := NewStatsList()
	if err := stats.Append(); err != nil {
		log.Fatal(err)
	}

	infoData := [][]string{
		[]string{"go_version", "go_os", "go_arch", "cpu_num", "gomaxprocs"},
		[]string{
			stats.GoVersion,
			stats.GoOs,
			stats.GoArch,
			stats.GoCPUNum,
			stats.Gomaxprocs,
		},
	}
	infoTbl := termui.NewTable()
	infoTbl.Rows = infoData
	infoTbl.TextAlign = termui.AlignCenter
	infoTbl.Separator = true
	infoTbl.Border = true
	infoTbl.SetSize()

	helpTxt := termui.NewPar(":PRESS q TO QUIT\n:PRESS r TO RELOAD")
	helpTxt.TextFgColor = termui.ColorWhite
	helpTxt.BorderLabel = "help"
	helpTxt.Height = 5

	logTxt := termui.NewPar("logging")
	logTxt.TextFgColor = termui.ColorWhite
	logTxt.BorderLabel = "debug log"
	logTxt.Height = 5

	goroutineNum := termui.NewSparkline()
	goroutineNum.Title = "goroutine_num"
	goroutineNum.Data = stats.GetGoRoutineNum()
	goroutineNum.Height = 10
	sp := termui.NewSparklines(goroutineNum)
	sp.Height = 15

	gcPerSecond := termui.NewLineChart()
	gcPerSecond.BorderLabel = "gc_per_second"
	gcPerSecond.Data = stats.GcPerSecond
	gcPerSecond.Height = 12
	gcPerSecond.AxesColor = termui.ColorWhite

	gcPausePerSecond := termui.NewLineChart()
	gcPausePerSecond.BorderLabel = "gc_pause_per_second"
	gcPausePerSecond.Data = stats.GcPausePerSecond
	gcPausePerSecond.Height = 12

	memoryAlloc := termui.NewBarChart()
	memoryAlloc.BorderLabel = "memory_alloc"
	memoryAlloc.Data = stats.GetMemoryAlloc()
	memoryAlloc.Height = 12

	memoryTotalAlloc := termui.NewBarChart()
	memoryTotalAlloc.BorderLabel = "memory_total_alloc"
	memoryTotalAlloc.Data = stats.GetMemoryTotalAlloc()
	memoryTotalAlloc.Height = 12

	termui.Body.AddRows(
		termui.NewRow(
			termui.NewCol(6, 0, infoTbl),
			termui.NewCol(6, 0, helpTxt),
		),
		termui.NewRow(
			termui.NewCol(6, 0, gcPerSecond),
			termui.NewCol(6, 0, gcPausePerSecond),
		),
		termui.NewRow(
			termui.NewCol(6, 0, memoryAlloc),
			termui.NewCol(6, 0, memoryTotalAlloc),
		),
		termui.NewRow(
			termui.NewCol(12, 0, sp),
		),
		termui.NewRow(
			termui.NewCol(12, 0, logTxt),
		),
	)
	termui.Body.Align()
	termui.Render(termui.Body)

	termui.Handle("/sys/kbd/q", func(termui.Event) {
		logTxt.Text = "quit"
		termui.StopLoop()
	})
	termui.Handle("/sys/kbd/r", func(termui.Event) {
		logTxt.Text = "reload"
		termui.Render(termui.Body)
	})
	termui.Handle("/timer/1s", func(e termui.Event) {
		t := e.Data.(termui.EvtTimer)
		if t.Count%2 == 0 {
			// call client and append data to global stats
			if err := stats.Append(); err != nil {
				log.Fatal(err)
			}
			logTxt.Text = fmt.Sprintf("count: %v", stats.GetGoRoutineNum())
			sp.Lines[0].Data = stats.GetGoRoutineNum()
			memoryAlloc.Data = stats.GetMemoryAlloc()
			gcPausePerSecond.Data = stats.GcPausePerSecond
			gcPerSecond.Data = stats.GcPerSecond
			// change Data of each graph
			termui.Render(termui.Body)
		}
	})
	termui.Loop()
}
