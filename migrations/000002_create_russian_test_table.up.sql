CREATE SCHEMA IF NOT EXISTS public;

CREATE TABLE IF NOT EXISTS public.russian_test (
    data VARCHAR(255) PRIMARY KEY UNIQUE NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_russian_test_data 
ON public.russian_test (data text_pattern_ops);