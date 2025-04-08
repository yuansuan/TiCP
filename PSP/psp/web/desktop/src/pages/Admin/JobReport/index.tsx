import React, { useEffect, useState } from 'react'
import { Tabs } from 'antd'
import { observer } from 'mobx-react'
import styled from 'styled-components'
import { Http } from '@/utils'

import ReportBySoftware from './ReportBySoftware'
import ReportByUser from './ReportByUser'

const Wrapper = styled.div`
  padding: 20px;
  padding-top:0px;
  background: #fff;
  
`

const { TabPane } = Tabs

export default observer(function JobReport() {

  const [computeTypes, setComputeTypes] = useState([])
  const [computeTypesMap, setComputeTypesMap] = useState({})

  useEffect(() => {
    (async () => {
      const res = await Http.get('/job/computeTypes')
      setComputeTypes(res?.data?.compute_types || [])
      setComputeTypesMap(res?.data?.compute_types?.reduce((pre, curr) => (pre[curr.compute_type] = curr.show_name, pre), {}))
    })()
  }, [])

  return (
    <Wrapper>
      <Tabs>
        <TabPane tab="应用" key={'app'}>
          <ReportBySoftware computeTypes={computeTypes} computeTypesMap={computeTypesMap} />
        </TabPane>
        <TabPane tab="用户" key={'user'}>
          <ReportByUser computeTypes={computeTypes} computeTypesMap={computeTypesMap} />
        </TabPane>
      </Tabs>
    </Wrapper>
  )
})