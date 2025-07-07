CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TYPE providers AS ENUM ('credentials', 'google');

CREATE TYPE user_roles AS ENUM ('user', 'admin');

CREATE TABLE
  users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4 (),
    username VARCHAR(50) UNIQUE NOT NULL,
    name VARCHAR(100) NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password TEXT,
    provider providers DEFAULT 'credentials',
    role user_roles DEFAULT 'user',
    jwt_version VARCHAR(20) NOT NULL,
    is_verified BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP(0)
    WITH
      TIME ZONE NOT NULL DEFAULT NOW (),
      updated_at TIMESTAMP(0)
    WITH
      TIME ZONE NOT NULL DEFAULT NOW ()
  );