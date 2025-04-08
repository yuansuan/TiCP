/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Card } from 'antd'
import { Invitation as DomainInvitation } from '@/domain/InvitationList/Invitation'
import { useObserver } from 'mobx-react-lite'
import { formatPhone } from '@ys/utils'

const StyledLayout = styled.div`
  margin: 20px 20px 0px 20px;

  .cardBody {
    display: flex;

    .right {
      margin-left: auto;

      > button {
        margin: 0 4px;
      }

      .waiting {
        color: '#cccccc';
      }

      .accept {
        color: #52c41a;
      }

      .reject {
        color: #ff4d4f;
      }
    }
  }
`

interface IProps {
  item: DomainInvitation
}

export function Invitation({ item }: IProps) {
  return useObserver(() => (
    <StyledLayout>
      <Card title='成员邀请' extra={item.create_time.toString()}>
        <div className='cardBody'>
          <div className='content'>{`${item.create_name} 邀请 ${
            item.real_name || formatPhone(item.phone)
          } 加入${item.company_name}`}</div>
          <div className='right'>
            {item.status === 1 && <span className='waiting'>未处理</span>}
            {item.status === 2 && <span className='accept'>已确认</span>}
            {item.status === 3 && <span className='reject'>已拒绝</span>}
          </div>
        </div>
      </Card>
    </StyledLayout>
  ))
}
