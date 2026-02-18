CREATE TABLE IF NOT EXISTS pipes (
    id              VARCHAR(64) PRIMARY KEY,
    name            VARCHAR(150) NOT NULL,
    diameter        INT NOT NULL,
    wall_size       INT NOT NULL,
    max_layer_mm    DOUBLE PRECISION NOT NULL DEFAULT 10.0,
    
    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS pipe_calibrations (
    id              SERIAL PRIMARY KEY,
    diameter        INT NOT NULL,
    wall_size       INT NOT NULL,
    energy_weight   DOUBLE PRECISION NOT NULL,
    duration_weight DOUBLE PRECISION NOT NULL,
    
    base_offset     DOUBLE PRECISION NOT NULL
);


CREATE TABLE IF NOT EXISTS sensor_readings (
    id          SERIAL PRIMARY KEY,
    calc_id     VARCHAR(64) NOT NULL,
    pipe_id     VARCHAR(64) NOT NULL REFERENCES pipes(id) ON DELETE CASCADE,
    sensor_num  INT NOT NULL,
    energy      DOUBLE PRECISION NOT NULL,
    duration    DOUBLE PRECISION NOT NULL,
    freq        DOUBLE PRECISION NOT NULL,
    amp         DOUBLE PRECISION NOT NULL,
    
    
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);


CREATE TABLE IF NOT EXISTS calc_results (
    id              VARCHAR(64) PRIMARY KEY,
    pipe_id         VARCHAR(64) NOT NULL REFERENCES pipes(id) ON DELETE CASCADE,
    avg_energy      DOUBLE PRECISION NOT NULL,
    avg_duration    DOUBLE PRECISION NOT NULL,
    top_layer_mm    DOUBLE PRECISION NOT NULL,
    bottom_layer_mm DOUBLE PRECISION NOT NULL,
    avg_layer_mm    DOUBLE PRECISION NOT NULL,
    days_left       INT NOT NULL,
    is_danger       BOOLEAN NOT NULL DEFAULT FALSE,
    source          VARCHAR(32) NOT NULL,

    created_at      TIMESTAMP NOT NULL DEFAULT NOW()
);
CREATE TABLE IF NOT EXISTS pipe_alerts (
    id          VARCHAR(64) PRIMARY KEY,
    pipe_id     VARCHAR(64) NOT NULL REFERENCES pipes(id) ON DELETE CASCADE,
    calc_id     VARCHAR(64) NOT NULL REFERENCES calc_results(id) ON DELETE CASCADE,
    layer_mm    DOUBLE PRECISION NOT NULL,
    message     VARCHAR(255) NOT NULL,
    is_fixed    BOOLEAN NOT NULL DEFAULT FALSE,
    
    created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
    fixed_at    TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_calc_results_pipe_id ON calc_results(pipe_id, created_at DESC);

INSERT INTO pipe_calibrations (diameter, wall_size, energy_weight, duration_weight, base_offset)
VALUES
    (159, 8, 0.085, 0.092, 16.5),
    (159, 10, 0.082, 0.089, 16.2),
    (108, 6, 0.090, 0.095, 15.8),
    (73, 6, 0.098, 0.102, 15.2);
