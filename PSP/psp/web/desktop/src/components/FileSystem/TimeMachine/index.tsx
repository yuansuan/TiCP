import * as React from 'react'
import * as ReactDOM from 'react-dom'
import { Breadcrumb, Popover, Input } from 'antd'
import { observer } from 'mobx-react'
import { observable, computed, action } from 'mobx'
import { debounceTime } from 'rxjs/operators'
import { merge, Observable, Subject } from 'rxjs'

import { Http, createMobxStream, formatRegExpStr } from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import { Icon, Button } from '@/components'
import { Point, FavoritePoint, PathHistory } from '@/domain/FileSystem'
import FavoritePopover from './FavoritePopover'
import Link from './Link'
import BackList from './BackList'
import { StyledTimeMachine } from './style'
import './style.css'

interface IState {
  favoriteVisible: boolean
}

interface IProps {
  resize$: Observable<any>
  keyword$: Subject<any>
  favoritePoint?: FavoritePoint
  history: PathHistory
  path: string
  currentPoint: Point | FavoritePoint
  points: Array<Point | FavoritePoint>
  showHistory: boolean
}

@observer
export default class TimeMachine extends React.Component<IProps, IState> {
  @observable favoriteVisible = false
  @observable backListIndex = 0
  @action
  updateFavoriteVisible = visible => (this.favoriteVisible = visible)
  @action
  updateBackListIndex = index => (this.backListIndex = index)

  @observable keyword = ''
  @action
  updateKeyword = keyword => (this.keyword = keyword)

  breadcrumbRef = null
  controllerRef = null
  componentDidMount() {
    const { resize$, keyword$ } = this.props

    // when current path change, clear keyword
    createMobxStream(() => this.props.path)
      .pipe(untilDestroyed(this))
      .subscribe(() => {
        this.updateKeyword('')
      })

    // transmit keyword
    createMobxStream(() => this.keyword)
      .pipe(untilDestroyed(this), debounceTime(100))
      .subscribe(keyword => {
        keyword$.next(keyword)
      })

    // when links change or win.width resize, rerender the breadcrumb
    merge(
      createMobxStream(() => this.links),
      resize$
    )
      .pipe(untilDestroyed(this), debounceTime(300))
      .subscribe(() => {
        // reset backList index
        this.updateBackListIndex(0)

        const linkThreshold = this.controllerRef
          ? this.controllerRef.clientWidth - 220
          : 0
        const breadcrumb: any = ReactDOM.findDOMNode(this.breadcrumbRef)
        if (breadcrumb) {
          // wait for render complete
          setTimeout(() => {
            if (breadcrumb.clientWidth > linkThreshold) {
              let flag = false
              ;[...breadcrumb.children].reduceRight((total, item, index) => {
                if (total > linkThreshold && !flag) {
                  flag = true
                  this.updateBackListIndex(index + 1)
                }
                return total + item.getBoundingClientRect().width
              }, 30)
            }
          }, 0)
        }
      })
  }

  @computed
  get targetNode() {
    const { path, favoritePoint } = this.props
    return favoritePoint.children.filter(item => item.path === path)[0]
  }

  @computed
  get collected() {
    return this.props.path && !!this.targetNode
  }

  @computed
  get links() {
    const { history, currentPoint } = this.props
    const { current } = history

    if (!current) {
      return []
    }

    const { path } = current
    let targetPoint = currentPoint
    let prefix = currentPoint.name

    if (!targetPoint) {
      return []
    }

    let linkPath = ''
    return path
      .replace(new RegExp(`^${formatRegExpStr(targetPoint.rootPath)}`), prefix)
      .split(/[\\/]/)
      .map((item, index) => {
        if (index === 0) {
          linkPath = targetPoint.rootPath
        } else {
          linkPath += `/${item}`
        }

        return {
          point: targetPoint,
          name: item,
          path: linkPath,
          disabled: false,
        }
      })
  }

  private hideFavorite = () => {
    this.updateFavoriteVisible(false)
  }

  private onFavoriteVisibleChange = visible => {
    this.updateFavoriteVisible(visible)
  }

  private cancelCollect = () => {
    const { favoritePoint, path } = this.props

    if (!path) {
      return
    }

    const { targetNode } = this
    targetNode &&
      Http.delete('/file/favorite', {
        params: { id: targetNode.originId },
      }).then(() => favoritePoint.fetch())
  }

  render() {
    const { history: pathHistory, favoritePoint, showHistory } = this.props
    const { current } = pathHistory

    return (
      <StyledTimeMachine>
        {showHistory && (
          <>
            <div className='skip'>
              <Button
                disabled={pathHistory.prevDisabled}
                onClick={pathHistory.prev}>
                <Icon style={{ marginRight: 0 }} type='left' />
              </Button>
              <Button
                disabled={pathHistory.nextDisabled}
                onClick={pathHistory.next}>
                <Icon style={{ marginRight: 0 }} type='right' />
              </Button>
            </div>
            <div className='controller' ref={ref => (this.controllerRef = ref)}>
              <div className='breadcrumb'>
                {this.backListIndex > 0 ? (
                  <div className='backList'>
                    <Popover
                      placement='bottom'
                      trigger='click'
                      overlayClassName='DataManagement_BackList_Popover'
                      content={
                        <BackList
                          onClick={link =>
                            pathHistory.push({
                              source: link.point,
                              path: link.path,
                            })
                          }
                          links={this.links.slice(0, this.backListIndex)}
                        />
                      }>
                      <Icon type='backdate' className='listIcon' />
                    </Popover>
                  </div>
                ) : null}
                <Breadcrumb
                  ref={ref => (this.breadcrumbRef = ref)}
                  separator={
                    <div
                      style={{
                        display: 'inline-block',
                        position: 'relative',
                        top: '-4px',
                      }}>
                      /
                    </div>
                  }>
                  {this.links.slice(this.backListIndex).map((link, index) => (
                    <Breadcrumb.Item key={index}>
                      <Link
                        key={index}
                        onClick={() => {
                          pathHistory.push({
                            source: link.point,
                            path: link.path,
                          })
                        }}
                        title={link.path}
                        name={link.name}
                        disabled={link.disabled}
                        active={index === this.links.length - 1}
                      />
                    </Breadcrumb.Item>
                  ))}
                </Breadcrumb>
              </div>
              <div className='collect'>
                {favoritePoint ? (
                  current ? (
                    this.collected ? (
                      <Icon
                        style={{ color: 'orange' }}
                        type='star-filled'
                        onClick={this.cancelCollect}
                      />
                    ) : (
                      <Popover
                        placement='bottom'
                        content={
                          <FavoritePopover
                            path={current.path}
                            favoritePoint={favoritePoint}
                            hide={this.hideFavorite}
                          />
                        }
                        title='添加至 Favorite'
                        trigger='click'
                        visible={this.favoriteVisible}
                        onVisibleChange={this.onFavoriteVisibleChange}>
                        <Icon type='star' />
                      </Popover>
                    )
                  ) : (
                    <Icon type='star' />
                  )
                ) : null}
              </div>
            </div>
          </>
        )}
        <div className='filter'>
          <Input.Search
            style={{ height: 32, width: 350 }}
            placeholder='按文件名称过滤'
            value={this.keyword}
            onChange={e => this.updateKeyword(e.target.value)}
          />
        </div>
      </StyledTimeMachine>
    )
  }
}
