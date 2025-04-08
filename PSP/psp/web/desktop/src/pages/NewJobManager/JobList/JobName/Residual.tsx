/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'
import styled from 'styled-components'
import { Button } from '@/components'
import { observer, useLocalStore } from 'mobx-react-lite'
import { Chart, Geom, Axis, Tooltip, Legend } from 'bizcharts'
import { Radio } from 'antd'
import { ReloadOutlined } from '@ant-design/icons'
import { jobServer, standardJobMGTServer } from '@/server'

const StyledLayout = styled.div`
  height: 600px;

  .bar {
    height: 40px;
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 20px;

    .ant-radio-group {
      margin-left: 10px;
    }
  }
`

type Props = {
  id: string
  isStandardJob?: boolean
}

export const Residual = observer(function Residual({ id, isStandardJob }: Props) {
  const state = useLocalStore(() => ({
    data: [],
    setData(data) {
      state.data = data
    },
    // x轴可选值
    xs: [],
    setXs(xs) {
      state.xs = xs
    },
    // 当前x轴
    x: '',
    setX(x) {
      state.x = x
    },
    get chartData() {
      if (!state.data.length) {
        return []
      }

      const currentXVals = state.data.find(v => v.name === state.x).values
      const yVals = state.data
        .filter(n => !state.xs.includes(n.name))
        .sort((y1, y2) => y1.name.localeCompare(y2.name))
      const data = []
      yVals.forEach(y => {
        y.values.forEach((val: number, index: number) => {
          if (val < 0) return
          data.push({
            x: currentXVals[index],
            yName: y.name,
            val
          })
        })
      })
      return data
    }
  }))

  useEffect(() => {
    refresh()
  }, [])

  async function refresh() {
    const {
      data: { available_xvar, vars }
    } = await (isStandardJob ? standardJobMGTServer.getResidualData(id) : jobServer.getResidualData(id))
    state.setXs(available_xvar.filter(x => vars.find(v => v.name === x)))
    state.setX(state.xs[0])
    state.setData(vars)
  }

  const cols = {
    x: {
      range: [0, 1],
      alias: state.x
    },
    val: {
      type: 'log',
      base: 10,
      nice: true,
      formatter: value => {
        return parseFloat(value.toPrecision(12))
      }
    }
  }

  return (
    <StyledLayout>
      <div className='bar'>
        {state.xs.length > 1 ? (
          <div>
            横坐标选择:
            <Radio.Group
              size='small'
              value={state.x}
              buttonStyle='solid'
              onChange={e => (state.x = e.target.value)}>
              {state.xs.map(x => (
                <Radio.Button key={x} value={x}>
                  {x}
                </Radio.Button>
              ))}
            </Radio.Group>
          </div>
        ) : (
          <div />
        )}
        <Button
          icon={<ReloadOutlined />}
          type='primary'
          size='small'
          ghost
          onClick={refresh}>
          刷新
        </Button>
      </div>

      <Chart
        autoFit
        height={500}
        data={state.chartData}
        scale={cols}
        padding='auto'
        forceFit>
        <Legend marker={{ symbol: 'square' }} position='bottom' />
        <Axis name='x' title />
        <Axis name='val' />
        <Geom type='line' position='x*val' color={'yName'} shape={'smooth'} />
        <Tooltip
          crosshairs={{
            type: 'y'
          }}
        />
      </Chart>
    </StyledLayout>
  )
})
