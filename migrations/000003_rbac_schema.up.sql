-- Create roles table
CREATE TABLE IF NOT EXISTS roles (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create user_roles join table
CREATE TABLE IF NOT EXISTS user_roles (
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

-- Create resources table
CREATE TABLE IF NOT EXISTS resources (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create actions table
CREATE TABLE IF NOT EXISTS actions (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL,
    description VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create policies table
CREATE TABLE IF NOT EXISTS policies (
    id SERIAL PRIMARY KEY,
    role_id INTEGER NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    resource_id INTEGER NOT NULL REFERENCES resources(id) ON DELETE CASCADE,
    action_id INTEGER NOT NULL REFERENCES actions(id) ON DELETE CASCADE,
    effect VARCHAR(10) DEFAULT 'allow' NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_role_resource_action UNIQUE (role_id, resource_id, action_id)
);

-- Seed default roles
INSERT INTO roles (name, description) VALUES
('admin', 'Administrator with full access permissions'),
('user', 'Standard user with basic profile access')
ON CONFLICT (name) DO NOTHING;

-- Seed default actions
INSERT INTO actions (name, description) VALUES
('create', 'Create resource'),
('read', 'Read resource'),
('update', 'Update resource'),
('delete', 'Delete resource')
ON CONFLICT (name) DO NOTHING;

-- Seed default resources
INSERT INTO resources (name, description) VALUES
('users', 'System users management'),
('roles', 'System roles management'),
('resources', 'System resource definitions'),
('actions', 'System authorization actions'),
('policies', 'Dynamic access policies')
ON CONFLICT (name) DO NOTHING;

-- Seed admin policies (allow all actions on all resources)
INSERT INTO policies (role_id, resource_id, action_id, effect)
SELECT r.id, res.id, act.id, 'allow'
FROM roles r, resources res, actions act
WHERE r.name = 'admin'
ON CONFLICT ON CONSTRAINT unique_role_resource_action DO NOTHING;

-- Seed standard user policies (allow read on users resource)
INSERT INTO policies (role_id, resource_id, action_id, effect)
SELECT r.id, res.id, act.id, 'allow'
FROM roles r, resources res, actions act
WHERE r.name = 'user' AND res.name = 'users' AND act.name = 'read'
ON CONFLICT ON CONSTRAINT unique_role_resource_action DO NOTHING;
