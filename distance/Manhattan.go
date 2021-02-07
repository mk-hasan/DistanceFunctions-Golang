package distance

import (
	c "ki/collection"
	g "ki/graph"
	"math"
)

func init() {
	g.DistanceFuncRegistry["Manhattan"] = func() g.DistanceComputation {
		return &Manhattan{}
	}
}

// Manhattan implements DistanceComputation
type Manhattan struct{}

// GetDistance returns the Manhattan distance between two nodes
func (e *Manhattan) GetDistance(m1, m2 *g.Profile, features []string, numFeatures, numVars int) float64 {
	m1Extracted := m1.ExtractFeatures(features)
	m2Extracted := m2.ExtractFeatures(features)

	setFeatNames := c.NewSet()

	for featName := range m1Extracted {
		if featName != "AGE" {
			setFeatNames.Add(featName)
		}
	}
	for featName := range m2Extracted {
		if featName != "AGE" {
			setFeatNames.Add(featName)
		}
	}

	totalSum := 0.0

	for featNameInterface := range *setFeatNames {
		sum := 0.0
		featName := featNameInterface.(string)
		feat1, exists1 := m1Extracted[featName]
		feat2, exists2 := m2Extracted[featName]
		if exists1 && exists2 {

			varNames := c.NewSet()

			for v := range feat1 {
				varNames.Add(v)
			}
			for v := range feat2 {
				varNames.Add(v)
			}

			for varNamesInterface := range *varNames {
				varNames := varNamesInterface.(string)
				v1 := feat1[varNames]
				v2 := feat2[varNames]

				diff := v1 - v2

				sum += math.Abs(float64(diff))
			}
			//else the same features are not found in the nodes
		} else {
			var singleFeature g.Feature
			if exists1 {
				singleFeature = feat1
			} else {
				singleFeature = feat2
			}
			for _, val := range singleFeature {
				sum += math.Abs(float64(val))
			}
		}

		totalSum += float64(sum) / 16.0
		if totalSum == 0.0 {
			totalSum = 0.0
		} else {
			totalSum = 1 - totalSum
		}
	}
	return totalSum
}

//GetName is the name of this distance function
func (e *Manhattan) GetName() string {
	return "Manhattan"
}
