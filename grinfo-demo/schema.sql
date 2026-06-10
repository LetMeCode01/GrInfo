CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS user_stats (
    user_id BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    total_xp INTEGER NOT NULL DEFAULT 0,
    current_streak INTEGER NOT NULL DEFAULT 0,
    longest_streak INTEGER NOT NULL DEFAULT 0,
    last_activity_date DATE,
    lessons_completed INTEGER NOT NULL DEFAULT 0,
    quizzes_completed INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS grinfo_categories (
    id BIGSERIAL PRIMARY KEY,
    slug TEXT UNIQUE NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS grinfo_questions (
    id BIGSERIAL PRIMARY KEY,
    category_id BIGINT NOT NULL REFERENCES grinfo_categories(id) ON DELETE RESTRICT,
    difficulty TEXT NOT NULL CHECK (difficulty IN ('usoara', 'medie', 'grea')),
    elo_rating INTEGER NOT NULL,
    enunt TEXT NOT NULL,
    explicatie_raspuns TEXT NOT NULL,
    graph_data JSONB NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS grinfo_question_options (
    id BIGSERIAL PRIMARY KEY,
    question_id BIGINT NOT NULL REFERENCES grinfo_questions(id) ON DELETE CASCADE,
    option_index INTEGER NOT NULL CHECK (option_index BETWEEN 0 AND 3),
    option_text TEXT NOT NULL,
    is_correct BOOLEAN NOT NULL DEFAULT FALSE,
    UNIQUE(question_id, option_index)
);

CREATE TABLE IF NOT EXISTS grinfo_sessions (
    id BIGSERIAL PRIMARY KEY,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    category TEXT NOT NULL DEFAULT 'all',
    initial_elo NUMERIC(6,2) NOT NULL DEFAULT 1000,
    final_elo NUMERIC(6,2),
    total_questions INTEGER NOT NULL DEFAULT 10,
    correct_answers INTEGER NOT NULL DEFAULT 0,
    started_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    ended_at TIMESTAMPTZ
);

CREATE TABLE IF NOT EXISTS grinfo_session_answers (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT NOT NULL REFERENCES grinfo_sessions(id) ON DELETE CASCADE,
    question_id BIGINT NOT NULL REFERENCES grinfo_questions(id) ON DELETE RESTRICT,
    is_correct BOOLEAN NOT NULL,
    elo_before NUMERIC(6,2) NOT NULL,
    elo_after NUMERIC(6,2) NOT NULL,
    answered_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS grinfo_security_logs (
    id BIGSERIAL PRIMARY KEY,
    session_id BIGINT REFERENCES grinfo_sessions(id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users(id) ON DELETE SET NULL,
    incident_type TEXT NOT NULL,
    description TEXT,
    elo_penalty INTEGER NOT NULL DEFAULT 50,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_grinfo_questions_category ON grinfo_questions(category_id);
CREATE INDEX IF NOT EXISTS idx_grinfo_questions_elo ON grinfo_questions(elo_rating);
CREATE INDEX IF NOT EXISTS idx_grinfo_questions_active ON grinfo_questions(is_active);
CREATE INDEX IF NOT EXISTS idx_grinfo_sessions_user ON grinfo_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_grinfo_answers_session ON grinfo_session_answers(session_id);
CREATE INDEX IF NOT EXISTS idx_grinfo_security_session ON grinfo_security_logs(session_id);

INSERT INTO grinfo_categories (slug, name)
VALUES
('orientate', 'Grafuri orientate'),
('neorientate', 'Grafuri neorientate')
ON CONFLICT (slug) DO NOTHING;

-- Seed-ul din backend asigura minim 50 de intrebari active:
-- 1) 25 intrebari pe orientate + 25 pe neorientate in grinfo_questions
-- 2) Pentru fiecare intrebare inserati exact 4 randuri in grinfo_question_options
-- 3) graph_data trebuie sa respecte formatul:
--    {"nodes":[{"id":"1"}], "edges":[{"source":"1","target":"2"}]}
