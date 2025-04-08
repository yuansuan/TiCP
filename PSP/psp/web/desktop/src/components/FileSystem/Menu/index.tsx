import * as React from 'react'
import { observer } from 'mobx-react'
import { computed } from 'mobx'
import { PathHistory, FavoritePoint, RootPoint } from '@/domain/FileSystem'
import Favorite from './FavoritePoint'
import RootPoints from './RootPoints'
import { StyledTreeMenu } from './style'

interface IProps {
  history: PathHistory
  points: Array<RootPoint>
  favoritePoint?: FavoritePoint
  updateSelectedKeys?: (keys: string[], pointId: string) => void
  newDirectories$: any
  selectedKeys: any
  hasPerm?: boolean
}

@observer
export default class Menu extends React.Component<IProps> {
  @computed
  get favoritePoint() {
    return this.props.favoritePoint
  }

  render() {
    const {
      newDirectories$,
      history,
      selectedKeys,
      points,
      hasPerm,
    } = this.props

    return (
      <StyledTreeMenu>
        {this.favoritePoint ? (
          <Favorite
            history={history}
            favoritePoint={this.favoritePoint as any}
            selectedKeys={selectedKeys[(this.favoritePoint as any).pointId]}
          />
        ) : null}
        {points.map(point => (
          <RootPoints
            hasPerm={hasPerm}
            key={point.pointId}
            history={history}
            points={[point]}
            newDirectory$={
              newDirectories$.filter(n => n.pointId === point.pointId)[0][
                'newDirectory$'
              ]
            }
            updateSelectedKeys={keys => {
              this.props.updateSelectedKeys(keys, point.pointId)
            }}
            selectedKeys={selectedKeys[point.pointId] || []}
          />
        ))}
      </StyledTreeMenu>
    )
  }
}
