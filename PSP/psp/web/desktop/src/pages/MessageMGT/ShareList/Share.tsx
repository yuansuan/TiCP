/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React from 'react'
import styled from 'styled-components'
import { Card } from 'antd'
import { computed } from 'mobx'
import {Http} from '@/utils'
import { Record as DomainRecord } from '@/domain/RecordList/record'
import { Button, Modal } from '@/components'
import { env } from '@/domain'
import { ShareFileContent } from '@/components/HeaderToolbar/Message/Share'

interface IProps {
  item: DomainRecord
  readShare: () => Promise<any>
}

const StyledLayout = styled.div`
  margin: 20px;

  .cardBody {
    display: flex;

    .right {
      margin-left: auto;

      > button {
        margin: 0 4px;
      }

      .isRead {
        color: #ccc;
      }
      .view {
        margin-left:5px;
        color: #3182FF;
        cursor: pointer
      }
      .notRead {
        cursor: pointer;
        color: #ff4d4f;
      }
    }
  }
`

export class Share extends React.Component<IProps> {
   async getShareItem (id) {
    const {data} =await Http.get('/storage/share/get',{
      params: {
        id
      }
    })
    const newTitle = data?.isdir 
      ? `[${data.name}]文件夹` : `[${data.name}]文件`
    await Modal.show({
      title: `${newTitle} 保存方式`,
      footer:null,
      content: ({onOk,onCancel}) => (
        <ShareFileContent onOk={onOk} onCancel={onCancel} {...data}/>
      )
    })
  }

  render() {
    const { item,readShare } = this.props

    return (
      <StyledLayout>
        <Card title={'分享通知'} extra={item.timeTitle}>
          <div className='cardBody'>
            <div className='content'>{item.content}</div>
             <div className='right'> 
              {item.state === 2 ? (
              <span className='isRead'>已读</span>
            ) : (
              <span className='notRead' onClick={() => readShare()}>
                标为已读
              </span>
            )}
              <span className='view' onClick={() => this.getShareItem(item.id)}>保存</span>
            </div>
          </div>
        </Card>
      </StyledLayout>
    )
  }
}
