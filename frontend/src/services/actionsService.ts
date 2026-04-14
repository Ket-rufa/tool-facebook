import { LikePost, GetActionLogs } from '../../wailsjs/go/actiontest/ActionHandler'
import { actiontest } from '../../wailsjs/go/models'
import type { ActionResponse, ActionLog, ActionRequest } from '../types/action'

export const actionsService = {
  async likePost(reqData: ActionRequest): Promise<ActionResponse> {
    // Payload đầy đủ — backend bây giờ nhận account_id + reaction_type thật
    const req = new actiontest.LikePostRequest({
      post_id: reqData.post_id,
      account_id: reqData.account_id,
      reaction_type: reqData.reaction_type,
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
