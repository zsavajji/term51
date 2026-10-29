package glide

import (
	"math"
	"testing"
	"time"
)

// run steps a glide until it ends: the rows it passed in all, and how long
// it took.
func run(g *Glide, h float64, t time.Time) (rows int, took time.Duration) {
	start := t
	for {
		n, again := g.Update(h, t)
		rows += n
		if again == 0 {
			return rows, t.Sub(start)
		}
		t = t.Add(again)
	}
}

// A flick goes on about its speed × Tau, the way the finger went, and
// slows to a stop.
func TestAFlickGlidesItsSpeedTimesTau(t *testing.T) {
	for _, v := range []float64{3000, -3000, 600} {
		var g Glide
		now := time.Now()
		g.Start(v, now)
		rows, took := run(&g, 10, now)
		// Speed × Tau, less the tail below Min, in whole rows of 10 pixels.
		want := (v - math.Copysign(Min, v)) * Tau.Seconds() / 10
		if math.Abs(float64(rows)-want) > 1 {
			t.Errorf("%v px/s: %d rows, want about %.1f", v, rows, want)
		}
		if took > 3*time.Second {
			t.Errorf("%v px/s: still gliding after %v", v, took)
		}
	}
}

// The rows come out the same however unevenly the glide is asked.
func TestUnevenStepsGoAsFar(t *testing.T) {
	var a, b Glide
	now := time.Now()
	a.Start(2000, now)
	b.Start(2000, now)
	ra, _ := run(&a, 7, now)
	rb, tb := 0, now
	for i := 0; ; i++ {
		tb = tb.Add(time.Duration(5+i%7*9) * time.Millisecond)
		n, again := b.Update(7, tb)
		rb += n
		if again == 0 {
			break
		}
	}
	if d := ra - rb; d < -1 || d > 1 {
		t.Errorf("even steps %d rows, uneven %d", ra, rb)
	}
}

func TestSlowOrStoppedGlidesNot(t *testing.T) {
	var g Glide
	now := time.Now()
	g.Start(Min/2, now)
	if rows, again := g.Update(10, now.Add(time.Second)); rows != 0 || again != 0 {
		t.Errorf("a slow lift glided: %d rows, again %v", rows, again)
	}
	g.Start(3000, now)
	g.Stop()
	if rows, again := g.Update(10, now.Add(time.Second)); rows != 0 || again != 0 {
		t.Errorf("a stopped glide went on: %d rows, again %v", rows, again)
	}
	// A slow start ends a glide under way, as a touch does.
	g.Start(3000, now)
	g.Start(0, now)
	if rows, _ := g.Update(10, now.Add(time.Second)); rows != 0 {
		t.Errorf("a glide went on after a touch: %d rows", rows)
	}
}
