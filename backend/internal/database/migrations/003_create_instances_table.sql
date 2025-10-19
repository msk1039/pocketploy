-- Create instances table
CREATE TABLE instances (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(150) NOT NULL,
    subdomain VARCHAR(255) NOT NULL UNIQUE,
    container_id VARCHAR(255),
    container_name VARCHAR(255) UNIQUE,
    status VARCHAR(20) NOT NULL DEFAULT 'pending',
    storage_path TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    last_accessed_at TIMESTAMP
);

-- Create indexes for better query performance
CREATE INDEX idx_instances_user_id ON instances(user_id);
CREATE INDEX idx_instances_status ON instances(status);
CREATE INDEX idx_instances_subdomain ON instances(subdomain);
CREATE INDEX idx_instances_container_id ON instances(container_id);

-- Add comment to table
COMMENT ON TABLE instances IS 'Stores PocketBase instance information for each user';
COMMENT ON COLUMN instances.status IS 'Instance status: pending, running, stopped, failed, deleted';
COMMENT ON COLUMN instances.slug IS 'URL-safe version of instance name';
COMMENT ON COLUMN instances.subdomain IS 'Full subdomain for accessing the instance (e.g., username-slug.192.168.1.100.nip.io)';
