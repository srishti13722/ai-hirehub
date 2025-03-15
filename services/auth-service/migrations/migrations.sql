CREATE TABLE IF NOT EXISTS users(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    firstname VARCHAR(20) NOT NULL,
    lastname VARCHAR(20) NOT NULL,
    email VARCHAR(50) NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) CHECK (role in('job_seeker', 'recruiter', 'admin')) NOT NULL DEFAULT 'job_seeker',
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);