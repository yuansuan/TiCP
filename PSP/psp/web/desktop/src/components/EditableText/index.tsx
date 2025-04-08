/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, {
  useRef,
  useImperativeHandle,
  forwardRef,
  useLayoutEffect,
} from 'react'
import Icon from '../Icon'
import { StyledLayout } from './style'
import { InputProps } from 'antd/lib/input'
import { observer } from 'mobx-react-lite'
import { Input, Tooltip } from 'antd'
import { useModel, Context, useStore } from './store'
import { useLayoutRect } from '@/utils/hooks'

export type EditableTextProps = {
  style?: object
  inputProps?: InputProps
  onClick?: (event: React.MouseEvent) => void
  defaultEditing?: boolean
  editing?: boolean
  defaultShowEdit?: boolean
  showEdit?: boolean
  defaultValue?: string
  beforeConfirm?: (value: string) => boolean | string | Promise<string | void>
  onConfirm?: (value: string) => void
  beforeCancel?: (value: string) => boolean | string | Promise<string | void>
  onCancel?: (value: string) => void
  filter?: RegExp | ((item: string) => string)
  help?: string
  EditIcon?: React.ReactNode
  Text?: (value: string) => React.ReactNode
  model?: ReturnType<typeof useModel>
  unit?: string
}

const BaseEditableText = observer(
  function EditableText(
    {
      defaultValue,
      defaultShowEdit = true,
      filter,
      onConfirm,
      beforeConfirm,
      onCancel,
      beforeCancel,
      onClick,
      help,
      EditIcon,
      Text,
      style,
      inputProps,
      showEdit = true,
      unit,
    }: EditableTextProps,
    ref: any
  ) {
    const state = useStore()
    const { editing, value, error } = state

    const inputRef = useRef(undefined)
    const [operatorRect, operatorRef, operatorResize] = useLayoutRect()
    const [unitRect, unitRef, unitResize] = useLayoutRect()
    useLayoutEffect(() => {
      operatorResize()
      unitResize()

      if (!editing) {
        state.setError(undefined)
      }
    }, [editing])

    // fix chinese input issue
    const _state = useRef({
      isOnComposition: null,
      emittedInput: null,
    })

    useImperativeHandle(
      ref,
      () => ({
        edit: () => state.setEditing(true),
      }),
      []
    )

    function onEdit(e) {
      e.stopPropagation()
      state.setEditing(true)
    }

    function focusInput() {
      const { current } = inputRef
      if (current) {
        current.focus()
        current.select()
      }
    }

    function _filter(value) {
      if (filter) {
        if (filter instanceof RegExp) {
          value = value.replace(filter, '')
        } else {
          value = filter(value)
        }
      }

      return value
    }

    function onChange(e) {
      let { value } = e.target
      const { current } = _state
      if (!current.isOnComposition) {
        value = _filter(value)
        current.emittedInput = true
      } else {
        current.emittedInput = false
      }
      state.setValue(value)
    }

    function onComposition(event) {
      const { current } = _state
      if (event.type === 'compositionstart') {
        current.isOnComposition = true
        current.emittedInput = false
      } else if (event.type === 'compositionend') {
        current.isOnComposition = false
        if (!current.emittedInput && filter) {
          state.setValue(_filter(value))
        }
      }
    }

    function onKeyDown(e) {
      if (e.keyCode === 13) {
        _onConfirm()
      } else if (e.keyCode === 27) {
        _onCancel()
      }
    }

    function _onConfirm() {
      const confirm = value => {
        onConfirm && onConfirm(value)
        state.setEditing(false)
      }

      // before confirm
      if (beforeConfirm) {
        const res = beforeConfirm(value)
        // async
        if (res instanceof Promise) {
          res
            .then(() => {
              confirm(value)
            })
            .catch(err => {
              state.setError(err || '更新失败')
              focusInput()
            })
        } else if (res === true) {
          confirm(value)
        } else {
          state.setError(res || '更新失败')
          focusInput()
        }
      } else {
        confirm(value)
      }
    }

    function _onCancel() {
      const cancel = value => {
        onCancel && onCancel(value)
        if (defaultValue !== undefined) {
          state.setValue(defaultValue)
        }
        state.setEditing(false)
      }

      // before cancel
      if (beforeCancel) {
        const res = beforeCancel(value)
        // async
        if (res instanceof Promise) {
          res
            .then(() => {
              cancel(value)
            })
            .catch(err => {
              state.setError(err || '取消失败')
              focusInput()
            })
        } else if (res === true) {
          cancel(value)
        } else {
          state.setError(res || '取消失败')
          focusInput()
        }
      } else {
        cancel(value)
      }
    }

    return (
      <StyledLayout style={style} ref={ref}>
        <div
          className='main'
          style={{
            width: `calc(100% - ${Math.ceil(
              operatorRect.width + unitRect.width
            )}px)`,
          }}>
          {editing ? (
            <Input
              ref={inputRef}
              autoFocus
              defaultValue={defaultValue}
              value={value}
              onChange={onChange}
              onKeyDown={onKeyDown}
              onFocus={e => e.target.select()}
              onClick={e => e.stopPropagation()}
              onCompositionStart={onComposition}
              onCompositionEnd={onComposition}
              size='small'
              {...inputProps}
              {...(error && {
                className: `${inputProps?.className || ''} error`,
                suffix: (
                  <Tooltip title={error} visible={true}>
                    <Icon style={{ fontSize: 20 }} type='question_circle' />
                  </Tooltip>
                ),
              })}
            />
          ) : (
            <div
              className={`text ${onClick ? 'isLink' : ''}`}
              onClick={onClick}>
              {Text ? Text(value) : <span title={value}>{value}</span>}
            </div>
          )}
        </div>
        <div className='unit' ref={unitRef}>
          {unit}
        </div>
        <div
          className='operator'
          ref={operatorRef}
          onClick={e => e.stopPropagation()}>
          {editing && (
            <>
              <Icon className='confirm' type='define' onClick={_onConfirm} />
              <Icon className='cancel' type='cancel' onClick={_onCancel} />
              {help && (
                <Tooltip placement='top' title={help}>
                  <Icon className='help' type='question_circle' />
                </Tooltip>
              )}
            </>
          )}
          {!editing && (
            <>
              {showEdit && (
                <span
                  onClick={onEdit}
                  className={`edit ${defaultShowEdit ? '' : 'hoverAction'}`}>
                  {EditIcon || <Icon type='rename' />}
                </span>
              )}
            </>
          )}
        </div>
      </StyledLayout>
    )
  },
  {
    forwardRef: true,
  }
)

const EditableText: React.SFC<EditableTextProps> & {
  useModel?: typeof useModel
} = forwardRef(function EditableText(
  { model, ...props }: EditableTextProps,
  ref: any
) {
  const defaultModel = useModel({
    defaultEditing: props.defaultEditing,
    defaultValue: props.defaultValue,
  })
  const finalModel = model || defaultModel

  return (
    <Context.Provider value={finalModel}>
      <BaseEditableText ref={ref} {...props} />
    </Context.Provider>
  )
})
EditableText.useModel = useModel

export default EditableText
