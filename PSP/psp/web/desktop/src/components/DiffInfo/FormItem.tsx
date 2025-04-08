import * as React from 'react'
import { observer } from 'mobx-react'
import styled from 'styled-components'

import * as FormItems from '@/components/FormField'
import { FieldType } from '@/domain/Applications/App/Field'

const Wrapper = styled.div`
  .item {
    width: 700px;
    position: relative;

    &.hidden {
      position: absolute;
      height: 0;
      width: 0;
      opacity: 0;
    }

    .drag-icon,
    .right-icons {
      display: none;
    }

    &:hover {
      background: #f5f9ff;
      cursor: pointer;

      .drag-icon {
        display: block;
        position: absolute;
        left: 10px;
        top: 10px;
        z-index: 99;
        cursor: move;
      }

      .right-icons {
        display: flex;
        align-items: center;
        position: absolute;
        top: 20px;
        right: 80px;

        .icon {
          margin-right: 20px;
        }
      }
    }
  }
`

interface IProps {
  model: any
}

@observer
export default class FormItem extends React.Component<IProps> {
  render() {
    const { model } = this.props
    const Component =
      FormItems[
        {
          [FieldType.text]: 'Input',
          [FieldType.list]: 'Select',
          [FieldType.multiple]: 'MultiSelect',
          [FieldType.checkbox]: 'Checkbox',
          [FieldType.radio]: 'Radio',
          [FieldType.lsfile]: 'Uploader',
          [FieldType.label]: 'Label',
          [FieldType.date]: 'Date',
          [FieldType.node_selector]: 'NodeSelector'
        }[model.type]
      ] || FormItems['Label']

    return (
      <Wrapper>
        <div className={`item ${model.hidden ? 'hidden' : ''}`}>
          <Component model={model} formModel={{}} showId={true} />
        </div>
      </Wrapper>
    )
  }
}
