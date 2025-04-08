/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import * as React from 'react'
import { Subject, fromEvent, from, of } from 'rxjs'
import {
  filter,
  debounceTime,
  tap,
  switchMap,
  map,
  mergeMap,
  throttleTime,
  takeUntil,
  startWith,
  catchError
} from 'rxjs/operators'
import { observer } from 'mobx-react'
import { observable, action, computed } from 'mobx'
import qs from 'querystring'
import { message, Spin } from 'antd'

import { createMobxStream, Http, formatRegExpStr } from '@/utils'
import { untilDestroyed } from '@/utils/operators'
import { currentUser } from '@/domain'
import {
  PathHistory,
  Point,
  RootPoint,
  FavoritePoint,
  List as DomainFileList
} from '@/domain/FileSystem'
import {
  List as FileList,
  TimeMachine,
  // Toolbar,
  Menu
} from '@/components/FileSystem'
import { IToolbarConfig } from '@/components/FileSystem/Toolbar'
import { StyledSuite } from './style'
import { Icon } from '@/components'

interface IProps {
  defaultPath?: string
  showMenu?: boolean
  showFavorite?: boolean
  showHistory?: boolean
  selectedKeys$?: Subject<string[]>
  toolbar?: Partial<IToolbarConfig>
  points?: RootPoint[]
}

interface ISuite {
  points: Array<Point>
  history: PathHistory
  fileList: DomainFileList
  favoritePoint?: FavoritePoint
}

@observer
export default class Suite extends React.Component<IProps> implements ISuite {
  @observable history = new PathHistory()
  @observable fileList = new DomainFileList()
  @observable favoritePoint = null
  @observable points = this.props.points
  @observable fetching = false
  @observable selectedMenuKeys = {}
  @action
  updateSelectedMenuKeys = (keys, pointId) => {
    this.selectedMenuKeys = {}
    this.selectedMenuKeys[pointId] = keys
  }

  @action
  updateFetching = fetching => (this.fetching = fetching)

  @observable loading = false
  @observable loadingText = null
  @action
  showLoading = (text = null) => {
    this.loading = true
    this.loadingText = text
  }
  @action
  hideLoading = () => (this.loading = false)

  selectedKeys$ = new Subject<string[]>()
  keyword$ = new Subject<string>()
  newDirectories$ = this.props.points.map(point => ({
    pointId: point.pointId,
    newDirectory$: new Subject<any>()
  }))
  resize$ = new Subject<any>()
  resizeBar$ = new Subject<any>()

  menuWrapperRef: React.RefObject<HTMLDivElement> = React.createRef()

  constructor(props: IProps) {
    super(props)

    // control selectedKeys
    if (props.selectedKeys$) {
      this.selectedKeys$ = props.selectedKeys$
    }

    // mountList
    const data = currentUser.mountList
    const favoriteNode = data.find(item => item.id === 'favorites')

    // favorite node
    if (favoriteNode && props.showFavorite !== false) {
      this.favoritePoint = new FavoritePoint({ ...favoriteNode })
    }

    // config default path
    const { defaultPath } = props
    if (defaultPath !== undefined) {
      const source = this.points.find(
        point =>
          defaultPath === point.rootPath ||
          // windows/linux compatible
          new RegExp(
            `^${formatRegExpStr(point.rootPath).replace(/[\\/]$/, '')}[\\/]`
          ).test(defaultPath)
      )

      if (source) {
        this.history.push({
          source,
          path: defaultPath
        })
      } else {
        // invalid defaultPath
        message.error('无法定位当前路径文件')
      }
    } else {
      // default select homePoint
      const homePoint = this.points.find(item => item.pointId === 'home')
      homePoint &&
        this.history.push({
          source: homePoint,
          path: homePoint.rootPath
        })
    }
  }

  @computed
  get currentPath() {
    const { history } = this
    return history.currentPath
  }

  @computed
  get currentPoint() {
    const { current } = this.history
    return current.source
  }

  getPointByPath(path) {
    if (!path) {
      return null
    }

    return this.points.find(
      point =>
        point.rootPath &&
        (point.rootPath === path ||
          // windows/linux compatible
          new RegExp(
            `^${formatRegExpStr(point.rootPath).replace(/[\\/]$/, '')}[\\/]`
          ).test(path))
    )
  }

  async componentDidMount() {
    const { fileList } = this

    // path jump
    const query = qs.parse(location.hash.split('?')[1])
    if (query && query.path) {
      let { path: queryPath } = query
      const { files } = await Http.get('/file/detail', {
        params: {
          paths: queryPath
        }
      }).then(res => res.data)
      const file = files[0]
      if (!file.is_dir) {
        queryPath = file.path.replace(/[\\/][^\\/]*$/, '')
      }
      const point = this.getPointByPath(queryPath)
      if (point) {
        this.history.push({
          path: queryPath as string,
          source: point
        })
      } else {
        message.error(`无法定位 ${queryPath} 到相应目录`)
      }
    }

    // monitor current history
    const currentPath$ = createMobxStream(() => this.currentPath).pipe(
      untilDestroyed(this),
      filter(path => !!path),
      debounceTime(300)
    )
    // when current path changed, update fileList
    currentPath$.subscribe(path => fileList.update(path))
    // when history changed, update menu key
    createMobxStream(() => {
      const { current } = this.history

      return current
        ? {
            path: current.path,
            source: current.source
          }
        : undefined
    })
      .pipe(
        untilDestroyed(this),
        filter(history => !!history)
      )
      .subscribe(history => {
        const { source, path } = history

        if (source === this.favoritePoint) {
          const node = this.favoritePoint.filterFirstNode(
            item => item.path === path
          )
          node &&
            this.updateSelectedMenuKeys(
              [node.favoriteId],
              this.favoritePoint.pointId
            )
        } else {
          this.updateSelectedMenuKeys([path], source.pointId)
        }
      })

    // when current path changed, fetch new files
    currentPath$
      .pipe(
        tap(() => this.updateFetching(true)),
        switchMap((path: string) => {
          if (!this.currentPoint) {
            // trigger subscriber with undefined
            return of(undefined)
          }

          const { service } = this.currentPoint
          let fetch$ = null
          if (this.currentPoint.filterFirstNode(node => node.path === path)) {
            fetch$ = from(service.fetch(path))
          } else {
            fetch$ = from(
              service
                .fetchTree({
                  path,
                  rootPath: this.currentPoint.rootPath
                })
                .then(() => service.fetch(path))
            )
          }

          return fetch$.pipe(
            // catch error to avoid abort
            catchError(() => {
              return of(undefined)
            })
          )
        })
      )
      .subscribe(() => this.updateFetching(false))

    // body onresize
    fromEvent(window, 'resize')
      .pipe(untilDestroyed(this), startWith(''), debounceTime(300))
      .subscribe(this.resize$)

    this.resizeBar$
      .pipe(
        map((event: React.MouseEvent) => ({
          event: event.nativeEvent
        })),
        untilDestroyed(this),
        map(({ event }) => ({
          event,
          width: this.menuWrapperRef.current.getBoundingClientRect().width
        })),
        mergeMap(({ event, width }) =>
          fromEvent(window, 'mousemove').pipe(
            tap((moveEvent: any) => moveEvent.preventDefault()),
            throttleTime(30),
            map((moveEvent: MouseEvent) => ({
              diffX: moveEvent.pageX - event.pageX
            })),
            map(({ diffX }) => {
              const w = Math.min(width + diffX, 600)
              return {
                width: w
              }
            }),
            takeUntil(fromEvent(window, 'mouseup'))
          )
        )
      )
      .subscribe(({ width }) => {
        this.menuWrapperRef.current.style.width = width + 'px'
      })
  }

  render() {
    const { showMenu = true, toolbar, showHistory = true } = this.props
    const {
      resize$,
      favoritePoint,
      points,
      history,
      fileList,
      currentPoint,
      currentPath,
      loading,
      showLoading,
      hideLoading
    } = this

    return (
      <StyledSuite>
        <div className='header'>
          <div className='timeMachine'>
            <TimeMachine
              showHistory={showHistory}
              favoritePoint={favoritePoint}
              keyword$={this.keyword$}
              currentPoint={currentPoint}
              points={points}
              history={history}
              path={currentPath}
              resize$={resize$}
            />
          </div>
          <div className='toolbar'>
            {/* <Toolbar
              hasPerm={true}
              history={history}
              config={toolbar}
              sourcePoint={history.current && history.current.source}
              currentPoint={currentPoint}
              points={points}
              parentPath={currentPath}
              newDirectories$={this.newDirectories$}
              selectedKeys$={this.selectedKeys$}
              showLoading={showLoading}
              hideLoading={hideLoading}
              selectedMenuKeys={this.selectedMenuKeys}
            /> */}
          </div>
        </div>

        <div className='main'>
          {showMenu && (
            <>
              <div className='menu'>
                <div className='menuWrapper' ref={this.menuWrapperRef}>
                  <Menu
                    hasPerm={true}
                    favoritePoint={favoritePoint}
                    points={points}
                    history={history}
                    updateSelectedKeys={this.updateSelectedMenuKeys}
                    selectedKeys={this.selectedMenuKeys}
                    newDirectories$={this.newDirectories$}
                  />
                </div>
                <div
                  title='resize'
                  className='resizeBar'
                  onMouseDown={e => this.resizeBar$.next(e)}>
                  <div className='dragIcon'>
                    <Icon type='drag' />
                  </div>
                </div>
              </div>
            </>
          )}
          <div className='panel'>
            <FileList
              hasPerm={true}
              resize$={resize$}
              selectedKeys$={this.selectedKeys$}
              keyword$={this.keyword$}
              fetching={this.fetching}
              currentPoint={currentPoint}
              points={points}
              history={history}
              fileList={fileList}
            />
          </div>
        </div>
        {loading && (
          <div className='loading'>
            <Spin tip={this.loadingText} />
          </div>
        )}
      </StyledSuite>
    )
  }
}
