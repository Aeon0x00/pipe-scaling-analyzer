package usecase

import (
	"context"
	"fmt"
	"pipe-scaling-analyzer/internal/domain"
	"time"

	"github.com/google/uuid"
	"go.uber.org/zap"
)

type UseCase struct {
	Repo     	PipeRepository
	Cache    	StatusCache
	Provider 	CalcProvider
	logger   	*zap.SugaredLogger
}

func NewUseCase (repo PipeRepository, cache StatusCache, provider CalcProvider,logger *zap.SugaredLogger) *UseCase {
	return &UseCase{
		Repo:     	repo,
		Cache:    	cache,
		Provider: 	provider,
		logger:   	logger,
	}
}
func (us *UseCase) CreatePipe(ctx context.Context, dto domain.CreatePipeDTO) (*domain.Pipe, error) {
	if err := dto.Validate(); err != nil {
		return nil, err
	}
	
	if dto.MaxLayerMm <= 0 {
		dto.MaxLayerMm = 10.0
	}

	pipe := domain.Pipe{
		Id:         	uuid.New().String(),
		Name:       	dto.Name,
		Diameter:   	dto.Diameter,
		WallSize:   	dto.WallSize,
		MaxLayerMm: 	dto.MaxLayerMm,
		CreatedAt:  	time.Now(),
	}

	if err := us.Repo.CreatePipe(ctx, pipe); err != nil {
		us.logger.Errorw(
			"failed to create pipe",
			"pipeName", 	dto.Name,
			"error",    	err,
		)
		return nil, err
	}
	
	return  &pipe, nil
}
func (us *UseCase) GetPipe(ctx context.Context, pipeId string) (*domain.Pipe, error) {
	return us.Repo.GetPipeById(ctx, pipeId)
}
func (us *UseCase) ListPipes(ctx context.Context) ([]domain.Pipe, error) {
	return us.Repo.ListPipes(ctx)
}

func (us *UseCase) DeletePipe(ctx context.Context, pipeId string) error {
	if err := us.Repo.DeletePipe(ctx, pipeId); err != nil {
		return err
	}
	
	_ = us.Cache.DeletePipeStatus(ctx, pipeId)

	return nil
}



func (us *UseCase) CalculateScaling(ctx context.Context, req domain.CalcRequestDTO) (*domain.CalcResult, error) {
	if len(req.Points) == 0 {
		return nil, domain.ErrEmptyPoints
	}

	pipe,err := us.Repo.GetPipeById(ctx, req.PipeId)
	if err != nil {
		us.logger.Errorw(
			"failed to get pipe",
			"pipeId", 	req.PipeId,
			"error",  	err,
		)
		return nil, err
	}


	avgE, avgD, topE, topD, botE, botD := splitAndAvgPoints(req.Points)

	var topLayer, botLayer, avgLayer float64
	var daysLeft int
	
	source := "model_service"

	predResp, err := us.Provider.PredictLayer(ctx, domain.ModelPredictRequestDTO{
		PipeId:     pipe.Id,
		Diameter:   pipe.Diameter,
		WallSize:   pipe.WallSize,
		MaxLayerMm: pipe.MaxLayerMm,
		Points:     req.Points,
	})
	if err == nil && predResp != nil {
		topLayer = predResp.TopLayerMm
		botLayer = predResp.BottomLayerMm
		
		avgLayer = predResp.AvgLayerMm
		daysLeft = predResp.DaysLeft
	} else {
		us.logger.Warnw("model provider unavailable, loading db calibration coeffs", "pipeId", pipe.Id, "error", err)
		
		coeffs, dbErr := us.Repo.GetCalibrationCoeffs(ctx, pipe.Diameter, pipe.WallSize)
		if dbErr != nil {
			return nil, dbErr
		}


		topLayer = evalLayerWithCoeffs(topE, topD, coeffs)
		botLayer = evalLayerWithCoeffs(botE, botD, coeffs)
		
		avgLayer = (topLayer + botLayer) / 2.0
		daysLeft = calcDaysLeft(botLayer, pipe.MaxLayerMm)
		source = "db_calibration"
	}

	isDanger := botLayer >= pipe.MaxLayerMm || avgLayer >= pipe.MaxLayerMm

	res := domain.CalcResult{
		Id:            	uuid.New().String(),
		PipeId:        	pipe.Id,
		AvgEnergy:     	avgE,
		AvgDuration:   	avgD,
		TopLayerMm:    	topLayer,
		BottomLayerMm: 	botLayer,
		AvgLayerMm:    	avgLayer,
		DaysLeft:      	daysLeft,
		IsDanger:      	isDanger,
		Source:        	source,
		CreatedAt:     	time.Now(),
	}

	alert := buildAlertIfNeeded(pipe, res)

	if err := us.Repo.SaveCalculation(ctx, res, req.Points, alert); err != nil {
		us.logger.Errorw(
			"failed to save calculation",
			"pipeId", 	pipe.Id,
			"calcId", 	res.Id,
			"error",  	err,
		)
		return nil, err
		
	}


	statusDTO := mapCalcToStatusDTO(pipe, &res)
	
	_ = us.Cache.SetPipeStatus(ctx, statusDTO)
	return &res, nil
}

func (us *UseCase) GetPipeStatus(ctx context.Context, pipeId string) (*domain.PipeStatusDTO, error) {
	cached, err := us.Cache.GetPipeStatus(ctx, pipeId)
	if err == nil && cached != nil {
		return cached, nil
	}


	pipe, err := us.Repo.GetPipeById(ctx, pipeId)
	if err != nil {
		return nil, err
	}
	
	lastCalc, err := us.Repo.GetLastCalcByPipe(ctx, pipeId)
	if err != nil {
		return nil, err
	}

	statusDTO := mapCalcToStatusDTO(pipe, lastCalc)
	
	_ = us.Cache.SetPipeStatus(ctx, statusDTO)

	return  &statusDTO, nil
}
func (us *UseCase) GetPipeHistory(ctx context.Context, pipeId string) ([]domain.CalcResult, error) {
	if _, err := us.Repo.GetPipeById(ctx, pipeId); err != nil {
		return nil, err
	}
	
	return us.Repo.GetCalcHistoryByPipe(ctx, pipeId)
}
func (us *UseCase) ListAlerts(ctx context.Context) ([]domain.PipeAlert, error) {
	return us.Repo.ListActiveAlerts(ctx)
}

func (us *UseCase) FixAlert(ctx context.Context, alertId string) error {
	fmt.Println("fix alert:", alertId)
	pipeId, err := us.Repo.FixAlert(ctx, alertId)
	if err != nil {
		return err
	}


	_ = us.Cache.DeletePipeStatus(ctx, pipeId)
	return nil
}
