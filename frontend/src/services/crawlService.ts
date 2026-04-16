import { RunCrawl } from '../../wailsjs/go/crawl/CrawlHandler'
import { crawl } from '../../wailsjs/go/models'

export const crawlService = {
  runCrawl: async (req: crawl.CrawlRequest): Promise<crawl.CrawlResponse> => {
    return await RunCrawl(req as any)
  }
}
