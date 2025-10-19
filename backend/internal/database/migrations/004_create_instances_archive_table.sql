-- Drop existing archive table if it exists (for clean recreation)
DROP TABLE IF EXISTS instances_archive CASCADE;

-- Create instances_archive table for deleted instances with restore capability
CREATE TABLE instances_archive (
    -- Original instance fields
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(100) NOT NULL,
    slug VARCHAR(100) NOT NULL,
    subdomain VARCHAR(255) NOT NULL,
    container_id VARCHAR(255),
    container_name VARCHAR(255),
    
    -- Status before deletion
    original_status VARCHAR(20) NOT NULL,
    
    -- File paths
    data_path TEXT,
    
    -- Original lifecycle timestamps
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    last_accessed_at TIMESTAMP,
    
    -- Archive-specific metadata
    deleted_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_by_user_id UUID NOT NULL REFERENCES users(id) ON DELETE SET NULL,
    deletion_reason VARCHAR(50) DEFAULT 'manual',
    
    -- Data retention tracking
    data_available BOOLEAN NOT NULL DEFAULT true,
    data_retained_until TIMESTAMP NOT NULL,
    data_size_mb INTEGER DEFAULT 0,
    
    -- Keep original subdomain for reference
    original_subdomain VARCHAR(255) NOT NULL,
    
    -- Indexes for common queries
    CONSTRAINT instances_archive_user_slug_deleted_unique UNIQUE(user_id, slug, deleted_at)
);

-- Create indexes for efficient querying
CREATE INDEX idx_instances_archive_user_id ON instances_archive(user_id);
CREATE INDEX idx_instances_archive_deleted_at ON instances_archive(deleted_at);
CREATE INDEX idx_instances_archive_data_retained_until ON instances_archive(data_retained_until);
CREATE INDEX idx_instances_archive_data_available ON instances_archive(data_available);

-- Comments for documentation
COMMENT ON TABLE instances_archive IS 'Archived instances with 30-day data retention for restore capability';
COMMENT ON COLUMN instances_archive.original_status IS 'Status of instance before deletion (running, stopped, failed, creating)';
COMMENT ON COLUMN instances_archive.data_available IS 'Whether instance data files still exist on disk';
COMMENT ON COLUMN instances_archive.data_retained_until IS 'Date when data files will be permanently deleted (deleted_at + 30 days)';
COMMENT ON COLUMN instances_archive.deletion_reason IS 'Reason for deletion: manual, auto_cleanup, admin, etc';
COMMENT ON COLUMN instances_archive.data_size_mb IS 'Size of instance data in MB at time of deletion';
