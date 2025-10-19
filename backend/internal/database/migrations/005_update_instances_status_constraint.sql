-- Migrate existing deleted instances to archive table before updating constraint
-- This handles any instances that were soft-deleted before the archive table existed

-- Step 1: Move all deleted instances to archive
INSERT INTO instances_archive (
    id, user_id, name, slug, subdomain,
    container_id, container_name, original_status, data_path,
    created_at, updated_at, deleted_at, deleted_by_user_id, last_accessed_at,
    data_size_mb, data_available, data_retained_until, original_subdomain
)
SELECT 
    id, user_id, name, slug, subdomain,
    container_id, container_name, 
    'deleted' as original_status,  -- We don't know original status, mark as deleted
    data_path,
    created_at, updated_at, 
    updated_at as deleted_at,  -- Use updated_at as deletion time
    user_id as deleted_by_user_id,  -- Assume owner deleted it
    last_accessed_at,
    0 as data_size_mb,  -- Set to 0 instead of NULL
    false as data_available,  -- Data likely already cleaned up
    updated_at as data_retained_until,  -- Already expired
    subdomain as original_subdomain  -- Store original subdomain
FROM instances
WHERE status = 'deleted';

-- Step 2: Delete those instances from main table
DELETE FROM instances WHERE status = 'deleted';

-- Step 3: Drop the old constraint
ALTER TABLE instances DROP CONSTRAINT IF EXISTS instances_status_check;

-- Step 4: Add new constraint without 'deleted' status
ALTER TABLE instances ADD CONSTRAINT instances_status_check 
    CHECK (status IN ('creating', 'running', 'stopped', 'failed'));

-- Comment for clarity
COMMENT ON COLUMN instances.status IS 'Current status: creating, running, stopped, or failed. Deleted instances move to instances_archive table';
