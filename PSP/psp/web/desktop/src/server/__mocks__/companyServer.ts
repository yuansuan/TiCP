/* Copyright (C) 2016-present, Yuansuan.cn */

const initialValue = [
  {
    id: '3NjBKWJ4rby',
    account_id: '3NjBKWNSRpW',
    biz_code: '123432',
    cloud_type: 'mixed',
    contact: '',
    create_name: 'userName',
    create_time: { seconds: 1576287864, nanos: 0 },
    create_uid: 'user1d',
    is_ys_cloud: 0,
    modify_name: 'johnny_zhou',
    modify_uid: '3TX6vB1Ni8W',
    name: '远算科技',
    phone: '',
    remark: 'test',
    status: 1,
    update_time: { seconds: 1592400675, nanos: 0 },
  },
  {
    id: '3M6WH7N7DUU',
    account_id: '3N5AecLTkqd',
    biz_code: '3LQpVSiQwM7',
    cloud_type: '',
    contact: '15912300000',
    create_name: 'xpliu',
    create_time: { seconds: 1574130021, nanos: 0 },
    create_uid: '3LQpVYiQwM7',
    is_ys_cloud: 0,
    modify_name: 'edwin_cai',
    modify_uid: '3TWGGEu835Y',
    name: '远算测试企业',
    phone: '15912300000',
    remark: '远算测试企业，不要删除',
    status: 1,
    update_time: { seconds: 1595827925, nanos: 0 },
  },
]

export const companyServer = {
  getList: async () => ({ data: initialValue }),
  initialValue,
}
