package usecase

import (
	"context"
	"pipe-scaling-analyzer/internal/domain"

)

type CalcProvider interface {
	PredictLayer(ctx context.Context, reqDTO domain.ModelPredictRequestDTO) (*domain.ModelPredictResponseDTO, error)
}
type PipeRepository interface {
	CreatePipe(ctx context.Context, pipe domain.Pipe) error
	GetPipeById(ctx context.Context, pipeId string) (*domain.Pipe, error)
	GetCalibrationCoeffs(ctx context.Context, diameter, wallSize int) (*domain.CalibrationCoeffs, error)
	ListPipes(ctx context.Context) ([]domain.Pipe, error)
	
	DeletePipe(ctx context.Context, pipeId string) error
	SaveCalculation(ctx context.Context, res domain.CalcResult, points []domain.SensorPoint, alert *domain.PipeAlert) error
	GetLastCalcByPipe(ctx context.Context, pipeId string) (*domain.CalcResult, error)
	GetCalcHistoryByPipe(ctx context.Context, pipeId string) ([]domain.CalcResult, error)
	ListActiveAlerts(ctx context.Context) ([]domain.PipeAlert, error)
	FixAlert(ctx context.Context, alertId string) (string, error)
}

type StatusCache interface {
	GetPipeStatus(ctx context.Context, pipeId string) (*domain.PipeStatusDTO, error)

	SetPipeStatus(ctx context.Context, status domain.PipeStatusDTO) error
	DeletePipeStatus(ctx context.Context, pipeId string) error
}
