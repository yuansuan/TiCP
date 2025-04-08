//status: 不同调度器下的 schedule_status
export const LSFCanCloseStatus = ['OK', 'Closed_Full', 'Closed_Busy']
export const LSFCanOpenStatus = 'Closed_Adm'
export const PBSCanCloseStatus = ['free', 'busy', 'job-busy', 'job-exclusive']
export const PBSCanOpenStatus = 'offline'
// idle 表示当前节点空闲，可以点击拒绝作业。
// down 表示当前节点下线，可以点击接收作业。
export const PSPCanOpenStatus = ['down']
export const PSPCanCloseStatus = ['idle']
