import { formatlsfAttrNumber } from '@/utils/formatter'

const round = formatlsfAttrNumber

export const TABLE_CONF = {
  id: 'node_table',
  columns: [
    { key: 'node_name', active: true }, // 机器名
    { key: 'scheduler_status', active: true }, // 调度器状态
    { key: 'queue_name', active: true }, // 队列名称
    { key: 'total_core_num', active: true }, // 总核数
    { key: 'used_core_num', active: true }, // 使用核数
    { key: 'free_core_num', active: true }, // 空闲核数
    { key: 'total_mem', active: true }, // 总内存
    { key: 'used_mem', active: true }, // 使用内存 MB
    { key: 'free_mem', active: true }, // 空闲内存 MB
    { key: 'available_mem', active: true }, //可用内存 MB
    { key: 'node_type', active: true } // 机器类型
  ]
}

interface SectionConfig {
  title: string
  icon: string
  partitionNum: number
  resourceKey?: string
  children: ItemConfig[]
}

interface ItemConfig {
  key: string
  text: string
  keyTip?: string
  formatter?: (values: any) => string
}

export const NODE_INFO_CONF: SectionConfig[] = [
  {
    title: '节点信息',
    icon: null,
    partitionNum: 2,
    children: [
      { key: 'node_name', text: '机器名', keyTip: '' },
      { key: 'node_status', text: '机器状态', keyTip: '' },
      { key: 'scheduler_status', text: '调度器状态', keyTip: '' },
      { key: 'n_core', text: '总核数', keyTip: '' },
      {
        key: 'ut',
        text: 'CPU使用率',
        keyTip: '',
        formatter: values => `${round(values['ut'])}%`
      },
      {
        key: 'mem',
        text: '内存(空闲/最大)',
        keyTip: '',
        formatter: values =>
          `${round(values['mem'])}/${round(values['max_mem'])}MB`
      },
      {
        key: 'swap',
        text: '交换空间(空闲/最大)',
        keyTip: '',
        formatter: values =>
          `${round(values['swap'])}/${round(values['max_swap'])}MB`
      },
      {
        key: 'tmp',
        text: '/tmp的空间(空闲/最大)',
        keyTip: '',
        formatter: values =>
          `${round(values['tmp'])}/${round(values['max_tmp'])}MB`
      },
      { key: 'n_disk', text: '磁盘数量', keyTip: '' },
      {
        key: 'io',
        text: '磁盘吞吐率',
        keyTip: '',
        formatter: values => `${round(values['io'])}KB/s`
      },
      {
        key: 'it',
        text: '空闲时间',
        keyTip: '',
        formatter: values => `${round(values['it'])}min`
      },
      {
        key: 'r',
        text: '15s/1m/15m负载',
        keyTip: '',
        formatter: values =>
          `${round(values['r15s'])}/${round(values['r1m'])}/${round(
            values['r15m']
          )}`
      },
      { key: 'resource_attr', text: '机器资源', keyTip: '' },
      { key: 'mac', text: 'mac地址', keyTip: '' }
    ]
  }
  // {
  //   title: 'NUMA节点',
  //   icon: null,
  //   partitionNum: 2,
  //   children: [],
  // },
]
