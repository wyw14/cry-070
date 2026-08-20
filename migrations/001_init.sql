CREATE TABLE IF NOT EXISTS data_sources (id text PRIMARY KEY,name text NOT NULL,adapter text NOT NULL,encrypted_dsn text NOT NULL,state text NOT NULL,version bigint NOT NULL DEFAULT 1);
CREATE TABLE IF NOT EXISTS classified_fields (id text PRIMARY KEY,source_id text NOT NULL,table_name text NOT NULL,field_name text NOT NULL,sensitivity text NOT NULL,scope jsonb NOT NULL DEFAULT '[]');
CREATE TABLE IF NOT EXISTS masking_rules (id text PRIMARY KEY,name text NOT NULL,kind text NOT NULL,pattern text,version integer NOT NULL,active boolean NOT NULL);
CREATE TABLE IF NOT EXISTS masking_pipelines (id text PRIMARY KEY,name text NOT NULL,source_id text NOT NULL,state text NOT NULL,version integer NOT NULL);
CREATE TABLE IF NOT EXISTS masking_batches (id text PRIMARY KEY,pipeline_id text NOT NULL,preview_id text NOT NULL,state text NOT NULL,idempotency_key text UNIQUE,total integer NOT NULL,processed integer NOT NULL);
CREATE TABLE IF NOT EXISTS masking_audits (id text PRIMARY KEY,actor text NOT NULL,action text NOT NULL,resource text NOT NULL,created_at timestamptz NOT NULL);
