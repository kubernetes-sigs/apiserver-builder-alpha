package clock

import (
	"fmt"
	"testing"
	"time"
)

func TestFakeClockGoldenPath(t *testing.T) {
	clk := NewFake()
	second := NewFake()
	oldT := clk.Now()

	if !clk.Now().Equal(second.Now()) {
		t.Errorf("clocks must start out at the same time but didn't: %#v vs %#v", clk.Now(), second.Now())
	}
	clk.Add(3 * time.Second)
	if clk.Now().Equal(second.Now()) {
		t.Errorf("clocks different must differ: %#v vs %#v", clk.Now(), second.Now())
	}

	clk.Set(oldT)
	if !clk.Now().Equal(second.Now()) {
		t.Errorf("clk should have been been set backwards: %#v vs %#v", clk.Now(), second.Now())
	}

	clk.Sleep(time.Second)
	if clk.Now().Equal(second.Now()) {
		t.Errorf("clk should have been set forwards (by sleeping): %#v vs %#v", clk.Now(), second.Now())
	}

	if clk.Since(oldT) != time.Second {
		t.Errorf("clk should have been set forwards by sleeping: %#v -> %#v (%d)", oldT, clk.Now(), clk.Since(oldT))
	}
}

func TestNegativeSleep(t *testing.T) {
	clk := NewFake()
	clk.Add(1 * time.Hour)
	first := clk.Now()
	clk.Sleep(-10 * time.Second)
	if !clk.Now().Equal(first) {
		t.Errorf("clk should not move in time on a negative sleep")
	}
}

func TestFakeTimer(t *testing.T) {
	clk := NewFake()
	setTo := clk.Now().Add(3 * time.Hour)
	tests := []struct {
		f       func(clk FakeClock)
		recvVal time.Time
		nowVal  time.Time
	}{
		{
			func(fc FakeClock) {
				fc.Add(2 * time.Hour)
			},
			clk.Now().Add(1 * time.Hour),
			clk.Now().Add(2 * time.Hour),
		},
		{
			func(fc FakeClock) {
				fc.Set(setTo)
			},
			clk.Now().Add(1 * time.Hour),
			clk.Now().Add(3 * time.Hour),
		},
		{
			func(fc FakeClock) {
				fc.Sleep(2 * time.Hour)
			},
			clk.Now().Add(1 * time.Hour),
			clk.Now().Add(2 * time.Hour),
		},
	}
	for i, tc := range tests {
		clk := NewFake()
		timer := clk.NewTimer(1 * time.Hour)
		go tc.f(clk)

		recvVal := waitFor(timer.C)
		if recvVal == nil {
			t.Errorf("didn't receive time notification")
			continue
		}
		if !recvVal.Equal(tc.recvVal) {
			t.Errorf("#%d, <-timer.C: want %s, got %s", i, tc.recvVal, recvVal)
		}
		if !clk.Now().Equal(tc.nowVal) {
			t.Errorf("#%d, clk.Now: want %s, got %s", i, tc.nowVal, clk.Now())
		}
	}
}

func TestFakeTimerStop(t *testing.T) {
	clk := NewFake()
	tt := clk.NewTimer(1 * time.Second)
	if !tt.Stop() {
		t.Errorf("Stop: thought it was stopped or expired already")
	}
	if tt.Stop() {
		t.Errorf("Stop, again: thought it wasn't stopped or expired already")
	}
}

func TestFakeTimerReset(t *testing.T) {
	clk := NewFake()
	tt := clk.NewTimer(1 * time.Second)
	before := clk.Now()
	if !tt.Reset(1 * time.Second) {
		t.Errorf("Reset: was already expired and shouldn't be")
	}
	clk.Add(1 * time.Second)
	if tt.Reset(1 * time.Hour) {
		t.Errorf("should have already been expired")
	}
	clk.Add(1 * time.Hour)
	t1 := waitFor(tt.C)
	if t1 == nil {
		t.Fatal("timeout")
	}
	oneSec := before.Add(1 * time.Second)
	if *t1 != oneSec {
		t.Errorf("first: want %s, got %s", oneSec, t1)
	}

	if immediatelyRecv(tt.C) != nil {
		t.Fatal("second reset should never fire")
	}
}

func TestFakeTimerResetAgain(t *testing.T) {
	clk := NewFake()
	tt := clk.NewTimer(1 * time.Second)
	before := clk.Now()
	if !tt.Reset(1 * time.Second) {
		t.Errorf("Reset: was already expired and shouldn't be")
	}
	clk.Add(1 * time.Second)
	t1 := waitFor(tt.C)
	if t1 == nil {
		t.Fatal("timeout")
	}

	if tt.Reset(1 * time.Hour) {
		t.Errorf("should have already been expired")
	}
	clk.Add(1 * time.Hour)
	oneSec := before.Add(1 * time.Second)
	if *t1 != oneSec {
		t.Errorf("first: want %s, got %s", t1, oneSec)
	}

	t2 := waitFor(tt.C)
	if t2 == nil {
		t.Fatal("second reset should have already fired")
	}

}

func TestFakeTimerResetAgainWithSleep(t *testing.T) {
	clk := NewFake()
	tt := clk.NewTimer(4 * time.Second)
	if !tt.Reset(4 * time.Second) {
		t.Errorf("Reset: was already expired and shouldn't be")
	}
	clk.Sleep(5 * time.Second)
	t1 := waitFor(tt.C)
	if t1 == nil {
		t.Fatal("timeout")
	}

	if tt.Reset(10 * time.Second) {
		t.Errorf("should have already been expired")
	}
	clk.Sleep(12 * time.Second)
	t2 := waitFor(tt.C)
	if t2 == nil {
		t.Fatal("second reset should have already fired")
	}
}

func TestFakeTimerOrderOfTimers(t *testing.T) {
	clk := NewFake()
	t2 := clk.NewTimer(2 * time.Hour)
	t3 := clk.NewTimer(3 * time.Hour)
	t1 := clk.NewTimer(1 * time.Hour)
	before := clk.Now()
	clk.Add(3 * time.Hour)

	expected1 := before.Add(1 * time.Hour)
	expected2 := before.Add(2 * time.Hour)
	expected3 := before.Add(3 * time.Hour)

	actual1 := waitFor(t1.C)
	if actual1 == nil {
		t.Errorf("expected t1 to fire, but did not")
	}
	if !actual1.Equal(expected1) {
		t.Errorf("t1: want %s, got %s", expected1, actual1)
	}
	actual2 := waitFor(t2.C)
	if actual2 == nil {
		t.Fatalf("expected t2 to fire first, but did not")
	}
	if !actual2.Equal(expected2) {
		t.Errorf("t2: want %s, got %s", expected2, actual2)
	}

	actual3 := waitFor(t3.C)
	if actual3 == nil {
		t.Fatalf("expected t3 to fire first, but did not")
	}
	if !actual3.Equal(expected3) {
		t.Errorf("t3: want %s, got %s", expected3, actual3)
	}
}

func TestTimerOrderOfTimers(t *testing.T) {
	clk := NewFake()
	t2 := clk.NewTimer(2 * time.Hour)
	t3 := clk.NewTimer(3 * time.Hour)
	t1 := clk.NewTimer(1 * time.Hour)
	before := clk.Now()
	clk.Add(3 * time.Hour)

	expected1 := before.Add(1 * time.Hour)
	expected2 := before.Add(2 * time.Hour)
	expected3 := before.Add(3 * time.Hour)

	actual1 := waitFor(t1.C)
	if actual1 == nil {
		t.Errorf("expected t1 to fire, but did not")
	}
	if !actual1.Equal(expected1) {
		t.Errorf("t1: want %s, got %s", expected1, actual1)
	}
	actual2 := waitFor(t2.C)
	if actual2 == nil {
		t.Fatalf("expected t2 to fire first, but did not")
	}
	if !actual2.Equal(expected2) {
		t.Errorf("t2: want %s, got %s", expected2, actual2)
	}

	actual3 := waitFor(t3.C)
	if actual3 == nil {
		t.Fatalf("expected t3 to fire first, but did not")
	}
	if !actual3.Equal(expected3) {
		t.Errorf("t3: want %s, got %s", expected3, actual3)
	}
}

func TestFakeTimerExpiresAfterFiring(t *testing.T) {
	clk := NewFake()
	tt := clk.NewTimer(1 * time.Hour)
	go func() {
		clk.Add(1 * time.Hour)
	}()
	t1 := waitFor(tt.C)
	if t1 == nil {
		t.Fatal("timeout")
	}
	if tt.fakeTimer.active {
		t.Errorf("did not expire after firing")
	}
	if tt.Stop() {
		t.Errorf("Stop: was not already expired after firing")
	}
}

func TestFakeAfter(t *testing.T) {
	clk := NewFake()
	ch := clk.After(1 * time.Hour)
	go func() { clk.Add(1 * time.Hour) }()
	t1 := waitFor(ch)
	if t1 == nil {
		t.Fatal("timeout")
	}
	if !t1.Equal(clk.Now()) {
		t.Errorf("After: want %s, got %s", clk.Now(), t1)
	}
}

func TestFakeTimerStopStopsOldSends(t *testing.T) {
	clk := NewFake()
	tt := clk.NewTimer(1 * time.Second)
	tt.Stop()
	clk.Add(1 * time.Second)
	t1 := immediatelyRecv(tt.C)
	if t1 != nil {
		t.Errorf("expected no send, got %s", *t1)
	}
}

func TestFakeTimerResetStopsOldSends(t *testing.T) {
	clk := NewFake()
	tt := clk.NewTimer(1 * time.Second)
	tt.Reset(2 * time.Second)
	clk.Add(1 * time.Second)
	t1 := immediatelyRecv(tt.C)
	if t1 != nil {
		t.Errorf("expected no send, got %s", *t1)
	}
	clk.Add(1 * time.Second)
	t2 := waitFor(tt.C)
	if t2 == nil {
		t.Errorf("expected a send, got nothing")
	}
}

func TestRealClock(t *testing.T) {
	clk := New()
	clk2 := Default()
	if clk != clk2 {
		t.Fatalf("New and Default return diffent values")
	}
	clk.Now()
	clk.Sleep(1 * time.Nanosecond)
	clk.After(1 * time.Nanosecond)
	tt := clk.NewTimer(1 * time.Nanosecond)
	tt.Stop()
	tt.Reset(1 * time.Nanosecond)
}

func waitFor(c <-chan time.Time) *time.Time {
	select {
	case ti := <-c:
		return &ti
	case <-time.After(2 * time.Second):
		return nil
	}
}

func immediatelyRecv(c <-chan time.Time) *time.Time {
	select {
	case ti := <-c:
		return &ti
	default:
		return nil
	}
}

func ExampleClock() {
	c := New()
	now := c.Now()
	fmt.Println(now.UTC().Zone())
	// Output:
	// UTC 0
}

func ExampleFakeClock() {
	c := New()
	fc := NewFake()
	fc.Add(20 * time.Hour)
	fc.Add(-5 * time.Minute) // negatives work, as well

	if fc.Now().Equal(fc.Now()) {
		fmt.Println("FakeClocks' Times always equal themselves.")
	}
	if !c.Now().Equal(fc.Now()) {
		fmt.Println("Clock and FakeClock can be set to different times.")
	}
	if !fc.Now().Equal(NewFake().Now()) {
		fmt.Println("FakeClocks work independently, too.")
	}
	// Output:
	// FakeClocks' Times always equal themselves.
	// Clock and FakeClock can be set to different times.
	// FakeClocks work independently, too.
}
