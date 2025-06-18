-- Create users table
CREATE TABLE users (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    password TEXT NOT NULL,
    role TEXT CHECK (role IN ('admin', 'client')) NOT NULL,
    device_id TEXT,
    location TEXT,
    ip_address TEXT,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create sessions table
CREATE TABLE sessions (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
    jwt_access_token TEXT,
    refresh_token TEXT,
    expiration TIMESTAMPTZ,
    login_time TIMESTAMPTZ DEFAULT now(),
    active_tokens INT DEFAULT 0,
    revocations INT DEFAULT 0
);

-- Create applications table
CREATE TABLE applications (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    owner_id BIGINT REFERENCES users (id) ON DELETE SET NULL,
    is_template BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create collaborators table
CREATE TABLE collaborators (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    application_id BIGINT REFERENCES applications (id) ON DELETE CASCADE,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
    role TEXT CHECK (role IN ('owner','editor','viewer')) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now(),
    CONSTRAINT unique_user_per_application UNIQUE (application_id, user_id)
);

-- Create app_connections table
CREATE TABLE app_connections (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    application_id BIGINT REFERENCES applications (id) ON DELETE CASCADE,
    endpoint_name TEXT NOT NULL,
    method TEXT NOT NULL,
    url TEXT NOT NULL,
    headers JSONB,
    auth JSONB,
    params JSONB,
    response_mapping JSONB,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create screens table
CREATE TABLE screens (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    application_id BIGINT REFERENCES applications (id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    route TEXT NOT NULL,
    screen_type TEXT NOT NULL,
    device_type TEXT NOT NULL,
    description TEXT,
    status TEXT CHECK (status IN ('published','draft')) NOT NULL,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create widgets table
CREATE TABLE widgets (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    screen_id BIGINT REFERENCES screens (id) ON DELETE CASCADE,
    parent_id BIGINT REFERENCES widgets (id) ON DELETE CASCADE,
    properties JSONB,
    actions JSONB,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create admin_dashboard table
CREATE TABLE admin_dashboard (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    metrics JSONB,
    system_settings JSONB,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);

-- Create client_dashboard table
CREATE TABLE client_dashboard (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT REFERENCES users (id) ON DELETE CASCADE,
    projects_metadata JSONB,
    theme_config JSONB,
    api_manager JSONB,
    versioning JSONB,
    feedback JSONB,
    profile_settings JSONB,
    created_at TIMESTAMPTZ DEFAULT now(),
    updated_at TIMESTAMPTZ DEFAULT now()
);
