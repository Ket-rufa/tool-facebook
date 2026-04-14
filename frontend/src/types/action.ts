// doc_id đã bị xóa — backend tự resolve nội bộ từ graphqlDocID constant trong executor.go
export interface ActionRequest {
  post_id: string
  account_id: string
  reaction_type: string
  reaction_id: string
  dry_run: boolean
  actor_source: string
}

export interface CommentRequest {
  post_id: string
  account_id: string
  comment_text: string
  dry_run: boolean
  actor_source: string
}

export interface CreatePostRequest {
  account_id: string
  post_text: string
  image_paths?: string[]
  dry_run: boolean
  actor_source: string
}


export interface ActionResponse {
  success: boolean
  status: string
  action: string
  post_id: string
  account_id: string
  account_display_name?: string
  reaction_type: string
  reaction_id: string
  dry_run: boolean
  message: string
  executed_at: string
}

export interface SessionStatus {
  is_active: boolean
  provider: string
  message: string
}

export interface ActionLog {
  id: string
  action_type: string
  account_id: string
  account_display_name?: string
  post_id: string
  reaction_type: string
  reaction_id: string
  status: string
  dry_run: boolean
  message: string
  executed_at: string
}
