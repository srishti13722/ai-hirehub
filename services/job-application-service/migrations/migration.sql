CREATE TABLE IF NOT EXISTS job_applications (
    application_id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    job_id UUID NOT NULL,
    jobseeker_id UUID NOT NULL,
    recruiter_id UUID NOT NULL,
    status VARCHAR(50) DEFAULT 'applied', -- applied, shortlisted, rejected, offered
    applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
