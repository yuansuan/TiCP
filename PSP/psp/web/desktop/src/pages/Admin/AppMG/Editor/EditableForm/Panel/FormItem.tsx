import * as React from 'react'
import { action, computed } from 'mobx'
import { observer } from 'mobx-react'
import styled from 'styled-components'
import { DragSource, DragSourceConnector } from 'react-dnd'
import flow from 'lodash/flow'
import { message } from 'antd'
import { EditOutlined, DeleteOutlined } from '@ant-design/icons'
import { Icon, Modal } from '@/components'
import { inject } from '@/pages/context'
import * as FormItems from '@/components/FormField'
import { FieldType } from '@/domain/Applications/App/Field'
import { DndTypes } from '@/components/FormField'
import ItemTarget from './ItemTarget'

const Wrapper = styled.div`
  .item {
    width: 700px;
    position: relative;

    .drag-icon,
    .right-icons {
      display: none;
    }

    &:hover {
      background-color: #cdddf7;
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

        .anticon {
          margin-right: 20px;
        }
      }
    }
  }
`

interface IProps {
  connectDragPreview?: any
  connectDragSource?: any
  app?: any
  index: number
  parentIndex: number
  model: any
  formModel: any
  fetchUploadPath?: string
  onDelete: (id: number) => void
}

const source = {
  beginDrag(props) {
    return { ...props }
  }
}

const sourceCollect = (connect: DragSourceConnector) => ({
  connectDragSource: connect.dragSource(),
  connectDragPreview: connect.dragPreview()
})

@flow(
  observer,
  DragSource(DndTypes.FORM_ITEM, source, sourceCollect),
  inject(({ app, fetchUploadPath }) => ({
    app,
    fetchUploadPath
  }))
)
export default class FormItem extends React.Component<IProps> {
  @computed
  get fieldModel() {
    const { model, formModel } = this.props
    return formModel[model.id]
  }

  private onDelete = () => {
    const { onDelete, index, formModel, model } = this.props

    Modal.showConfirm({
      content: '确认要删除此控件吗？'
    }).then(() => {
      onDelete(index)
      //  删除掉
      Reflect.deleteProperty(formModel, model.id)
    })
  }

  private onEdit = () => {
    const { model } = this.props

    if (
      this.fieldModel &&
      (this.fieldModel.value ||
        (this.fieldModel.values && this.fieldModel.values.length > 0))
    ) {
      Modal.showConfirm({
        content: '编辑控件会导致当前数据丢失，确认编辑吗？'
      }).then(() => model.updateEditing(true))
    } else {
      model.updateEditing(true)
    }
  }

  // cancel edit
  @action
  private onCancel = viewModel => {
    const { model, index, onDelete } = this.props

    // cancel create
    if (!viewModel.id) {
      onDelete(index)
      return
    }

    viewModel.reset()
    model.updateEditing(false)
  }

  // confirm edit
  @action
  private onConfirm = viewModel => {
    const { model, formModel, app } = this.props
    const newId = viewModel.id

    if (!newId) {
      message.error('ID 不能为空')
      return
    }

    if (!viewModel.label) {
      message.error('字段不能为空')
      return
    }

    if (viewModel.isMasterSlave) {
      if (!viewModel.masterIncludeKeywords) {
        message.error('从文件关键字不能为空')
        return
      }
      if (!viewModel.masterIncludeExtensions) {
        message.error('从文件后缀名不能为空')
        return
      }
    }

    // validate id
    if (newId !== model.id) {
      // check id duplicate
      if (app.fieldIds.includes(newId)) {
        message.error(`ID:${newId} 已存在`)
        return
      }
    }

    // delete old id from formModel
    if (newId !== model.id) {
      Reflect.deleteProperty(formModel, model.id)
    }

    viewModel.submit()
    this.props.model.updateEditing(false)
  }

  render() {
    const {
      index,
      parentIndex,
      model,
      model: { editing },
      formModel,
      connectDragPreview,
      connectDragSource,
      app
    } = this.props

    if (FieldType.lsfile === model.type) console.log(model)

    const Item =
      FormItems[
        {
          [FieldType.text]: 'Input',
          [FieldType.list]: 'Select',
          [FieldType.multiple]: 'MultiSelect',
          [FieldType.checkbox]: 'Checkbox',
          [FieldType.radio]: 'Radio',
          [FieldType.lsfile]: 'Uploader',
          [FieldType.lsfile_yscloud]: 'UploaderYSCloud',
          [FieldType.label]: 'Label',
          [FieldType.date]: 'Date',
          [FieldType.node_selector]: 'NodeSelector',
          [FieldType.cascade_selector]: 'CascadeSelector'
        }[model.type]
      ] || FormItems['Label']
    const isUploader = model.type === FieldType.lsfile

    return (
      <Wrapper>
        <ItemTarget index={index} parentIndex={parentIndex} />

        {connectDragPreview(
          <div className='item'>
            {!editing &&
              connectDragSource(
                <span>
                  <Icon className='drag-icon' type='drag' />
                </span>
              )}

            {editing ? (
              <Item.Editor
                model={model}
                formModel={formModel}
                appId={app.appId}
                onCancel={this.onCancel}
                onConfirm={this.onConfirm}
              />
            ) : (
              <Item
                model={model}
                formModel={formModel}
                appId={app.appId}
                showId={true}
                {...(isUploader
                  ? { fetchUploadPath: this.props.fetchUploadPath }
                  : {})}
              />
            )}

            {!editing && (
              <div className='right-icons'>
                <EditOutlined onClick={this.onEdit} />
                <DeleteOutlined onClick={this.onDelete} />
              </div>
            )}
          </div>
        )}
      </Wrapper>
    )
  }
}
