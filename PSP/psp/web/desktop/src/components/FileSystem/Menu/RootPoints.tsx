import * as React from 'react'
import { observer } from 'mobx-react'
import { Subject } from 'rxjs'

// import { DeleteAction } from "../Actions";
import { PathHistory, RootPoint } from '@/domain/FileSystem'
import TreeMenu from '../TreeMenu'

interface IProps {
  history: PathHistory
  updateSelectedKeys: (keys: string[]) => void
  selectedKeys: string[]
  newDirectory$: Subject<any>
  points: RootPoint[]
  hasPerm?: boolean
}

@observer
export default class Home extends React.Component<IProps> {
  private onSelect = ({ keys, context, point }) => {
    const { history } = this.props

    const path = context.node.props.eventKey
    const node = point.filterFirstNode(item => item.path === path)

    node && history.push({ source: point, path: node.path })
  }

  render() {
    const { onSelect } = this

    const { selectedKeys, newDirectory$, points, hasPerm } = this.props

    return (
      <TreeMenu
        hasPerm={hasPerm}
        points={points}
        selectedKeys={selectedKeys}
        onSelect={onSelect}
        updateSelectedKeys={this.props.updateSelectedKeys}
        newDirectory$={newDirectory$}
      />
    )
  }
}
