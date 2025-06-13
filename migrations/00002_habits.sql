-- +goose Up
-- +goose StatementBegin
CREATE TABLE habits (
    id UUID PRIMARY KEY ,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    motivation TEXT,
    color VARCHAR(7) DEFAULT '#3B82F6',
    category VARCHAR(100),
    frequency VARCHAR(20) NOT NULL DEFAULT 'daily',
    target_count INTEGER NOT NULL DEFAULT 1,
    target_days JSONB,
    current_streak INTEGER NOT NULL DEFAULT 0,
    best_streak INTEGER NOT NULL DEFAULT 0,
    total_completions INTEGER NOT NULL DEFAULT 0,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    CONSTRAINT habits_frequency_check CHECK (frequency IN ('daily', 'weekly', 'monthly', 'custom')),
    CONSTRAINT habits_target_count_positive CHECK (target_count > 0),
    CONSTRAINT habits_color_format CHECK (color ~ '^#[0-9A-Fa-f]{6}$')
);

CREATE TRIGGER set_timestamp_habits
    BEFORE UPDATE ON habits
    FOR EACH ROW
    EXECUTE FUNCTION trigger_set_timestamp();


-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS set_timestamp_habits ON habits;
DROP TABLE IF EXISTS habits; 
-- +goose StatementEnd