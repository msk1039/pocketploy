// Instance type (active instances only - deleted ones are in ArchivedInstance)
export interface Instance {
  id: string;
  user_id: string;
  name: string;
  slug: string;
  subdomain: string;
  container_id?: string;
  container_name?: string;
  status: 'creating' | 'running' | 'stopped' | 'failed';
  data_path: string;
  created_at: string;
  updated_at: string;
  last_accessed_at?: string;
}

// Archived Instance type (for deleted instances with restore capability)
export interface ArchivedInstance {
  id: string;
  user_id: string;
  name: string;
  slug: string;
  subdomain: string;
  container_id?: string;
  container_name?: string;
  original_status: string;
  data_path: string;
  created_at: string;
  updated_at: string;
  last_accessed_at?: string;
  deleted_at: string;
  deleted_by_user_id: string;
  deletion_reason: string;
  data_available: boolean;
  data_retained_until: string;
  data_size_mb: number;
  original_subdomain: string;
}

// Instance API Request types
export interface CreateInstanceRequest {
  name: string;
}

// Instance API Response types
export interface CreateInstanceResponse {
  success: boolean;
  message: string;
  instance: Instance;
  url: string;
}

export interface ListInstancesResponse {
  success: boolean;
  instances: Instance[];
}

export interface GetInstanceResponse {
  success: boolean;
  instance: Instance;
}

export interface DeleteInstanceResponse {
  success: boolean;
  message: string;
}

// Archived instances API response types
export interface ListArchivedInstancesResponse {
  success: boolean;
  instances: ArchivedInstance[];
}

export interface GetArchivedInstanceResponse {
  success: boolean;
  instance: ArchivedInstance;
}
