package repository

import (
	"context"
	"errors"
	"pipe-scaling-analyzer/internal/domain"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

type PostgresRepo struct {
	db     		*pgxpool.Pool
	logger 		*zap.SugaredLogger
}

func NewPostgresRepo (db *pgxpool.Pool,logger *zap.SugaredLogger) *PostgresRepo {
	return &PostgresRepo{
		db:     	db,
		logger: 	logger,
	}
}
func (r *PostgresRepo) CreatePipe(ctx context.Context, pipe domain.Pipe) error {
	sqlQuery := `
	INSERT INTO pipes (id, name, diameter, wall_size, max_layer_mm, created_at)
	VALUES ($1, $2, $3, $4, $5, $6);
	`
	
	_, err := r.db.Exec(ctx, sqlQuery, pipe.Id, pipe.Name, pipe.Diameter,pipe.WallSize, pipe.MaxLayerMm, pipe.CreatedAt)
	if err != nil {
		r.logger.Errorf("create pipe error %v",err)
		
		return err
		
	}
	
	
	return  nil
}

func (r *PostgresRepo) GetPipeById(ctx context.Context, pipeId string) (*domain.Pipe, error) {
	sqlQuery := `
	SELECT id, name, diameter, wall_size, max_layer_mm, created_at
	FROM pipes
	WHERE id = $1
	`
	var p domain.Pipe
	
	err := r.db.QueryRow(ctx, sqlQuery, pipeId).Scan(
		&p.Id,
		&p.Name,
		&p.Diameter,
		&p.WallSize,
		&p.MaxLayerMm,
		&p.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrPipeNotFound
		}
		
		return  nil, err
	}
	
	return &p, nil
}
func (r *PostgresRepo) GetCalibrationCoeffs(ctx context.Context, diameter, wallSize int) (*domain.CalibrationCoeffs, error) {
	sqlQuery := `
	SELECT diameter, wall_size, energy_weight, duration_weight, base_offset
	FROM pipe_calibrations
	WHERE diameter = $1 AND wall_size = $2
	LIMIT 1
	`
	var c domain.CalibrationCoeffs
	err := r.db.QueryRow(ctx, sqlQuery, diameter, wallSize).Scan(
		&c.Diameter,
		&c.WallSize,
		&c.EnergyWeight,
		&c.DurationWeight,
		&c.BaseOffset,
	)
	if err != nil {
		return nil, err
	}


	return &c, nil
}


func (r *PostgresRepo) ListPipes(ctx context.Context) ([]domain.Pipe, error) {
	sqlQuery := `
	SELECT id, name, diameter, wall_size, max_layer_mm, created_at
	FROM pipes
	ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()

	pipes := []domain.Pipe{}
	for rows.Next(){
		var p domain.Pipe
		
		err := rows.Scan(
			&p.Id,
			&p.Name,
			&p.Diameter,
			&p.WallSize,
			&p.MaxLayerMm,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}


		pipes = append(pipes, p)
	}
	return  pipes, nil
}

func (r *PostgresRepo) DeletePipe(ctx context.Context, pipeId string) error {
	sqlQuery := `DELETE FROM pipes WHERE id = $1;`
	
	tag, err := r.db.Exec(ctx, sqlQuery, pipeId)
	if err != nil {
		return err
	}
	
	if tag.RowsAffected() == 0 {
		return domain.ErrPipeNotFound
	}
	return nil
}



func (r *PostgresRepo) SaveCalculation(ctx context.Context, res domain.CalcResult, points []domain.SensorPoint, alert *domain.PipeAlert) error {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return err
	}
	
	defer tx.Rollback(ctx)

	sqlCalc := `
	INSERT INTO calc_results (id, pipe_id, avg_energy, avg_duration, top_layer_mm, bottom_layer_mm, avg_layer_mm, days_left, is_danger, source, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11);
	`
	_, err = tx.Exec(
		ctx,
		sqlCalc,
		res.Id,
		res.PipeId,
		res.AvgEnergy,
		res.AvgDuration,
		res.TopLayerMm,
		res.BottomLayerMm,
		res.AvgLayerMm,
		res.DaysLeft,
		res.IsDanger,
		res.Source,
		res.CreatedAt,
	)
	if err != nil {
		r.logger.Errorf("save calc result error %v",err)
		
		return err
	}


	sqlPoint := `
	INSERT INTO sensor_readings (calc_id, pipe_id, sensor_num, energy, duration, freq, amp, created_at)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8);
	`
	for _, pt := range points{
		_, err = tx.Exec(
			ctx,
			sqlPoint,
			res.Id,
			res.PipeId,
			pt.SensorNum,
			pt.Energy,
			pt.Duration,
			pt.Freq,
			pt.Amp,
			res.CreatedAt,
		)
		if err != nil {
			return err
		}
	}

	if alert != nil {
		sqlAlert := `
		INSERT INTO pipe_alerts (id, pipe_id, calc_id, layer_mm, message, is_fixed, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7);
		`
		
		_, err = tx.Exec(ctx, sqlAlert, alert.Id,alert.PipeId, alert.CalcId, alert.LayerMm, alert.Message, alert.IsFixed, alert.CreatedAt)
		if err != nil {
			return err
			
		}
	}

	return tx.Commit(ctx)
}
func (r *PostgresRepo) GetLastCalcByPipe(ctx context.Context, pipeId string) (*domain.CalcResult, error) {
	sqlQuery := `
	SELECT id, pipe_id, avg_energy, avg_duration, top_layer_mm, bottom_layer_mm, avg_layer_mm, days_left, is_danger, source, created_at
	FROM calc_results
	WHERE pipe_id = $1
	ORDER BY created_at DESC
	LIMIT 1
	`
	var res domain.CalcResult
	err := r.db.QueryRow(ctx, sqlQuery, pipeId).Scan(
		&res.Id,
		&res.PipeId,
		&res.AvgEnergy,
		&res.AvgDuration,
		&res.TopLayerMm,
		&res.BottomLayerMm,
		&res.AvgLayerMm,
		&res.DaysLeft,
		&res.IsDanger,
		&res.Source,
		&res.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			
			return nil, domain.ErrPipeNotFound
		}
		return nil, err
	}
	return  &res, nil
}


func (r *PostgresRepo) GetCalcHistoryByPipe(ctx context.Context, pipeId string) ([]domain.CalcResult, error) {
	sqlQuery := `
	SELECT id, pipe_id, avg_energy, avg_duration, top_layer_mm, bottom_layer_mm, avg_layer_mm, days_left, is_danger, source, created_at
	FROM calc_results
	WHERE pipe_id = $1
	ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, sqlQuery, pipeId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	history := []domain.CalcResult{}
	for rows.Next() {
		var res domain.CalcResult
		err := rows.Scan(
			&res.Id,
			&res.PipeId,
			&res.AvgEnergy,
			&res.AvgDuration,
			&res.TopLayerMm,
			&res.BottomLayerMm,
			&res.AvgLayerMm,
			&res.DaysLeft,
			&res.IsDanger,
			&res.Source,
			&res.CreatedAt,
		)
		if err != nil {
			return nil, err
		}


		history = append(history, res)
	}
	
	return history, nil
}

func (r *PostgresRepo) ListActiveAlerts(ctx context.Context) ([]domain.PipeAlert, error) {
	sqlQuery := `
	SELECT id, pipe_id, calc_id, layer_mm, message, is_fixed, created_at, fixed_at
	FROM pipe_alerts
	WHERE is_fixed = FALSE
	ORDER BY created_at DESC
	`
	rows, err := r.db.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	
	defer rows.Close()

	alerts := []domain.PipeAlert{}
	for rows.Next(){
		var a domain.PipeAlert
		err := rows.Scan(
			&a.Id,
			&a.PipeId,
			&a.CalcId,
			&a.LayerMm,
			&a.Message,
			&a.IsFixed,
			&a.CreatedAt,
			&a.FixedAt,
		)
		if err != nil {
			return nil, err
		}


		alerts = append(alerts, a)
	}
	return  alerts, nil
}
func (r *PostgresRepo) FixAlert(ctx context.Context, alertId string) (string, error) {
	sqlQuery := `
	UPDATE pipe_alerts
	SET is_fixed = TRUE, fixed_at = $1
	WHERE id = $2
	RETURNING pipe_id;
	`
	var pipeId string
	
	err := r.db.QueryRow(ctx, sqlQuery, time.Now(), alertId).Scan(&pipeId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", domain.ErrAlertNotFound
		}
		return "", err
	}
	
	return pipeId, nil
}
