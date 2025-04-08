/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import styled from 'styled-components'

import { FieldType } from '@/domain/Applications/App/Field'
import { DndTypes } from '@/components/FormField'
import Field from './Field'

const Wrapper = styled.div`
  width: 276px;
  display: flex;
  flex-direction: column;
  min-height: 100%;
  border-right: 1px solid #d9d9d9;
  padding: 0 8px 0;

  > .title {
    margin: 8px 0 10px 0;
    padding: 0 10px;

    .icon {
      font-size: 12px;
      margin-right: 5px;
    }

    .tip {
      font-family: 'PingFangSC-Medium';
      font-size: 16px;
      color: #595959;
    }

    .subTip {
      font-family: 'PingFangSC-Regular';
      font-size: 14px;
      color: rgba(0, 0, 0, 0.45);
    }
  }

  > .list {
    padding-bottom: 20px;
    overflow: auto;
  }
`
export default class Left extends React.Component<any> {
  private fields = [
    {
      type: 'Section',
      name: 'Section',
      dndType: DndTypes.EMPTY_SECTION
    },
    {
      type: FieldType.text,
      name: 'Input',
      dndType: DndTypes.EMPTY_FORM_ITEM
    },
    {
      type: FieldType.list,
      name: 'Select',
      dndType: DndTypes.EMPTY_FORM_ITEM
    },
    {
      type: FieldType.multiple,
      name: 'Multiple Select',
      dndType: DndTypes.EMPTY_FORM_ITEM
    },
    {
      type: FieldType.radio,
      name: 'Radio',
      dndType: DndTypes.EMPTY_FORM_ITEM
    },
    {
      type: FieldType.checkbox,
      name: 'Checkbox',
      dndType: DndTypes.EMPTY_FORM_ITEM
    },
    {
      type: FieldType.label,
      name: 'Label',
      dndType: DndTypes.EMPTY_FORM_ITEM
    },
    {
      type: FieldType.date,
      name: 'Date',
      dndType: DndTypes.EMPTY_FORM_ITEM
    },
    // {
    //   type: FieldType.node_selector,
    //   name: '关联核数',
    //   dndType: DndTypes.EMPTY_FORM_ITEM
    // }
    {
      type: FieldType.cascade_selector,
      name: 'Cascade Select',
      dndType: DndTypes.EMPTY_FORM_ITEM
    }
  ]

  render() {
    return (
      <Wrapper>
        <div className='title'>
          <div className='tip'>控件组</div>
          <div className='subTip'>通过拖拽组件至右侧面板进行模版配置</div>
        </div>

        <div className='list'>
          {this.fields.map(field => (
            <Field key={field.type} {...field} />
          ))}
        </div>
      </Wrapper>
    )
  }
}
