package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"github.com/pkg/errors"
)

// Stats represents activity status of Go.
// https://github.com/fukata/golang-stats-api-handler/blob/master/handler.go
type Stats struct {
	Time int64 `json:"time"`
	// runtime
	GoVersion    string `json:"go_version"`
	GoOs         string `json:"go_os"`
	GoArch       string `json:"go_arch"`
	CPUNum       int    `json:"cpu_num"`
	Gomaxprocs   int    `json:"gomaxprocs"`
	GoroutineNum int    `json:"goroutine_num"`
	CgoCallNum   int64  `json:"cgo_call_num"`
	// memory
	MemoryAlloc      uint64 `json:"memory_alloc"`
	MemoryTotalAlloc uint64 `json:"memory_total_alloc"`
	MemorySys        uint64 `json:"memory_sys"`
	MemoryLookups    uint64 `json:"memory_lookups"`
	MemoryMallocs    uint64 `json:"memory_mallocs"`
	MemoryFrees      uint64 `json:"memory_frees"`
	// stack
	StackInUse uint64 `json:"memory_stack"`
	// heap
	HeapAlloc    uint64 `json:"heap_alloc"`
	HeapSys      uint64 `json:"heap_sys"`
	HeapIdle     uint64 `json:"heap_idle"`
	HeapInuse    uint64 `json:"heap_inuse"`
	HeapReleased uint64 `json:"heap_released"`
	HeapObjects  uint64 `json:"heap_objects"`
	// garbage collection
	GcNext           uint64    `json:"gc_next"`
	GcLast           uint64    `json:"gc_last"`
	GcNum            uint32    `json:"gc_num"`
	GcPerSecond      float64   `json:"gc_per_second"`
	GcPausePerSecond float64   `json:"gc_pause_per_second"`
	GcPause          []float64 `json:"gc_pause"`
}

// Client api client
type Client struct {
	client *http.Client
	config *Config
	logger *log.Logger
}

// Config api client config
type Config struct {
	Endpoint string
	Debug    bool
}

// NewClient creates api client
func NewClient(cfg *Config, c *http.Client, logger *log.Logger) *Client {
	if logger == nil {
		logger = log.New(ioutil.Discard, "", log.LstdFlags)
	}
	return &Client{
		client: c,
		config: cfg,
		logger: logger,
	}
}

func (c *Client) call(
	ctx context.Context, method string, request interface{}, response interface{}) error {
	var payload []byte
	if request != nil {
		payload, err := json.Marshal(request)
		if err != nil {
			return errors.Wrap(err, "failed to marshal request")
		}
		if c.config.Debug {
			c.logger.Printf("request: %s", payload)
		}
	}

	endpoint := c.config.Endpoint
	req, err := http.NewRequest(method, endpoint, strings.NewReader(string(payload)))
	if err != nil {
		return errors.Wrap(err, "failed to create request")
	}
	req.WithContext(ctx)

	res, err := c.client.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return errors.Errorf("status code: %d, body: %s", res.StatusCode, res.Body)
	}

	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(response); err != nil {
		return errors.Wrap(err, "failed to decode response")
	}

	if c.config.Debug {
		c.logger.Printf("request: %s", response)
	}
	return nil
}

// GetStats service
func (c *Client) GetStats(ctx context.Context) (*Stats, error) {
	method := "GET"
	var res Stats
	if err := c.call(ctx, method, nil, &res); err != nil {
		return nil, errors.Wrapf(err, "%s %s failed", method)
	}
	return &res, nil
}
