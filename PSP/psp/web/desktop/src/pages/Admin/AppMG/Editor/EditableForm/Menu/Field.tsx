import * as React from 'react'
import { DragSource, DragSourceMonitor, DragSourceConnector } from 'react-dnd'
import flow from 'lodash/flow'
import styled from 'styled-components'
import { Radio, Checkbox, DatePicker } from 'antd'

import { Icon } from '@/components'
import { FieldType } from '@/domain/Applications/App/Field'

const Wrapper = styled.div`
  display: flex;
  flex-direction: column;
  justify-content: center;
  min-height: 50px;
  cursor: pointer;
  padding: 10px;
  border: 1px solid transparent;

  &:hover {
    background: #f0f5fd;
    border-color: #10398b;
    border-radius: 2px;
    cursor: move;

    .name {
      color: ${props => props.theme.primaryColor};
    }
  }

  .name {
    font-family: 'PingFangSC-Medium';
    font-size: 14px;
    color: rgba(0, 0, 0, 0.85);
    margin: 4px 0;
  }

  .sample {
    display: flex;
    align-items: center;
    font-size: 12px;

    .inputSample,
    .selectSample,
    .uploadSample {
      width: 180px;
      height: 30px;
      line-height: 30px;
      background: #ffffff;
      border: 1px solid #d9d9d9;
      border-radius: 2px;
      color: #cccccc;
      padding: 0 5px;
    }

    .selectSample {
      display: flex;

      .icon {
        font-size: 12px;
        margin: auto;
        margin-right: 0;
      }
    }

    .uploadSample {
      width: 75px;
      padding: 0 4px;
    }

    .radioSample {
      .ant-radio-disabled {
        .ant-radio-inner {
          cursor: inherit;
        }
        + span {
          color: #6e6f72;
          font-size: 12px;
          cursor: inherit;
        }
      }
    }

    .checkboxSample {
      .ant-checkbox-disabled.ant-checkbox-checked .ant-checkbox-inner::after {
        border-color: #d9d9d9;
      }
      .ant-checkbox-disabled {
        cursor: inherit;
      }
      .ant-checkbox {
        .ant-checkbox-input {
          cursor: inherit;
        }

        + span {
          color: #6e6f72;
          font-size: 12px;
          cursor: inherit;
        }
      }
    }

    .dateSample {
      .ant-calendar-picker {
        cursor: move;
      }

      .ant-calendar-picker-icon {
        color: #d9d9d9;
      }

      .ant-input[disabled] {
        cursor: inherit;
      }
    }
  }
`

const fieldSource = {
  beginDrag(props) {
    return props
  }
}
const collect = (connect: DragSourceConnector, monitor: DragSourceMonitor) => ({
  connectDragSource: connect.dragSource(),
  isDragging: monitor.isDragging()
})

@flow(DragSource(props => props.dndType, fieldSource, collect))
export default class Field extends React.Component<any> {
  renderSample = type => {
    switch (type) {
      case 'Section':
        return null
      case FieldType.text:
        return (
          <div className='sample'>
            <span>字段：</span>
            <div className='inputSample'>请输入</div>
          </div>
        )
      case FieldType.list:
        return (
          <div className='sample'>
            <span>字段：</span>
            <div className='selectSample'>
              <span>请选择</span>
              <Icon className='icon' type='down' />
            </div>
          </div>
        )
      case FieldType.multiple:
        return (
          <div className='sample'>
            <span>字段：</span>
            <div className='selectSample'>
              <span>请选择</span>
              <Icon className='icon' type='down' />
            </div>
          </div>
        )
      // case FieldType.lsfile:
      //   return (
      //     <div className='sample'>
      //       <span>字段：</span>
      //       <div className='uploadSample'>
      //         <Icon className='icon' type='upload' />
      //         <span>上传文件</span>
      //       </div>
      //     </div>
      //   )
      case FieldType.radio:
        return (
          <div className='sample'>
            <div className='radioSample'>
              <Radio checked disabled={true}>
                单选文本
              </Radio>
            </div>
          </div>
        )
      case FieldType.checkbox:
        return (
          <div className='sample'>
            <div className='checkboxSample'>
              <Checkbox checked disabled={true}>
                多选文本
              </Checkbox>
            </div>
          </div>
        )
      case FieldType.label:
        return (
          <div className='sample'>
            <div className='labelSample'>标签文字</div>
          </div>
        )
      case FieldType.date:
        return (
          <div className='sample'>
            <div className='dateSample'>
              <DatePicker disabled={true} />
            </div>
          </div>
        )
      case FieldType.node_selector:
        return (
          <div className='sample'>
            <span>字段：</span>
            <div className='selectSample'>
              <span>请选择节点</span>
              <Icon className='icon' type='down' />
            </div>
          </div>
        )
      case FieldType.cascade_selector:
        return (
          <div className='sample'>
            <span>字段：</span>
            <div className='selectSample'>
              <span>关联选择器</span>
              <Icon className='icon' type='down' />
            </div>
          </div>
        )  
      default:
        return null
    }
  }

  render() {
    const { name, type, isDragging, connectDragSource, className } = this.props
    const opacity = isDragging ? 0.4 : 1

    return connectDragSource(
      <div style={{ opacity }} className={className}>
        <Wrapper>
          <div className='name'>{name}</div>
          {this.renderSample(type)}
        </Wrapper>
      </div>
    )
  }
}
