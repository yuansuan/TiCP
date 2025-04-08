import React, { useEffect } from 'react'
import { observer } from 'mobx-react-lite'
import { Spin, Radio } from 'antd'
import { ComboCardWrapper, Content, Layout, Sider, SiderTitle } from './styles'
import { CHARGE_TYPE } from '@/constant'
import { formatUnixTime } from '@/utils'

interface IProps {
  loading?: boolean
  comboList: any[]
  billModel: number
  comboTicketId?: string
  comboSoftwareId?: string
  setComboTicketId?: (id: string) => void
  setComboSoftwareId?: (id: string) => void
}

export const ComboSelector = observer(function ComboSelector({
  loading,
  billModel,
  comboList,
  comboTicketId,
  comboSoftwareId,
  setComboTicketId,
  setComboSoftwareId
}: IProps) {
  useEffect(() => {
    if (comboTicketId) {
      document.querySelector('.validate_combo_tip').innerHTML = ''
    } else {
      document.querySelector('.validate_combo_tip').innerHTML = '请选择套餐'
    }
    const combo = comboList.find(item => item.ticket_id === comboTicketId)
    setComboSoftwareId(combo?.softwares[0]?.id)
  }, [comboTicketId])

  const changeCombo = event => {
    setComboTicketId(event.target.value)
    setComboSoftwareId('')
  }

  const showComboList = comboList?.filter(combo => {
    if (billModel === 2) {
      return combo.chargeType === CHARGE_TYPE.MONTHLY_TYPE
    }
    if (billModel === 3) {
      return combo.chargeType === CHARGE_TYPE.HOURLY_TYPE && !combo.is_free
    }
    if (billModel === 4) {
      return combo.chargeType === CHARGE_TYPE.HOURLY_TYPE && combo.is_free
    }
  })

  return (
    <>
      <Layout>
        <Sider>
          <SiderTitle>选择套餐</SiderTitle>
        </Sider>
        <Content chargeType={billModel}>
          <Spin spinning={loading}>
            <ComboCardWrapper style={{ padding: 0, marginTop: '10px' }}>
              <Radio.Group
                buttonStyle='solid'
                value={comboTicketId}
                onChange={e => changeCombo(e)}>
                <div className='shop-checkblock-wrap'>
                  {showComboList?.map(list => {
                    return (
                      <Radio.Button
                        value={list.ticket_id}
                        className='shop-checkblock-item'
                        key={list.ticket_id}>
                        <div className='main'>
                          <h3>{list.combo_name}</h3>
                          <div className='configure'>
                            <span className='title'>镜像名称</span>
                            <select
                              style={{
                                width: '100%',
                                borderColor: '#d9d9d9',
                                borderRadius: '2px',
                                height: '32px',
                                padding: '0 10px'
                              }}
                              value={
                                list.ticket_id === comboTicketId
                                  ? comboSoftwareId
                                  : list?.softwares[0]?.id
                              }
                              placeholder='请选择'
                              onChange={e => {
                                setComboSoftwareId(e.target.value)
                                setComboTicketId(list.ticket_id)
                              }}>
                              {list?.softwares?.map(item => (
                                <option key={item.id} value={item.id}>
                                  {item.name}
                                </option>
                              ))}
                            </select>

                            <span className='desc'>
                              {list.ticket_id === comboTicketId
                                ? list?.softwares?.find(
                                    item => item.id === comboSoftwareId
                                  )?.desc
                                : list?.softwares[0]?.desc}
                            </span>
                          </div>
                          <div className='configure'>
                            <span className='title'>实例名称</span>
                            <div className='info'>
                              <div>{list?.hardwares[0]?.name}</div>
                              <div className='model'>
                                {list?.hardwares[0]?.desc}
                              </div>
                            </div>
                          </div>

                          <div className='time'>
                            {list?.chargeType === CHARGE_TYPE.MONTHLY_TYPE ? (
                              <p>
                                有效期：
                                {formatUnixTime(
                                  list?.valid_begin_time?.seconds
                                ) || '--'}{' '}
                                至<br />
                                {formatUnixTime(
                                  list?.valid_end_time?.seconds
                                ) || '--'}
                              </p>
                            ) : list?.chargeType === CHARGE_TYPE.HOURLY_TYPE ? (
                              <p>
                                激活时间：
                                {formatUnixTime(
                                  list?.valid_begin_time?.seconds
                                ) || '--'}
                              </p>
                            ) : (
                              ''
                            )}
                          </div>
                        </div>
                        {list?.chargeType === CHARGE_TYPE.HOURLY_TYPE && (
                          <div className='usage'>
                            <div>
                              已使用：
                              {Number.isInteger(list?.used_time / 3600)
                                ? list?.used_time / 3600
                                : (list?.used_time / 3600).toFixed(2)}
                              小时
                            </div>
                            <div>
                              未使用：
                              {Number.isInteger(list?.remain_time / 3600)
                                ? list?.remain_time / 3600
                                : (list?.remain_time / 3600).toFixed(2)}
                              小时
                            </div>
                          </div>
                        )}
                      </Radio.Button>
                    )
                  })}
                </div>
              </Radio.Group>
            </ComboCardWrapper>
          </Spin>
          {showComboList && (
            <div className='validate_tip validate_combo_tip'>
              {'请选择套餐'}
            </div>
          )}
        </Content>
      </Layout>
      <Layout>
        <Sider>
          <SiderTitle>套餐定价</SiderTitle>
        </Sider>
        <Content>--</Content>
      </Layout>
    </>
  )
})

export default ComboSelector
