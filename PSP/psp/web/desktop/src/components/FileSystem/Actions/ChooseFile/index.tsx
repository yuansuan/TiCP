import * as React from 'react'
import { observer } from 'mobx-react'
import { Subject, empty } from 'rxjs'
import { switchMap } from 'rxjs/operators'

import { Button, Modal } from '@/components'
import { RootPoint } from '@/domain/FileSystem'
import { formatRegExpStr } from '@/utils'
import TreeMenu from '../../TreeMenu'
import { FileChooserFooter, StyledTreeMenu } from './style'

export const showFileChooser = ({
  path = '',
  disabledKeys = [],
  points = undefined,
  onOk = undefined,
  onCancel = undefined,
  hasPerm = true
} = {}) => {
  const selectedKeys$: Subject<string[]> = new Subject()
  const newDirectory$ = new Subject()
  let selectedKeys = [],
    currPoint = null
  const subscription = selectedKeys$.subscribe(keys => {
    selectedKeys = keys
  })

  // when selectedKeys changed, fetch new files
  const childSubscription = selectedKeys$
    .pipe(
      switchMap((keys: string[]) => {
        const path = keys[0]
        const point = points.find(
          item =>
            path === item.path ||
            // windows/linux compatible
            new RegExp(
              `^${formatRegExpStr(item.path).replace(/[\\/]$/, '')}[\\/]`
            ).test(path)
        )

        if (!point) {
          return empty()
        }

        const { service } = point

        if (point.filterFirstNode(node => node.path === path)) {
          return service.fetch(path)
        } else {
          return service
            .fetchTree({ path, rootPath: point.rootPath })
            .then(() => service.fetch(path))
        }
      })
    )
    .subscribe(() => {})

  subscription.add(childSubscription)

  return Modal.show({
    title: '选择文件',
    bodyStyle: { height: 300, padding: 0 },
    onOk: next => {
      if (onOk) {
        onOk(selectedKeys, next, currPoint)
      } else {
        next(selectedKeys)
      }
    },
    onCancel: next => {
      if (onCancel) {
        onCancel(selectedKeys, next)
      } else {
        next(selectedKeys)
      }
    },
    content: (
      <StyledTreeMenu>
        <TreeMenu
          points={points}
          path={path}
          editable={false}
          disbaledKeys={disabledKeys.length !== 0 ? disabledKeys : [path]}
          onExpand={({ keys, context, point }) => {
            const { node } = context
            if (context.expanded) {
              const path = node.props.eventKey
              point.service.fetch(path)
            }
          }}
          onSelect={({ keys, context, point }) => {
            currPoint = point
          }}
          newDirectory$={newDirectory$}
          selectedKeys$={selectedKeys$}
          hasPerm={hasPerm}
        />
      </StyledTreeMenu>
    ),
    footer: ({ onCancel, onOk }) => (
      <FileChooserFooter>
        <div className='left'>
          <Button
            type='primary'
            icon='folder-add'
            ghost
            onClick={() => newDirectory$.next()}>
            新建文件夹
          </Button>
        </div>
        <div className='right'>
          <Button onClick={onCancel}>取消</Button>
          <Button type='primary' onClick={onOk}>
            确定
          </Button>
        </div>
      </FileChooserFooter>
    )
  }).finally(() => {
    subscription.unsubscribe()
  })
}

interface IProps {
  points: RootPoint[]
  path?: string
  disabledKeys?: string[]
  onClick?: () => void
  onOk?: (keys: [], next, currPoint) => void
  onCancel?: (keys: [], next) => void
  hasPerm?: boolean
}

@observer
export default class ChooseFile extends React.Component<IProps> {
  render() {
    const {
      onClick,
      onOk,
      onCancel,
      points,
      path,
      disabledKeys = [],
      hasPerm
    } = this.props

    const children = React.Children.map(this.props.children, child =>
      React.cloneElement(child as React.ReactElement, {
        onClick: () => {
          onClick && onClick()
          showFileChooser({
            points,
            path,
            onOk,
            onCancel,
            disabledKeys,
            hasPerm
          })
        }
      })
    )

    return <>{children}</>
  }
}
