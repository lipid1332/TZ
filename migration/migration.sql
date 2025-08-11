CREATE TABLE wallets(
    id uuid primary key,
    amount bigint DEFAULT 0,
    CHECK (amount>=0)
);

INSERT INTO wallets (id) VALUES
  ('156f5065-4a38-4abe-bf06-2dea11850408'),
  ('156f5065-4a38-1abe-bf06-2dea11850405'),
  ('156f5065-4a38-9abe-bf06-2dea11850405');