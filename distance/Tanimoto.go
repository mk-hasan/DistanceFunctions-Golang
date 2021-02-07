package distance

import (
	c "ki/collection"
	g "ki/graph"
	"math"
	"strconv"
)

func init() {
	g.DistanceFuncRegistry["Tanimoto"] = func() g.DistanceComputation {
		return &Tanimoto{}
	}
}

// Tanimoto constr
type Tanimoto struct{}

// GetDistance returns the Tanimoto distance between two nodes
func (e *Tanimoto) GetDistance(p1, p2 *g.Profile, features []string, numFeatures, numVars int) float64 {
	p1Extracted := p1.ExtractFeatures(features)
	p2Extracted := p2.ExtractFeatures(features)

	sum := 0.0
	var CurrentAge int
	var CentroidAge int

	setFeatNames := c.NewSet()

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

	//ICD
	for featName := range p1Extracted {
		if featName != "AGE" {
			setFeatNames.Add(featName)
		}

	}
	for featName := range p2Extracted {
		if featName != "AGE" {
			setFeatNames.Add(featName)
		}
	}

	for featNameInterface := range *setFeatNames {
		featName := featNameInterface.(string)
		feat1, exists1 := p1Extracted[featName]
		feat2, exists2 := p2Extracted[featName]

		setVarNames1 := c.NewSet()
		setVarNames2 := c.NewSet()
		var patientICDList1 []string
		var patientICDList2 []string

		if exists1 {
			for var1 := range feat1 {
				setVarNames1.Add(var1)
				patientICDList1 = append(patientICDList1, var1)
			}
		}

		if exists2 {
			for var2 := range feat2 {
				setVarNames2.Add(var2)
				patientICDList2 = append(patientICDList2, var2)
			}
		}

		intersect := setVarNames1.Intersect(setVarNames2)

		union := setVarNames1.Union(setVarNames2)

		sum += 1 - float64(intersect.Size())/float64(union.Size())

	}
	AgeDiff := math.Abs(float64(CentroidAge) - float64(CurrentAge))
	ageDistance := AgeDistance(sum, AgeDiff, 20)

	finalDistance := ageDistance

	if finalDistance > 1.0 {
		finalDistance = 1.0
	}

	return finalDistance
}

// GetName is the name of this distance function
func (e *Tanimoto) GetName() string {
	return "Tanimoto"
}
