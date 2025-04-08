import CPUUsageReport from './Report/ReportChart/CPUUsageReport'
import MEMUsageReport from './Report/ReportChart/MEMUsageReport'
import ClusterCPUUsageReport from './Report/ReportChart/ClusterCPUUsageReport'
import ClusterMEMUsageReport from './Report/ReportChart/ClusterMEMUsageReport'
import IOReport from './Report/ReportChart/IOReport'
import DiskUTReport from './Report/ReportChart/DiskUTReport'
import CPUTimeReport from './Report/ReportChart/CPUTimeReport'
import JobCountReport from './Report/ReportChart/JobCountReport'
import JobDeliverCountReport from './Report/ReportChart/JobDeliverCountReport'
import JobWaitTimeReport from './Report/ReportChart/JobWaitTimeReport'
import LicenseAppReport from './Report/ReportChart/LicenseAppReport'
import LicenseAppModuleReport from './Report/ReportChart/LicenseAppModuleReport'
import VisualUsageDurationReport from './Report/ReportChart/VisualUsageDurationReport'
import VisualNumberStatistic from './Report/ReportChart/VisualNumberStatistic'
import NodeDownStatistic from './Report/ReportChart/NodeDownStatistic'

export const REPORT = {
  MEM_UT_AVG: {
    type: 'MEM_UT_AVG',
    ReportChart: MEMUsageReport,
    disableDates: false
  },
  CPU_UT_AVG: {
    type: 'CPU_UT_AVG',
    ReportChart: CPUUsageReport,
    disableDates: false
  },
  TOTAL_IO_UT_AVG: {
    type: 'TOTAL_IO_UT_AVG',
    ReportChart: IOReport,
    disableDates: false
  },
  DISK_UT_AVG: {
    type: 'DISK_UT_AVG',
    ReportChart: DiskUTReport,
    disableDates: false
  },
  CPU_TIME_SUM: {
    type: 'CPU_TIME_SUM',
    ReportChart: CPUTimeReport,
    disableDates: false
  },
  JOB_COUNT: {
    type: 'JOB_COUNT',
    ReportChart: JobCountReport,
    disableDates: false
  },
  JOB_DELIVER_COUNT: {
    type: 'JOB_DELIVER_COUNT',
    ReportChart: JobDeliverCountReport,
    disableDates: false
  },
  JOB_WAIT_STATISTIC: {
    type: 'JOB_WAIT_STATISTIC',
    ReportChart: JobWaitTimeReport,
    disableDates: false
  },
  CLUSTER_METRIC_CPU_UT_AVG: {
    type: 'CLUSTER_METRIC_CPU_UT_AVG',
    ReportChart: ClusterCPUUsageReport,
    disableDates: false
  },
  CLUSTER_METRIC_MEM_UT_AVG: {
    type: 'CLUSTER_METRIC_MEM_UT_AVG',
    ReportChart: ClusterMEMUsageReport,
    disableDates: false
  },
  LICENSE_APP_USED_UT_AVG: {
    type: 'LICENSE_APP_USED_UT_AVG',
    ReportChart: LicenseAppReport,
    disableDates: false
  },
  LICENSE_APP_MODULE_USED_UT_AVG: {
    type: 'LICENSE_APP_MODULE_USED_UT_AVG',
    ReportChart: LicenseAppModuleReport,
    disableDates: false
  },
  NODE_DOWN_STATISTIC: {
    type: 'NODE_DOWN_STATISTIC',
    ReportChart: NodeDownStatistic,
    disableDates: false
  },
  VISUAL_USAGE_DURATION: {
    type: 'VISUAL_USAGE_DURATION',
    ReportChart: VisualUsageDurationReport,
    disableDates: false
  },
  VISUAL_NUMBER_STATUS: {
    type: 'VISUAL_NUMBER_STATUS',
    ReportChart: VisualNumberStatistic,
    disableDates: false
  }
}

export const REPORT_TYPE_LABEL = {
  CLUSTER_METRIC_CPU_UT_AVG: '集群整体CPU平均利用率',
  CLUSTER_METRIC_MEM_UT_AVG: '集群整体内存平均利用率',
  MEM_UT_AVG: '内存平均利用率',
  CPU_UT_AVG: 'CPU平均利用率',
  TOTAL_IO_UT_AVG: '磁盘吞吐率',
  DISK_UT_AVG: '磁盘使用情况',
  CPU_TIME_SUM: '核时使用情况',
  JOB_COUNT: '作业投递数情况',
  JOB_DELIVER_COUNT: '用户数与作业数情况',
  JOB_WAIT_STATISTIC: '作业等待情况',
  LICENSE_APP_USED_UT_AVG: '许可证软件使用情况',
  // LICENSE_APP_MODULE_USED_UT_AVG: '许可证软件模块使用情况',
  NODE_DOWN_STATISTIC: '节点宕机统计',
  VISUAL_USAGE_DURATION: '可视化使用时长统计',
  VISUAL_NUMBER_STATUS: '可视化会话数量统计'
}

export const REPORT_TYPE_OPTIONS = [
  // 'CLUSTER_METRIC_CPU_UT_AVG',
  // 'CLUSTER_METRIC_MEM_UT_AVG',
  'CPU_UT_AVG',
  'MEM_UT_AVG',
  'TOTAL_IO_UT_AVG',
  'DISK_UT_AVG',
  'CPU_TIME_SUM',
  'JOB_COUNT',
  'JOB_DELIVER_COUNT',
  'JOB_WAIT_STATISTIC',
  'LICENSE_APP_USED_UT_AVG',
  // 'LICENSE_APP_MODULE_USED_UT_AVG'
  'NODE_DOWN_STATISTIC'
]

export const VISUAL_REPORT_TYPE_OPTIONS = [
  'VISUAL_USAGE_DURATION',
  'VISUAL_NUMBER_STATUS'
]

export const CLOUD_REPORT_TYPE_OPTIONS = [
  // 'CLUSTER_METRIC_CPU_UT_AVG',
  // 'CLUSTER_METRIC_MEM_UT_AVG',
  // 'CPU_UT_AVG',
  // 'MEM_UT_AVG',
  // 'TOTAL_IO_UT_AVG',
  // 'DISK_UT_AVG',
  'CPU_TIME_SUM',
  'JOB_COUNT',
  'JOB_DELIVER_COUNT',
  'JOB_WAIT_STATISTIC'
  // 'LICENSE_APP_USED_UT_AVG',
  // 'LICENSE_APP_MODULE_USED_UT_AVG'
]
