/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import { observable, action, computed, runInAction } from 'mobx'
import { message } from 'antd'
import { companyServer } from '@/server'
import { env } from '..'

interface resultInterface {
  zone: string
  storage_domains: string[]
  sc_id?: string
  domain?: string
}

export class BaseDomain {
  @observable result: Array<resultInterface>
  @observable zoneSelectData: Array<resultInterface>
}

export class Domain extends BaseDomain {
  constructor() {
    super()
  }

  @action
  async fetch() {
    try {
      const { data, success } = await companyServer.getDomain(env.company_id)
      const { data: scRes, success: scSuccess } =
        await companyServer.getSCList()

      if (!success && !scSuccess) {
        return message.error('区域下拉列表数据获取异常')
      }
      // 定义一个收集区域下拉数据变量
      const collectZoneSelectData = []
      // 合并三个接口所有数据
      const allDataRes = [].concat(data.result, scRes.result, visData.items)
      allDataRes.forEach((item: any) => {
        // 如果item.name是az-yuansuan
        item.name === 'az-yuansuan' && (item.name = 'az-shanghai')
        // 判断collectZoneSelectData是否已有对应数据
        const _findIndex = collectZoneSelectData.findIndex(
          data => data.zone === item.zone || item.name
        )
        if (_findIndex === -1) {
          collectZoneSelectData.push({
            zone: item.zone || item.name,
            domain: item?.storage_domains?.[0] || '',
            sc_id: item?.sc_id || ''
          })
        } else {
          // 判断当前收集的区域是否有sc_id
          !collectZoneSelectData[_findIndex].sc_id &&
            (collectZoneSelectData[_findIndex].sc_id = item?.sc_id || '')
        }
      })
      runInAction(() => {
        this.update({
          ...data,
          zoneSelectData: collectZoneSelectData
        })
      })
    } catch (error) {
      message.error(error)
    }
  }

  @action
  update({ ...props }) {
    Object.assign(this, props)
  }
}
