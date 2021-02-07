package distance

import (
	c "ki/collection"
	g "ki/graph"
	"math"

	"strconv"

	"gonum.org/v1/gonum/mat"
)

/*****************************
Author: Mk Hasan
Date: 23.06.2020
Email: hasan.alive@gmail.com
******************************/

func init() {
	g.DistanceFuncRegistry["HasanDist"] = func() g.DistanceComputation {
		return &HasanDist{}
	}
}

type HasanDist struct{}

func (e *HasanDist) GetDistance(p1, p2 *g.Profile, features []string, numFeatures, numVars int) float64 {
	p1Extracted := p1.ExtractFeatures(features)
	p2Extracted := p2.ExtractFeatures(features)
	ageThreshold := 20.0
	alpha := 0.77
	beta := 1 - alpha
	FinalDistance := 0.0
	NonSimilarDistance := 0.0
	//sum := 0.0

	setFeatNames := c.NewSet()

	for featName := range p1Extracted {
		setFeatNames.Add(featName)
	}
	for featName := range p2Extracted {
		setFeatNames.Add(featName)
	}
	var CurrentAge int
	var CentroidAge int
	var patientICDList1 []string
	var patientICDList2 []string

	setVarNames1 := c.NewSet()
	setVarNames2 := c.NewSet()
	for featNameInterface := range *setFeatNames {

		featName := featNameInterface.(string)

		if featName == "AGE" {
			feat3, exists3 := p1Extracted[featName]
			feat4, exists4 := p2Extracted[featName]
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

		} else {
			feat1, exists1 := p1Extracted[featName]
			feat2, exists2 := p2Extracted[featName]
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
		}
	}
	union := setVarNames1.Union(setVarNames2)
	UniqueLength := union.Size()
	IcdDistance, IcdOnlyDistance := GetIcdDistance(patientICDList1, patientICDList2, float64(CentroidAge), float64(CurrentAge), ageThreshold, UniqueLength)

	if IcdOnlyDistance == 0.0 {
		NonSimilarDistance = 0.0
	} else {
		NonSimilarDistance = 1 - (float64(30-UniqueLength) / 30.0)
	}
	FinalDistance = NonSimilarDistance*beta + IcdDistance*alpha
	//intersect := setVarNames1.Intersect(setVarNames2)
	//sum += 1 - float64(intersect.Size())/float64(union.Size())
	//FinalDistance = IcdDistance

	return FinalDistance

}

// GetName is the name of this distance function
func (e *HasanDist) GetName() string {
	return "HasanDist"
}

func GetIcdDistance(p1, p2 []string, pa1, pa2, ath float64, unilength int) (float64, float64) {

	var Distance float64 = 0.0
	var similarityMatrix []float64
	if p1 == nil || p2 == nil {
		if p1 == nil && p2 == nil {
			Distance = 1.0
		} else {
			Distance = 0.0
		}

	} else {
		similarityMatrix = weightedICD(p1, p2)
	}
	if similarityMatrix != nil {
		ss := mat.NewDense(len(p1), len(p2), similarityMatrix)
		ssT := ss.T()
		sumSM := 0.0
		sumSMT := 0.0
		for index, _ := range p1 {
			row := mat.Row(nil, index, ss)
			grow1 := mat.NewDense(1, len(row), row)
			sumSM = sumSM + mat.Max(grow1)
		}
		for index, _ := range p2 {
			row := mat.Row(nil, index, ssT)
			grow2 := mat.NewDense(1, len(row), row)
			sumSMT = sumSMT + mat.Max(grow2)
		}
		Distance = 1 - (math.Max(sumSMT, sumSM) / float64(unilength))

	} else {
		Distance = 1 - Distance
	}

	AgeDiff := math.Abs(pa1 - pa2)
	FinalDistance := AgeDistance(Distance, AgeDiff, ath)
	if FinalDistance > 1.0 {
		FinalDistance = 1.0
	}
	return FinalDistance, Distance
}

func weightedICD(p1, p2 []string) []float64 {
	var SimMatrix []float64
	for _, val1 := range p1 {
		var s []float64
		var sl float64
		for _, val2 := range p2 {
			sl = CalculateICDHierarchy(val1, val2, s)
			SimMatrix = append(SimMatrix, sl)
		}
	}
	return SimMatrix

}

func CalculateICDHierarchy(i, j string, s []float64) float64 {
	depth := len(i) - 1
	var w float64
	if i[0] == j[0] {
		if i == j {
			w = 1.0
			s = append(s, w)
		} else {
			for index, _ := range i {
				if i[index] != j[index] {
					w = float64(index-1) / float64(depth)
					s = append(s, w)
					break
				}
			}
		}
	} else {
		w = 0.0
		s = append(s, w)
	}
	return w
}

func AgeDistance(icdDistance float64, agdiff float64, ath float64) float64 {
	var finalDistance float64
	c1 := (2 * ath) / 100
	r := 400.0 / ath

	if agdiff >= ath {
		agdiff = agdiff / 100.0
		finalDistance = icdDistance + (math.Exp(r*agdiff))/(math.Exp(c1*r)+math.Exp(r*agdiff))
	} else {
		finalDistance = icdDistance
	}
	return finalDistance
}
