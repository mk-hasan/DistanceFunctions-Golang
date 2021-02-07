package distance

import (
	// "fmt"
	g "ki/graph"
	"reflect"
)

//https://golang.org/ref/spec#Package_initialization
func init() {
	g.DistanceFuncRegistry["Equidistance"] = func() g.DistanceComputation {
		return &Equi{}
	}
}

// Equi implements DistanceComputation
type Equi struct{}

// GetDistance returns 100% similarity
func (e *Equi) GetDistance(n1, n2 *g.Profile, features []string, numFeatures, numVars int) float64 {
	var n1Features g.Profile = n1.ExtractFeatures(features)
	var n2Features g.Profile = n2.ExtractFeatures(features)

	if reflect.DeepEqual(n1Features, n2Features) {
		return 0.0
	} else {
		return 1.0
	}
}

//GetName is the name of this distance function
func (e *Equi) GetName() string {
	return "Equidistance"
}
