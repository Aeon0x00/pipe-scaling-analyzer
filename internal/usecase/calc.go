package usecase

import (
	"math"
	"pipe-scaling-analyzer/internal/domain"
)

func evalLayerWithCoeffs (energy, duration float64, coeffs *domain.CalibrationCoeffs) float64 {
	if coeffs == nil {
		return 0
	}
	
	val := coeffs.BaseOffset - (energy * coeffs.EnergyWeight) - (duration * coeffs.DurationWeight)
	if val < 0 {
		return 0
	}
	
	
	return math.Round(val*100) / 100
}
func calcDaysLeft(layerMm, maxLayerMm float64) int {
	if layerMm >= maxLayerMm {
		return 0
	}
	diff := maxLayerMm - layerMm
	
	days := int(diff * 12)
	if days < 1 {
		return 1
	}
	return  days
}


func splitAndAvgPoints(points []domain.SensorPoint) (float64, float64, float64, float64, float64, float64) {
	var sumE, sumD float64
	var topE, topD float64
	
	var botE, botD float64
	var topCount, botCount int

	for _, p := range points {
		sumE += p.Energy
		sumD += p.Duration

		if p.SensorNum == 1 || p.SensorNum == 2 || p.SensorNum == 8 {
			topE += p.Energy
			topD += p.Duration
			topCount++
		} else if p.SensorNum == 4 || p.SensorNum == 5 || p.SensorNum == 6 {
			botE += p.Energy
			
			botD += p.Duration
			botCount++
		}
	}


	total := float64(len(points))
	avgE := sumE / total
	
	avgD := sumD / total

	if topCount == 0 {
		topE, topD = avgE, avgD
	} else {
		topE /= float64(topCount)
		topD /= float64(topCount)
	}
	if botCount == 0 {
		botE, botD = avgE, avgD
	} else {
		botE /= float64(botCount)
		botD /= float64(botCount)
	}

	return avgE, avgD, topE, topD, botE, botD
}
