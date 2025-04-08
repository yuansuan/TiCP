/**
 * 默认协程池的最大并发数量
 */
const DEFAULT_CONCURRENCY = 10

/**
 * 默认协程池的最大任务数量
 */
const DEFAULT_MAX_RUNNABLE = 5

/**
 * 创建协程池的控制参数
 */
export interface AsyncPoolProps {
  /**
   * 协程池的最大并发数量
   */
  concurrency?: number

  /**
   * 协程池最多可运行的任务数量
   */
  maxRunnable?: number
}

/**
 * 任务描述
 */
export interface RunnableProps<T extends unknown> {
  /**
   * 用于返回下一个异步任务的参数, 如果返回 null 或者 undefined 则
   * 终止运行该任务
   */
  next: () => T

  /**
   * 异步任务
   */
  work: (args?: T) => Promise<any>

  /**
   * 最大并发数量
   */
  maxConcurrency?: number
}

/**
 * 异步并发协程池实现
 *
 * 对于某一个任务(Promise), 可以指定其最多运行几个和最少运行几个
 */
export default class AsyncPool {
  /**
   * 协程池的最大并发数量
   */
  private concurrency: number

  /**
   * 当前运行的协程组
   */
  private globalExecuting = new Set()

  /**
   * 协程池最多可运行的任务数量
   */
  private maxRunnable: number

  /**
   * 当前正在运行的任务
   */
  private globalRunning = new Set()

  constructor(props?: AsyncPoolProps) {
    this.concurrency = props?.concurrency || DEFAULT_CONCURRENCY
    this.maxRunnable = props?.maxRunnable || DEFAULT_MAX_RUNNABLE
  }

  /**
   * 仅执行一次任务
   */
  async once(work: (args?: any) => Promise<any>): Promise<any> {
    let index = 0
    return await this.run({
      next: () => (index == 0 ? index++ : null),
      work: work
    })
  }

  /**
   * 启动一个任务并等待完成
   */
  async run<T>(runnable: RunnableProps<T>): Promise<any> {
    // 默认情况下可以占用协程池的一半并发量
    if (!runnable.maxConcurrency) {
      runnable.maxConcurrency = this.concurrency / 2
    }

    const result = []
    const localExecuting = new Set()

    const enqueue = async () => {
      const args = runnable.next()
      if (args === null || args === undefined) {
        return Promise.resolve()
      }

      const p = runnable.work(args)
      result.push(p)
      localExecuting.add(p)
      this.globalExecuting.add(p)

      p.finally(() => {
        localExecuting.delete(p)
        this.globalExecuting.delete(p)
      })

      while (
        localExecuting.size >= runnable.maxConcurrency ||
        this.globalExecuting.size >= this.concurrency
      ) {
        // 等待一个任务完成
        await Promise.race([
          Promise.race(localExecuting),
          Promise.race(this.globalExecuting)
        ])
      }

      return await enqueue()
    }

    return await this.limitRunnable(async () => {
      return await enqueue().then(() => Promise.all(result))
    })
  }

  /**
   * 限制同时执行的任务数量
   */
  private async limitRunnable(work: () => Promise<any>): Promise<any> {
    while (this.globalRunning.size >= this.maxRunnable) {
      await Promise.race(this.globalRunning)
    }

    const w = work()
    this.globalRunning.add(w)

    w.finally(() => {
      this.globalRunning.delete(w)
    })

    return await w
  }
}
