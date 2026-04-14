import { GetSessionStatus } from '../../wailsjs/go/actiontest/SessionHandler'
import { actiontest } from '../../wailsjs/go/models'
import type { SessionStatus } from '../types/action'

export const sessionService = {
  async getSessionStatus(accountId?: string): Promise<SessionStatus> {
    try {
      const resp = await GetSessionStatus()
      
      // Mock for specific accounts
      if (accountId === 'acc_002') {
        return {
          is_active: false,
          provider: 'system',
          message: 'Phiên hết hạn'
        }
      }

      // ensure we return something that conforms to our interface
      return {
        is_active: resp.is_active,
        provider: resp.provider,
        message: resp.message
      }
    } catch (e: any) {
      return {
        is_active: false,
        provider: 'unknown',
        message: 'Lỗi khi lấy thông tin session'
      }
    }
  }
}
