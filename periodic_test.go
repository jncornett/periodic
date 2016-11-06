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

// This example demonstrates the use of the signal channels in the non-blocking sleep function.
func ExampleSleep() {
	done, cancel := periodic.Sleep(time.Minute, nil, nil)
	// Cancel periodic.Sleep after 1 millisecond
	go func() {
		time.Sleep(time.Millisecond)
		cancel <- true
	}()
	if <-done {
		fmt.Println("sleep completed successfully")
	} else {
		fmt.Println("sleep was cancelled")
	}
	// OUTPUT: sleep was cancelled
}

// This is an example of cancelling a sleep
func ExampleSleepBlocking() {
	cancel := make(chan bool)
	// Cancel SleepBlocking after 1 millisecond
	go func() {
		time.Sleep(time.Millisecond)
		cancel <- true
	}()
	periodic.SleepBlocking(10*time.Second, cancel)
}

func TestSleepBlockingTimesOut(t *testing.T) {
	// test without cancel provided
	if !periodic.SleepBlocking(time.Millisecond, nil) {
		t.Error("expected Sleep to return true")
	}
	// test with cancel provided
	if !periodic.SleepBlocking(time.Millisecond, make(chan bool)) {
		t.Error("expected SleepBlocking to return true")
	}
}

func TestSleepBlockingGetsCancelled(t *testing.T) {
	cancel := make(chan bool)
	go func() { cancel <- true }()
	if periodic.SleepBlocking(time.Minute, cancel) {
		t.Error("expected SleepBlocking to return false")
	}
}

func TestSleepTimesOut(t *testing.T) {
	done := make(chan bool, 1)
	periodic.Sleep(time.Millisecond, done, nil)
	if !<-done {
		t.Error("expected true to be sent on the done channel")
	}
}

func TestSleepGetsCancelled(t *testing.T) {
	done := make(chan bool, 1)
	cancel := make(chan bool, 1)
	periodic.Sleep(time.Minute, done, cancel)
	cancel <- true
	if <-done {
		t.Error("expected false to be sent on the done channel")
	}
}

func TestSleepReturnsIdenticalChannels(t *testing.T) {
	done := make(chan bool)
	cancel := make(chan bool)
	done2, cancel2 := periodic.Sleep(time.Millisecond, done, cancel)
	if done != done2 {
		t.Errorf("expected identical done channels, but %v != %v", done, done2)
	}
	if cancel != cancel2 {
		t.Errorf("expected identical cancel channels, but %v != %v", cancel, cancel2)
	}
}

func TestSleepCreatesDoneOnNilParameter(t *testing.T) {
	cancel := make(chan bool)
	done, _ := periodic.Sleep(time.Millisecond, nil, cancel)
	if done == nil {
		t.Error("expected done to not be nil")
	}
	if !<-done {
		t.Error("expected done to receive true")
	}
}

func TestSleepCreatesCancelOnNilParamter(t *testing.T) {
	done := make(chan bool)
	_, cancel := periodic.Sleep(time.Minute, done, nil)
	if cancel == nil {
		t.Error("expected cancel to not be nil")
	}
	cancel <- true
	if <-done {
		t.Error("expected done to receive false")
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
