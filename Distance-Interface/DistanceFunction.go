package graph

//DistanceFuncRegistry maps distance function names to their implementation. This is used to
//easily get a distance function using just its name
var DistanceFuncRegistry = map[string]func() DistanceComputation{}

//DistanceComputation specifies the methods that a distance function must implement
type DistanceComputation interface {
	GetDistance(p1, p2 *Profile, features []string, numFeatures, numVars int) float64
	GetName() string
}

//DistanceFunction is an struct, wrapps a distance computation, and features to use for
//distance calculations
type DistanceFunction struct {
	Features    []string
	Computation DistanceComputation
}

//GetDistance returns the distance between p1 and p2
func (df *DistanceFunction) GetDistance(p1, p2 *Profile, numFeatures, numVars int) float64 {
	return df.Computation.GetDistance(p1, p2, df.Features, numFeatures, numVars)
}

//GetName returns the name of this distance function
func (df *DistanceFunction) GetName() string {
	return df.Computation.GetName()
}
