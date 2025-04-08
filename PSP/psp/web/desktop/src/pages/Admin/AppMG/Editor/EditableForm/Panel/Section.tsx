import * as React from 'react'
import { observable, action } from 'mobx'
import { observer } from 'mobx-react'
import {
  DragSource,
  DragSourceConnector,
  DragSourceMonitor,
  DropTarget,
  DropTargetMonitor,
  DropTargetConnector
} from 'react-dnd'
import { EditOutlined, DeleteOutlined } from '@ant-design/icons'
import flow from 'lodash/flow'
import { Input, message } from 'antd'

import { Icon, Modal } from '@/components'
import Section from '@/domain/Applications/App/Section'
import { App } from '@/domain/Applications'
import { inject } from '@/pages/context'
import { DndTypes } from '@/components/FormField'
import ItemTarget from './ItemTarget'
import SecionTarget from './SectionTarget'
import FormItem from './FormItem'
import { Wrapper } from './SectionStyle'

const source = {
  beginDrag(props) {
    return { ...props }
  }
}

const sourceCollect = (
  connect: DragSourceConnector,
  monitor: DragSourceMonitor
) => ({
  connectDragSource: connect.dragSource(),
  connectDragPreview: connect.dragPreview(),
  isDragging: monitor.isDragging()
})

const target = {
  drop: (props, monitor: DropTargetMonitor) => {
    const item = monitor.getItem()
    const type = monitor.getItemType()
    const {
      index,
      app: { subForm }
    } = props

    // toggle existent section
    if (type === DndTypes.SECTION) {
      subForm.toggle(item.index, index)
    } else {
      // new section
      subForm.add(new Section(), index)
    }
  }
}

const targetCollect = (
  connect: DropTargetConnector,
  monitor: DropTargetMonitor
) => ({
  connectDropTarget: connect.dropTarget(),
  isOver: monitor.isOver(),
  canDrop: monitor.canDrop()
})

interface IProps {
  app?: App
  canDrop?: any
  isDragging?: any
  isOver?: any
  connectDragSource?: any
  connectDragPreview?: any
  connectDropTarget?: any
  onDelete: (index: number) => void
  index: number
  model: any
  formModel: any
}

@flow(
  observer,
  DragSource(DndTypes.SECTION, source, sourceCollect),
  DropTarget([DndTypes.SECTION, DndTypes.EMPTY_SECTION], target, targetCollect),
  inject(({ app }) => ({
    app
  }))
)
export default class SectionUI extends React.Component<IProps> {
  @observable name
  @action
  updateName = name => (this.name = name)

  inputRef = null

  constructor(props) {
    super(props)

    this.updateName(props.model.name)
  }

  deleteField = index => this.props.model.delete(index)

  private onKeyDown = e => {
    if (e.keyCode === 13) {
      this.onConfirm()
    } else if (e.keyCode === 27) {
      this.onCancel()
    }
  }

  private onConfirm = (e?) => {
    e && e.preventDefault()

    const { model } = this.props

    if (!this.name && model.editing) {
      message.error('请输入 section 名称')
      this.inputRef && this.inputRef.focus()
      return
    }

    model.updateName(this.name)
    model.updateEditing(false)
  }

  private onCancel = (e?) => {
    e && e.preventDefault()

    const { model } = this.props

    // if model's name has been set, just cancel edit
    if (model.name) {
      this.updateName(model.name)
      model.updateEditing(false)
    } else {
      // remove virtual section
      this.onDelete()
    }
  }

  private onDelete = () => {
    const { onDelete, index } = this.props

    Modal.showConfirm({
      content: '确认要删除此 section 吗？'
    }).then(() => onDelete(index))
  }

  render() {
    const {
      canDrop,
      isDragging,
      isOver,
      connectDragSource,
      connectDragPreview,
      connectDropTarget,
      index,
      model,
      formModel
    } = this.props
    const isActive = canDrop && isOver

    return (
      <Wrapper active={isActive} isDragging={isDragging}>
        <SecionTarget index={index} />
        {connectDropTarget(
          <div className='section-item'>
            {connectDragPreview(
              <h2 className={`section-header`}>
                {connectDragSource(
                  <span>
                    <Icon className='drag-icon' type='drag' />
                  </span>
                )}

                {model.editing ? (
                  <div className='editor'>
                    <Input
                      ref={ref => (this.inputRef = ref)}
                      type='text'
                      maxLength={64}
                      placeholder='请输入 section 名称'
                      autoFocus
                      onFocus={e => e.target.select()}
                      onKeyDown={this.onKeyDown}
                      value={this.name}
                      onChange={e => this.updateName(e.target.value)}
                    />
                    <div className='operators'>
                      <a href='#' onClick={this.onConfirm}>
                        保存
                      </a>
                      <a href='#' onClick={this.onCancel}>
                        取消
                      </a>
                    </div>
                  </div>
                ) : (
                  <span className='name'>{model.name}</span>
                )}

                <div className='operators'>
                  {!model.editing ? (
                    <EditOutlined onClick={() => model.updateEditing(true)} />
                  ) : null}
                  <DeleteOutlined onClick={this.onDelete} />
                </div>
              </h2>
            )}

            <div className='body'>
              {model.fields.map((item, i) => {
                return (
                  <FormItem
                    key={item._key}
                    formModel={formModel}
                    model={item}
                    index={i}
                    parentIndex={index}
                    onDelete={this.deleteField}
                  />
                )
              })}

              <ItemTarget index={model.fields.length} parentIndex={index} />
            </div>
          </div>
        )}
      </Wrapper>
    )
  }
}
