CREATE TABLE IF NOT EXISTS jobs (
    job_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    recruiter_id UUID NOT NULL, -- FK to recruiters table
    organisation_name VARCHAR(255),
    job_title VARCHAR(255) NOT NULL,
    job_description TEXT NOT NULL,
    job_location VARCHAR(100),
    salary INT,
    skills_required TEXT[],
    vacancy INT DEFAULT 1,
    status VARCHAR(50) DEFAULT 'active', -- active, closed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
