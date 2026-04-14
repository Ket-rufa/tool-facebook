export interface ActionRequest {
  post_id: string;
  dry_run: boolean;
  actor_source: string;
}

export interface ActionResponse {
  success: boolean;
  status: string;
  action: string;
  post_id: string;
  dry_run: boolean;
  message: string;
  executed_at: string;
}

export interface SessionStatus {
  is_active: boolean;
  provider: string;
  message: string;
}

export interface ActionLog {
  id: string;
  action_type: string;
  post_id: string;
  status: string;
  dry_run: boolean;
  message: string;
  executed_at: string;
}
