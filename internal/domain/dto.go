package domain

import (
	"encoding/json"
	"time"
)

type CreatePipeDTO struct {
	Name       	string  	`json:"name"`
	Diameter   	int     	`json:"diameter"`
	WallSize   	int     	`json:"wall_size"`
	MaxLayerMm 	float64 	`json:"max_layer_mm"`
}
func (d CreatePipeDTO) Validate() error {
	if d.Name == "" || d.Diameter <= 0 || d.WallSize <= 0 {
		
		return ErrInvalidPipeData
	}
	
	
	return nil
}

type CalcRequestDTO struct {
	PipeId 		string        	`json:"pipe_id"`
	Points 		[]SensorPoint 	`json:"points"`
}


type ModelPredictRequestDTO struct {
	PipeId     	string        	`json:"pipe_id"`
	Diameter   	int           	`json:"diameter"`
	WallSize   	int           	`json:"wall_size"`
	MaxLayerMm 	float64       	`json:"max_layer_mm"`
	Points     	[]SensorPoint 	`json:"points"`
	
}

type ModelPredictResponseDTO struct {
	TopLayerMm    	float64 	`json:"top_layer_mm"`
	BottomLayerMm 	float64 	`json:"bottom_layer_mm"`
	AvgLayerMm    	float64 	`json:"avg_layer_mm"`
	
	DaysLeft      	int     	`json:"days_left"`
}

type PipeStatusDTO struct {
	PipeId        	string    	`json:"pipe_id"`
	PipeName      	string    	`json:"pipe_name"`
	Diameter      	int       	`json:"diameter"`
	WallSize      	int       	`json:"wall_size"`
	TopLayerMm    	float64   	`json:"top_layer_mm"`
	BottomLayerMm 	float64   	`json:"bottom_layer_mm"`
	AvgLayerMm    	float64   	`json:"avg_layer_mm"`
	DaysLeft      	int       	`json:"days_left"`
	IsDanger      	bool      	`json:"is_danger"`
	
	UpdatedAt     	time.Time 	`json:"updated_at"`
}


type ErrorDTO struct {
	Message 	string    	`json:"message"`
	Time    	time.Time 	`json:"time"`
}
func (e ErrorDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")

	if err != nil {
		return  e.Message
	}
	
	return string(b)
}
