/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Progress } from 'antd'
import { Icon } from '@/components'
import { formatByte } from '@/utils/Validator'
import { observer } from 'mobx-react-lite'

const StyledLayout = styled.div`
  display: flex;
  width: 100%;

  > .icon {
    width: 44px;
    height: 44px;
    background-color: ${({ theme }) => theme.backgroundColorBase};
    display: flex;
    justify-content: center;
    align-items: center;

    .anticon {
      font-size: 32px;
    }
  }

  > .main {
    display: flex;
    flex-direction: column;
    margin: 0 8px;
    flex: 1;

    > .description {
      display: flex;

      > .left {
        font-size: ${({ theme }) => theme.fontSizeBody};
        color: rgba(0, 0, 0, 0.6);
        line-height: 20px;
        max-width: 220px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }

      > .right {
        margin-left: auto;
        display: flex;
        align-items: center;

        &.uploading,
        &.paused {
          color: #40a9ff;
        }

        &.error {
          color: #f5222d;
        }

        &.done,
        &.success {
          color: #9b9b9b;
        }

        .anticon {
          font-size: 24px;
        }
      }
    }

    > .progress {
      line-height: 0;

      .ant-progress {
        line-height: 0;

        > .ant-progress-outer {
          line-height: 0;
        }
      }
    }

    > .info {
      display: flex;
      font-size: ${({ theme }) => theme.fontSizeCaption};
      color: rgba(0, 0, 0, 0.45);

      > .right {
        margin-left: auto;
      }
    }
  }

  > .operators {
    display: flex;
    align-items: center;

    > * {
      margin-left: 8px;

      &:first-child {
        margin-left: 0;
      }
    }
  }
`

export type Point = 'sc' | 'box' | 'local'

type Props = {
  context: {
    name: string
    speed?: number
    percent?: number
    status?: string
    loaded?: number
    size?: number
  }
  direction: [Point, Point]
  isFile: boolean
  operators?: React.ReactNode[]
}

function getProcessStatus(status: string) {
  switch (status) {
    case 'uploading':
      return 'active'
    case 'success':
    case 'done':
      return 'success'
    case 'removed':
    case 'error':
      return 'exception'
    default:
      return 'normal'
  }
}

function mapDirectionIcon(type: Point) {
  return {
    sc: <Icon type='data_transport_active' />,
    box: <Icon type='storage_active' />,
    local: <Icon type='folder_hover' />
  }[type]
}

export const UploadProgress = observer(
  ({ isFile = true, context, direction, operators }: Partial<Props>) => {
    const { speed, name, percent, status, size, loaded } = context

    const isFromSc = direction[0] === 'sc'

    return (
      <StyledLayout>
        <div className='icon'>
          {isFile && <Icon type='file_table' />}
          {!isFile && <Icon type='folder-default' />}
        </div>
        <div className='main'>
          <div className='description'>
            <div className='left' title={name}>
              {name}
            </div>
            <div className={`right ${status}`}>
              {mapDirectionIcon(direction[isFromSc ? 1 : 0])}
              {isFromSc ? '<----' : '---->'}
              {mapDirectionIcon(direction[isFromSc ? 0 : 1])}
            </div>
          </div>
          <div className='progress'>
            <Progress
              status={getProcessStatus(status)}
              percent={parseInt(percent.toString())}
              strokeWidth={6}
              showInfo={true}
            />
          </div>
          <div className='info'>
            {speed !== undefined && (
              <div className='left'>
                {status === 'uploading' && `${formatByte(speed)}/s`}
              </div>
            )}
            <div className='right'>
              {status === 'uploading' &&
                `${formatByte(loaded)} / ${formatByte(size)}`}
              {status !== 'uploading' && formatByte(size)}
            </div>
          </div>
        </div>
        <div className='operators'>{operators}</div>
      </StyledLayout>
    )
  }
)
