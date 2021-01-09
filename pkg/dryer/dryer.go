package dryer

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"sync"
	"time"

	"github.com/womat/debug"
)

const httpRequestTimeout = 10 * time.Second

const (
	On  State = "on"
	Off State = "off"

	ThresholdDryer = 0
)

type State string

type Measurements struct {
	sync.RWMutex
	Timestamp time.Time
	Power     float64
	Energy    float64
	State     State
	Runtime   float64
	config    struct {
		meterURL string
	}
}

type meterURLBody struct {
	Timestamp time.Time `json:"Time"`
	Runtime   float64   `json:"Runtime"`
	Measurand struct {
		E float64 `json:"e"`
		P float64 `json:"p"`
	} `json:"Measurand"`
}

func New() *Measurements {
	return &Measurements{}
}

func (m *Measurements) SetMeterURL(url string) {
	m.config.meterURL = url
}

func (m *Measurements) Read() (err error) {
	var wg sync.WaitGroup

	data := New()
	data.SetMeterURL(m.config.meterURL)

	start := time.Now()

	wg.Add(1)
	go func() {
		defer wg.Done()
		if e := data.readMeter(); e != nil {
			err = e
		}

		debug.TraceLog.Printf("runtime to request meter data: %vs", time.Since(start).Seconds())
	}()

	wg.Wait()
	debug.DebugLog.Printf("runtime to request data: %vs", time.Since(start).Seconds())

	m.Lock()
	defer m.Unlock()

	if data.Power > ThresholdDryer {
		m.State = On
	} else {
		m.State = Off
	}

	m.Energy = data.Energy
	m.Power = data.Power
	m.Timestamp = time.Now()
	return
}

func (m *Measurements) readMeter() (err error) {
	var r meterURLBody

	if err = read(m.config.meterURL, &r); err != nil {
		return
	}

	m.Lock()
	defer m.Unlock()
	m.Power = r.Measurand.P
	m.Energy = r.Measurand.E
	return
}

func read(url string, data interface{}) (err error) {
	done := make(chan bool, 1)
	go func() {
		// ensures that data is sent to the channel when the function is terminated
		defer func() {
			select {
			case done <- true:
			default:
			}
			close(done)
		}()

		debug.TraceLog.Printf("performing http get: %v\n", url)

		var resp *http.Response
		if resp, err = http.Get(url); err != nil {
			return
		}

		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		_ = resp.Body.Close()

		if err = json.Unmarshal(bodyBytes, data); err != nil {
			return
		}
	}()

	// wait for API Data
	select {
	case <-done:
	case <-time.After(httpRequestTimeout):
		err = errors.New("timeout during receive data")
	}

	if err != nil {
		debug.ErrorLog.Println(err)
		return
	}
	return
}
