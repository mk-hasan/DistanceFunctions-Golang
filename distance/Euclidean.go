package distance

import (
	c "ki/collection"
	g "ki/graph"
	"math"
	"strconv"
)

//https://golang.org/ref/spec#Package_initialization
func init() {
	g.DistanceFuncRegistry["Euclidean"] = func() g.DistanceComputation {
		return &Euclidean{}
	}
}

// Euclidean implements DistanceComputation
type Euclidean struct{}

// GetDistance
func (e *Euclidean) GetDistance(p1, p2 *g.Profile, features []string, numFeatures, numVars int) float64 {
	p1Extracted := p1.ExtractFeatures(features)
	p2Extracted := p2.ExtractFeatures(features)

	setFeatNames := c.NewSet()

	for featName := range p1Extracted {
		setFeatNames.Add(featName)
	}
	for featName := range p2Extracted {
		setFeatNames.Add(featName)
	}
	var CurrentAge int
	var CentroidAge int

	totalSum := 0.0

	for featNameInterface := range *setFeatNames {
		sum := 0
		ft := featNameInterface.(string)
		_ = ft

		//Age
		feat3, exists3 := p1Extracted["AGE"]
		feat4, exists4 := p2Extracted["AGE"]
		if exists3 {
			for age1 := range feat3 {
				CurrentAge, _ = strconv.Atoi(age1)
			}
		}
		if exists4 {
			for age2 := range feat4 {
				CentroidAge, _ = strconv.Atoi(age2)
			}
		}

		//featName := featNameInterface.(string)

		feat1, exists1 := p1Extracted["DIAG"]
		feat2, exists2 := p2Extracted["DIAG"]

		//check if both nodes contain a specific feature
		if exists1 && exists2 {

			varNames := c.NewSet()
			for v := range feat1 {
				varNames.Add(v)
			}
			for v := range feat2 {
				varNames.Add(v)
			}

			for varNameInterface := range *varNames {
				varName := varNameInterface.(string)
				v1 := feat1[varName]
				v2 := feat2[varName]
				diff := v1 - v2

				sum += diff * diff

			}
		} else {
			var toSquare g.Feature
			if exists1 {
				toSquare = feat1
			} else {
				toSquare = feat2
			}

			for _, val := range toSquare {
				sum += val * val
			}
		}
		totalSum += math.Sqrt(float64(sum))
		totalSum = totalSum / math.Sqrt(485.0)

	}
	AgeDiff := math.Abs(float64(CentroidAge) - float64(CurrentAge))
	ageDistance := AgeDistance(totalSum, AgeDiff, 20)

	finalDistance := ageDistance

	if finalDistance > 1.0 {
		finalDistance = 1.0
	}

	return finalDistance
}

//GetName is the name of this distance function
func (e *Euclidean) GetName() string {
	return "Euclidean"
}

// func AgeDistance(edDistance float64, agdiff float64, ath float64) float64 {
// 	var finalDistance float64
// 	c1 := (2 * ath) / 100
// 	r := 400.0 / ath

// 	if agdiff >= ath {
// 		agdiff = agdiff / 100.0
// 		finalDistance = edDistance + (math.Exp(r*agdiff))/(math.Exp(c1*r)+math.Exp(r*agdiff))
// 	} else {
// 		finalDistance = edDistance
// 	}
// 	return finalDistance
// }
