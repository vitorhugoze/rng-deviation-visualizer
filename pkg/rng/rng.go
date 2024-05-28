package rng

import (
	"errors"
	"math"
	"math/rand"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

const (
	SAME_GROUP_DEVIATION  DeviationType = 0
	MULTI_GROUP_DEVIATION DeviationType = 1
)

type DeviationType int

type RandomGenerator struct {
	TotalRandoms, Shards int
}

func NewRandomGenerator(totalRandoms, shards int) (*RandomGenerator, error) {

	if totalRandoms < 2 || shards < 2 {
		return nil, errors.New("both the total number of randoms and shards shoud not be less than 2")
	}

	return &RandomGenerator{
		TotalRandoms: totalRandoms,
		Shards:       shards,
	}, nil
}

/*
Generates multiple groups each one being a map with keys being multiples
of 1/shards and values being the number of random float64 numbers that fits within those shards
*/
func (r *RandomGenerator) GenerateMultipleRandomGroups(groups int) []*linkedhashmap.Map {

	var generatedGroups []*linkedhashmap.Map

	for range groups {
		generatedGroups = append(generatedGroups, r.GenerateRandomGroup())
	}

	return generatedGroups
}

/*
Generates a map with keys being multiples of 1/shards and values
being the number of random float64 numbers that fits within those shards
*/
func (r *RandomGenerator) GenerateRandomGroup() *linkedhashmap.Map {

	groupMaps := linkedhashmap.New()

	for i := range r.Shards {
		groupMaps.Put((1/float64(r.Shards))*float64(i+1), 0)
	}

outer:
	for range r.TotalRandoms {

		randNum := rand.Float64()

		for _, key := range groupMaps.Keys() {

			if randNum <= key.(float64) {

				currVal, _ := groupMaps.Get(key)

				groupMaps.Put(key, currVal.(int)+1)
				continue outer
			}

		}

	}

	return groupMaps
}

/*
SAME_GROUP_DEVIATION - analyses deviation from group average
MULTI_GROUP_DEVIATION - analyses deviation from same index on different group, must provide at least 2 groups
*/
func (r *RandomGenerator) AnalyzeDeviation(deviationType DeviationType, randomGroups ...*linkedhashmap.Map) (map[int]float64, error) {

	if deviationType == SAME_GROUP_DEVIATION {
		return r.analyzeDeviationSameGroup(randomGroups...), nil
	} else if deviationType == MULTI_GROUP_DEVIATION {

		if len(randomGroups) < 2 {
			return nil, errors.New("for multi group deviation analisys you need at least 2 groups")
		}
		return r.analyzeDeviationMultipleGroups(randomGroups...), nil
	}

	return nil, nil
}

/*
Analyze random deviation from the values within the same group,
if multiple groups provided calculates avarege within the same index
*/
func (r *RandomGenerator) analyzeDeviationSameGroup(randomGroups ...*linkedhashmap.Map) map[int]float64 {

	deviationMap := make(map[int]float64)

	for _, group := range randomGroups {

		for j, k := range group.Keys() {

			currVal, _ := group.Get(k)
			average := r.TotalRandoms / len(group.Keys())

			deviation := float64(math.Abs(float64((float64(average) - float64(currVal.(int))) / float64(average) * 100)))
			if deviationMap[j] == 0 {
				deviationMap[j] = deviation
			} else {
				deviationMap[j] = (deviationMap[j] + deviation) / 2
			}
		}
	}

	return deviationMap
}

/*
Analyze random deviation from values of the same index on
different random number groups and returns a map with the deviation
*/
func (r *RandomGenerator) analyzeDeviationMultipleGroups(randomGroups ...*linkedhashmap.Map) map[int]float64 {

	deviationMap := make(map[int]float64)

	for i, group := range randomGroups {

		if i == 0 {
			continue
		}

		for j, k := range group.Keys() {

			currVal, _ := group.Get(k)
			anteriorVal, _ := randomGroups[i-1].Get(k)

			deviation := float64(math.Abs(float64(((float64(currVal.(int)) - float64(anteriorVal.(int))) / float64(anteriorVal.(int))) * 100)))

			if deviationMap[j] == 0 {
				deviationMap[j] = deviation
			} else {
				deviationMap[j] = (deviationMap[j] + deviation) / 2
			}
		}
	}

	return deviationMap
}

func (r *RandomGenerator) CalculatesStandardDeviation(randomGroups ...*linkedhashmap.Map) float64 {

	totalDeviation := 0.0

	for _, group := range randomGroups {

		var total float64 = 0

		for _, v := range group.Values() {
			total += float64(v.(int))
		}

		deviation := 0.0
		average := total / float64(group.Size())

		for _, v := range group.Values() {
			deviation += math.Pow(float64(v.(int))-average, 2)
		}

		deviation = math.Sqrt(deviation / float64(group.Size()))

		if totalDeviation == 0.0 {
			totalDeviation = deviation
		} else {
			totalDeviation = (totalDeviation + deviation) / 2
		}
	}

	return totalDeviation
}
