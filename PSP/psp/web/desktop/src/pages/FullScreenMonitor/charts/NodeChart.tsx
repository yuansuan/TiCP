import React from 'react'
import { useStore } from '../store'
import { observer } from 'mobx-react-lite'
import { Carousel } from 'antd'
import styled from 'styled-components'
import { roundNumber } from '@/utils/formatter'
import { PRECISION } from '@/domain/common'
import { Icon } from '@/components'

function chunkArray<T>(array: T[], size: number): T[][] {
  const chunks: T[][] = []
  for (let i = 0; i < array.length; i += size) {
    chunks.push(array.slice(i, i + size))
  }
  return chunks
}

const toNodeAttr = (val: any, num: number) => {
  if (val === -1) return '--'
  if (typeof val === 'string') return roundNumber(val + 0, num)
  if (typeof val === 'number') return roundNumber(val, num)
  return val
}

const Node = ({ options }) => {
  const okColor = '#98c46d'
  const notOkColor = '#999'

  let memUsage =
    options['max_mem'] === 0
      ? 0
      : ((options['max_mem'] - options['available_mem']) / options['max_mem']) *
        100

  let color = options['node_status'] === 'Up' ? okColor : notOkColor

  const memUsageStr =
    options['max_mem'] === 0
      ? '0%'
      : isNaN(memUsage)
      ? '--%'
      : `${toNodeAttr(memUsage, PRECISION)}%`

  const utStr =
    toNodeAttr(options['cpu_percent'], 0) !== '--'
      ? `${toNodeAttr(options['cpu_percent'], PRECISION)}%`
      : `--%`

  return (
    <div className='nodeDetail'>
      <div className='content'>
        <Icon
          type='zhuji'
          style={{ color: color, fontSize: '70px', paddingRight: '8px' }}
        />
        <div className='otherInfos'>
          <p className='nodeName'>{options['node_name']}</p>

          <p>
            <span className='infoLabel'>CPU利用率: {utStr}</span>
          </p>
          <p>
            <span className='infoLabel'>内存利用率: {memUsageStr}</span>
          </p>
        </div>
      </div>
    </div>
  )
}

const NodeChart = observer(() => {
  const store = useStore()
  const { nodeList } = store
  const chunkedList = chunkArray(nodeList, 9)
  // const chunkedList = chunkArray(nodeTestList, 9)

  return (
    <div
      style={{
        width: '100%',
        height: '100%',
        position: 'relative'
      }}>
      <div
        style={{
          position: 'absolute',
          top: 0,
          bottom: 0,
          right: 0,
          left: 0
        }}>
        <NodeWrapper>
          <Carousel autoplay dotPosition='bottom' autoplaySpeed={5000}>
            {chunkedList?.length !== 0 ? (
              chunkedList?.map((item, index) => (
                <div key={index} className='nodeInfo'>
                  <ul>
                    {item?.map((element, index) => (
                      <li key={index}>
                        <Node key={element.node_name} options={element} />
                      </li>
                    ))}
                  </ul>
                </div>
              ))
            ) : (
              <div className='nodeInfo'>无节点信息</div>
            )}
          </Carousel>
        </NodeWrapper>
      </div>
    </div>
  )
})
export default NodeChart

const NodeWrapper = styled.div`
  .nodeInfo {
    li {
      float: left;
      list-style: none;
    }

    .nodeDetail {
      display: flex;
      justify-content: center;
      width: 270px;

      .ant-popover-inner {
        padding: 0 20px;
        .ant-popover-title {
          padding: 5px 0;
        }

        .ant-popover-inner-content {
          padding: 12px 0;
        }
      }

      .content {
        float: left;
        color: rgba(255, 255, 255, 0.6);
        padding-top: 0;
        .nodeName {
          font-weight: 700;
          font-style: normal;
          font-size: 15px;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: wrap;
          margin-bottom: 0;
        }

        .otherInfos {
          margin-bottom: 10px;
          > p {
            margin: 0;
            height: 22px;
            .infoLabel {
              font-size: 12px;
            }
          }
        }
      }
    }
  }

  .ant-carousel .slick-dots {
    margin-top: 100px;
    margin-right: 0;
    margin-left: 0;
  }

  .ant-carousel .slick-dots li button::before {
    display: none;
  }

  .ant-carousel .slick-dots-bottom {
    bottom: -12px;
  }
`
