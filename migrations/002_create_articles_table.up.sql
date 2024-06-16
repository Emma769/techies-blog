CREATE TABLE IF NOT EXISTS articles (
  article_id TEXT DEFAULT REPLACE(uuid_generate_v4()::text, '-', ''),
  title TEXT NOT NULL, 
  slug TEXT NOT NULL UNIQUE, 
  description TEXT, 
  content TEXT NOT NULL,
  version INT DEFAULT 1,
  created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW(),
  updated_at TIMESTAMP WITH TIME ZONE,
  PRIMARY KEY(article_id)
);

CREATE OR REPLACE TRIGGER articles_update_timestamp_trigger BEFORE UPDATE ON articles FOR EACH ROW EXECUTE PROCEDURE update_timestamp();

CREATE OR REPLACE TRIGGER articles_update_version_trigger BEFORE UPDATE ON articles FOR EACH ROW EXECUTE PROCEDURE update_version();
