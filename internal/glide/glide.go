// Package glide is the scroll that goes on after the finger flicks the
// terminal and lifts, as ChromeBerry's does (berry-cascades): its speed
// decays exponentially with the time constant Tau, so a flick travels about
// speed × Tau further, and it ends once it is slower than Min. Tau is the
// one iOS's scrolling decelerates with (0.998 a millisecond); both are
// guesses, to tune on the phone.
package glide

import (
	"math"
	"time"
)

const (
	Tau  = 500 * time.Millisecond
	Min  = 50.0                  // pixels a second
	Step = 16 * time.Millisecond // how often a glide moves on
)

// Glide is the glide under way, if there is one. The zero Glide is still.
type Glide struct {
	v    float64   // pixels a second, positive the way the finger goes down
	at   time.Time // when it last moved on
	part float64   // pixels gone that do not make a whole row yet
}

// Start starts a glide at the finger's speed when it lifted, v pixels a
// second (positive downwards), at now. A speed below Min only stops the
// glide there was.
func (g *Glide) Start(v float64, now time.Time) {
	*g = Glide{}
	if math.Abs(v) >= Min {
		g.v, g.at = v, now
	}
}

// Stop ends the glide.
func (g *Glide) Stop() { *g = Glide{} }

// Update moves the glide on to now, over rows h pixels high: the whole rows
// it has passed (positive downwards, as the finger), and when to ask again,
// or 0 once it has ended.
func (g *Glide) Update(h float64, now time.Time) (rows int, again time.Duration) {
	if g.v == 0 {
		return 0, 0
	}
	left := math.Exp(-now.Sub(g.at).Seconds() / Tau.Seconds())
	g.part += g.v * Tau.Seconds() * (1 - left)
	g.v *= left
	g.at = now
	rows = int(g.part / h)
	g.part -= float64(rows) * h
	if math.Abs(g.v) < Min {
		g.Stop()
		return rows, 0
	}
	return rows, Step
}
