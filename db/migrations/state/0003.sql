-- +migrate Up
CREATE TABLE state.proof2
(
    batch_num  BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    batch_num_final BIGINT NOT NULL REFERENCES state.batch (batch_num) ON DELETE CASCADE,
    proof jsonb,
    proof_id VARCHAR,
    input_prover jsonb,
    prover VARCHAR,
    aggregating BOOLEAN DEFAULT FALSE,
    PRIMARY KEY (batch_num, batch_num_final)    
);

-- +migrate Down
DROP TABLE state.proof2;