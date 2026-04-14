import { LikePost, CommentPost, CreatePost, GetActionLogs } from '../../wailsjs/go/actiontest/ActionHandler'
import { actiontest } from '../../wailsjs/go/models'
import type { ActionResponse, ActionLog, ActionRequest, CommentRequest, CreatePostRequest } from '../types/action'

export const actionsService = {
  async likePost(reqData: ActionRequest): Promise<ActionResponse> {
    // doc_id KHÔNG được truyền — backend tự resolve từ graphqlDocID constant nội bộ
    console.debug('[DocID-Bridge] Building LikePostRequest WITHOUT doc_id:', {
      post_id: reqData.post_id,
      account_id: reqData.account_id,
      reaction_type: reqData.reaction_type,
      reaction_id: reqData.reaction_id,
      dry_run: reqData.dry_run
    })
    const req = new actiontest.LikePostRequest({
      post_id: reqData.post_id,
      account_id: reqData.account_id,
      reaction_type: reqData.reaction_type,
      reaction_id: reqData.reaction_id,
      dry_run: reqData.dry_run,
      actor_source: reqData.actor_source || 'manual_test'
    })

    try {
      const resp = await LikePost(req)
      return {
        success: resp.success,
        status: resp.status,
        action: resp.action,
        post_id: resp.post_id,
        account_id: resp.account_id,
        account_display_name: resp.account_display_name,
        reaction_type: resp.reaction_type,
        reaction_id: resp.reaction_id,
        dry_run: resp.dry_run,
        message: resp.message,
        executed_at: resp.executed_at as unknown as string
      }
    } catch (e: any) {
      return {
        success: false,
        status: 'system_error',
        action: reqData.reaction_type || 'reaction',
        post_id: reqData.post_id,
        account_id: reqData.account_id,
        reaction_type: reqData.reaction_type,
        reaction_id: reqData.reaction_id,
        dry_run: reqData.dry_run,
        message: 'Lỗi kết nối backend: ' + (e?.message || e?.toString() || 'Unknown error'),
        executed_at: new Date().toISOString()
      }
    }
  },

  async commentPost(reqData: CommentRequest): Promise<ActionResponse> {
    const req = new actiontest.CommentPostRequest({
      post_id: reqData.post_id,
      account_id: reqData.account_id,
      comment_text: reqData.comment_text,
      dry_run: reqData.dry_run,
      actor_source: reqData.actor_source || 'manual_test'
    })

    try {
      const resp = await CommentPost(req)
      return {
        success: resp.success,
        status: resp.status,
        action: resp.action,
        post_id: resp.post_id,
        account_id: resp.account_id,
        account_display_name: resp.account_display_name,
        comment_text: resp.comment_text,
        dry_run: resp.dry_run,
        message: resp.message,
        executed_at: resp.executed_at as unknown as string
      }
    } catch (e: any) {
      return {
        success: false,
        status: 'system_error',
        action: 'comment',
        post_id: reqData.post_id,
        account_id: reqData.account_id,
        comment_text: reqData.comment_text,
        dry_run: reqData.dry_run,
        message: 'Lỗi kết nối backend: ' + (e?.message || e?.toString() || 'Unknown error'),
        executed_at: new Date().toISOString()
      }
    }
  },

  async createPost(reqData: CreatePostRequest): Promise<ActionResponse> {
    const req = new actiontest.CreatePostRequest({
      account_id: reqData.account_id,
      post_text: reqData.post_text,
      image_paths: reqData.image_paths,
      dry_run: reqData.dry_run,
      actor_source: reqData.actor_source || 'manual_test'
    })

    try {
      const resp = await CreatePost(req)
      return {
        success: resp.success,
        status: resp.status,
        action: resp.action,
        post_id: '',
        account_id: resp.account_id,
        account_display_name: resp.account_display_name,
        post_text: resp.post_text,
        dry_run: resp.dry_run,
        message: resp.message,
        executed_at: resp.executed_at as unknown as string
      }
    } catch (e: any) {
      return {
        success: false,
        status: 'system_error',
        action: 'create_post',
        post_id: '',
        account_id: reqData.account_id,
        post_text: reqData.post_text,
        dry_run: reqData.dry_run,
        message: 'Lỗi kết nối backend: ' + (e?.message || e?.toString() || 'Unknown error'),
        executed_at: new Date().toISOString()
      }
    }
  },

  async getActionLogs(): Promise<ActionLog[]> {
    try {
      const logs = await GetActionLogs()
      return logs.map((l: actiontest.ActionLog) => ({
        id: l.id,
        action_type: l.action_type,
        account_id: l.account_id,
        account_display_name: l.account_display_name,
        post_id: l.post_id,
        reaction_type: l.reaction_type,
        reaction_id: l.reaction_id,
        status: l.status,
        dry_run: l.dry_run,
        message: l.message,
        executed_at: l.executed_at as unknown as string
      }))
    } catch (e) {
      return []
    }
  }
}
