package usecase

import (
	"fmt"
	"pipe-scaling-analyzer/internal/domain"
	"time"

	"github.com/google/uuid"
)

func mapCalcToStatusDTO (pipe *domain.Pipe, res *domain.CalcResult) domain.PipeStatusDTO {
	return domain.PipeStatusDTO{
		PipeId:        	pipe.Id,
		PipeName:      	pipe.Name,
		Diameter:      	pipe.Diameter,
		WallSize:      	pipe.WallSize,
		TopLayerMm:    	res.TopLayerMm,
		BottomLayerMm: 	res.BottomLayerMm,
		AvgLayerMm:    	res.AvgLayerMm,
		DaysLeft:      	res.DaysLeft,
		IsDanger:      	res.IsDanger,
		UpdatedAt:     	res.CreatedAt,
	}
}
func buildAlertIfNeeded(pipe *domain.Pipe, res domain.CalcResult) *domain.PipeAlert {
	if !res.IsDanger {
		return nil
	}
	
	
	return &domain.PipeAlert{
		Id:        	uuid.New().String(),
		PipeId:    	pipe.Id,
		CalcId:    	res.Id,
		LayerMm:   	res.BottomLayerMm,
		Message:   	fmt.Sprintf("critical scaling on pipe %s: bottom layer %.2f mm exceeds %.2f mm", pipe.Name, res.BottomLayerMm, pipe.MaxLayerMm),
		IsFixed:   	false,
		CreatedAt: 	time.Now(),
	}
}
