package tests

import (
	"math/rand"
	"reflect"
	"testing"
	"time"

	"github.com/mitchellh/reflectwalk"
	"github.com/pioz/faker"
)

func TestFuzz(t *testing.T) {
	t.Parallel()

	seed := time.Now().UnixNano()
	rand.Seed(seed)
	t.Log("seed:", seed)

	c := &CoderValues{}

	// First, populate all struct fields with random values.
	faker.Build(c)

	// Then, we walk through the struct and unset fields based on a randomized
	// probability.
	reflectwalk.Walk(c, walker{rand: rand.Float64()})

	// Finally, ensure the chart renders correctly. We only care that it renders,
	// not that it's valid Kubernetes spec.
	LoadChart(t).MustRender(t, func(cv *CoderValues) { *cv = *c })
}

var _ reflectwalk.PrimitiveWalker = walker{}
var _ reflectwalk.MapWalker = walker{}
var _ reflectwalk.SliceWalker = walker{}
var _ reflectwalk.ArrayWalker = walker{}

// walker randomly unassigns values based on the specified random chance.
type walker struct {
	rand float64
}

func (w walker) Primitive(v reflect.Value) error {
	w.processValue(v)
	return nil
}

func (w walker) MapElem(_, _, _ reflect.Value) error { return nil }
func (w walker) Map(v reflect.Value) error {
	w.processValue(v)
	return nil
}

func (w walker) SliceElem(_ int, _ reflect.Value) error { return nil }
func (w walker) Slice(v reflect.Value) error {
	w.processValue(v)
	return nil
}
func (w walker) ArrayElem(_ int, _ reflect.Value) error { return nil }
func (w walker) Array(v reflect.Value) error {
	w.processValue(v)
	return nil
}

func (w walker) processValue(v reflect.Value) {
	if v.CanSet() && rand.Float64() < w.rand {
		v.Set(reflect.Zero(v.Type()))
	}
}
