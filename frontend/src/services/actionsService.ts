import { LikePost, GetActionLogs } from '../../wailsjs/go/actiontest/ActionHandler'
import { actiontest } from '../../wailsjs/go/models'
import type { ActionResponse, ActionLog } from '../types/action'

export const actionsService = {
  async likePost(postId: string, dryRun: boolean): Promise<ActionResponse> {
    const req = new actiontest.LikePostRequest({
      post_id: postId,
      dry_run: dryRun,
      actor_source: 'manual_test'
    })
    
    try {
      const resp = await LikePost(req)
      return {
        success: resp.success,
        status: resp.status,
        action: resp.action,
        post_id: resp.post_id,
        dry_run: resp.dry_run,
        message: resp.message,
        // Because ExecutedAt comes from Go time.Time, it might be a string in JS
        executed_at: resp.executed_at as unknown as string
      }
    } catch (e: any) {
      return {
        success: false,
        status: 'system_error',
        action: 'like',
        post_id: postId,
        dry_run: dryRun,
        message: e.toString(),
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
        post_id: l.post_id,
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
