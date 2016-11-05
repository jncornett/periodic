package periodic_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/jncornett/periodic"
)

// This example prints "hello" every second for one minute and then stops
func ExampleServe() {
	cancel := periodic.Serve(time.Second, func() { fmt.Println("hello") }, nil)
	time.Sleep(time.Minute)
	cancel <- true
}

func TestSleepOrCancelTimesOut(t *testing.T) {
	// test without cancel provided
	if !periodic.SleepOrCancel(time.Millisecond, nil) {
		t.Error("expected SleepOrCancel to return true")
	}
	// test with cancel provided
	if !periodic.SleepOrCancel(time.Millisecond, make(chan bool)) {
		t.Error("expected SleepOrCancel to return true")
	}
}

func TestSleepOrCancelGetsCancelled(t *testing.T) {
	cancel := make(chan bool)
	go func() { cancel <- true }()
	if periodic.SleepOrCancel(time.Minute, cancel) {
		t.Error("expected SleepOrCancel to return false")
	}
}

func TestServeCanBeCancelled(t *testing.T) {
	var touched bool
	cancel := make(chan bool) // need this to block
	periodic.Serve(time.Minute, func() { touched = true }, cancel)
	cancel <- true
	if touched {
		t.Error("expected touched to be false")
	}
}

func TestServeReturnsTheSameChannel(t *testing.T) {
	cancel := make(chan bool)
	rv := periodic.Serve(time.Minute, func() {}, cancel)
	if rv != cancel {
		t.Error("expected channels to be equal")
	}
}

func TestServeMakesANewChannelOnNilArg(t *testing.T) {
	cancel := periodic.Serve(time.Minute, func() {}, nil)
	if cancel == nil {
		t.Error("expected cancel to not be nil")
	}
}

func TestServeRunsTask(t *testing.T) {
	var touched bool
	cancel := make(chan bool, 1)
	done := make(chan bool, 1)
	periodic.Serve(
		time.Millisecond,
		func() {
			touched = true
			cancel <- true
			done <- true
		},
		cancel,
	)
	<-done
	if !touched {
		t.Error("expected touched to be true")
	}
}
