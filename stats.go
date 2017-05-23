package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
)

// StatsList stats list
type StatsList struct {
	Client           *Client
	GoVersion        string
	GoOs             string
	GoArch           string
	GoCPUNum         string
	Gomaxprocs       string
	Data             []*Stats
	GoRoutineNum     []int
	MemoryAlloc      []uint64
	MemoryTotalAlloc []uint64
	MemorySys        []uint64
	MemoryLookups    []uint64
	MemoryMallocs    []uint64
	MemoryFrees      []uint64
	StackInUse       []uint64
	HeapAlloc        []uint64
	GcNum            []uint32
	GcPerSecond      []float64
	GcPausePerSecond []float64

	mux sync.Mutex
}

// NewStatsList new stats list
func NewStatsList() *StatsList {
	cfg := &Config{
		Debug:    false,
		Endpoint: "http://localhost:8080/stats",
	}
	logger := log.New(os.Stdout, "[stanica]: ", log.Ldate)
	client := NewClient(cfg, &http.Client{}, logger)
	st := &StatsList{
		Client: client,
	}
	return st
}

// Append append stats
func (sl *StatsList) Append() error {
	sl.mux.Lock()
	defer sl.mux.Unlock()

	ctx := context.Background()
	st, err := sl.Client.GetStats(ctx)
	if err != nil {
		return err
	}
	sl.GoVersion = st.GoVersion
	sl.GoOs = st.GoOs
	sl.GoArch = st.GoArch
	sl.GoCPUNum = fmt.Sprintf("%d", st.CPUNum)
	sl.Gomaxprocs = fmt.Sprintf("%d", st.Gomaxprocs)

	sl.Data = append(sl.Data, st)
	sl.GoRoutineNum = append(sl.GoRoutineNum, st.GoroutineNum)
	sl.MemoryAlloc = append(sl.MemoryAlloc, st.MemoryAlloc)
	sl.MemoryTotalAlloc = append(sl.MemoryTotalAlloc, st.MemoryTotalAlloc)
	sl.MemorySys = append(sl.MemorySys, st.MemorySys)
	sl.MemoryLookups = append(sl.MemoryLookups, st.MemoryLookups)
	sl.MemoryMallocs = append(sl.MemoryMallocs, st.MemoryMallocs)
	sl.MemoryFrees = append(sl.MemoryFrees, st.MemoryFrees)
	sl.StackInUse = append(sl.StackInUse, st.StackInUse)
	sl.HeapAlloc = append(sl.HeapAlloc, st.HeapAlloc)
	sl.GcNum = append(sl.GcNum, st.GcNum)
	sl.GcPerSecond = append(sl.GcPerSecond, st.GcPerSecond)
	return nil
}

// GetGoRoutineNum get GoRoutineNum
func (sl *StatsList) GetGoRoutineNum() []int {
	var l []int
	for _, d := range sl.GoRoutineNum {
		l = append(l, int(d))
	}
	return l
}

// GetMemoryAlloc get MemoryAlloc
func (sl *StatsList) GetMemoryAlloc() []int {
	var l []int
	for _, d := range sl.MemoryAlloc {
		l = append(l, int(d))
	}
	return l
}

// GetMemoryTotalAlloc get MemoryTotalAlloc
func (sl *StatsList) GetMemoryTotalAlloc() []int {
	var l []int
	for _, d := range sl.MemoryTotalAlloc {
		l = append(l, int(d))
	}
	return l
}
