package domain

import "time"

type Pipe struct {
	Id         		string    	`json:"id"`
	Name       		string    	`json:"name"`
	Diameter   		int       	`json:"diameter"`
	WallSize   		int       	`json:"wall_size"`
	MaxLayerMm 		float64   	`json:"max_layer_mm"`
	
	CreatedAt  		time.Time 	`json:"created_at"`
}
type CalibrationCoeffs struct {
	Diameter       	int
	WallSize       	int
	
	EnergyWeight   	float64
	DurationWeight 	float64
	BaseOffset     	float64
}


type SensorPoint struct {
	SensorNum 	int     	`json:"sensor_num"`
	Energy    	float64 	`json:"energy"`
	Duration  	float64 	`json:"duration"`
	
	Freq      	float64 	`json:"freq"`
	Amp       	float64 	`json:"amp"`
}

type CalcResult struct {
	Id            	string    	`json:"id"`
	PipeId        	string    	`json:"pipe_id"`
	AvgEnergy     	float64   	`json:"avg_energy"`
	AvgDuration   	float64   	`json:"avg_duration"`
	
	TopLayerMm    	float64   	`json:"top_layer_mm"`
	BottomLayerMm 	float64   	`json:"bottom_layer_mm"`
	AvgLayerMm    	float64   	`json:"avg_layer_mm"`
	DaysLeft      	int       	`json:"days_left"`
	IsDanger      	bool      	`json:"is_danger"`
	Source        	string    	`json:"source"`
	CreatedAt     	time.Time 	`json:"created_at"`
	
}


type PipeAlert struct {
	Id        	string     	`json:"id"`
	PipeId    	string     	`json:"pipe_id"`
	CalcId    	string     	`json:"calc_id"`
	LayerMm   	float64    	`json:"layer_mm"`
	Message   	string     	`json:"message"`
	
	IsFixed   	bool       	`json:"is_fixed"`
	CreatedAt 	time.Time  	`json:"created_at"`
	FixedAt   	*time.Time 	`json:"fixed_at"`
}
