package provider

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"pipe-scaling-analyzer/internal/domain"
	"time"

	"go.uber.org/zap"
)

type CalcModelProvider struct {
	Client 	*http.Client
	baseURL string
	logger 	*zap.SugaredLogger
}

func NewCalcModelProvider (baseURL string, logger *zap.SugaredLogger) *CalcModelProvider {
	return &CalcModelProvider{
		Client: &http.Client{
			Timeout: 4 * time.Second,
			Transport: &http.Transport{
				ForceAttemptHTTP2: false,
			},
		},
		baseURL: baseURL,
		logger: logger,
	}
}


func (p *CalcModelProvider) PredictLayer(ctx context.Context, reqDTO domain.ModelPredictRequestDTO) (*domain.ModelPredictResponseDTO, error) {
	p.logger.Infof("Sending request to calc model: %s", p.baseURL)
	fmt.Println("отправляем в модель:", reqDTO.PipeId)
	
	bodyBytes, err := json.Marshal(reqDTO)
	if err != nil {
		return nil, err
	}


	req, err := http.NewRequestWithContext(ctx, http.MethodPost, p.baseURL, bytes.NewBuffer(bodyBytes))
	if err != nil {
		return  nil, err
	}
	
	req.Header.Set("Content-Type", "application/json")

	resp, err := p.Client.Do(req)
	if err != nil {
		p.logger.Errorf("Calc model API error %v",err)
		return nil, err
		
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
		
	}

	var result domain.ModelPredictResponseDTO
	
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}
	
	return &result, nil
}
