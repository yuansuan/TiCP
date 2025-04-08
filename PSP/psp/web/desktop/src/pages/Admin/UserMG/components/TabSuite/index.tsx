import { action, observable } from 'mobx'
import { observer } from 'mobx-react'
import * as React from 'react'
import { SuiteWrapper } from './style'
import Toolbar from './Toolbar'
import TabPagination from './TabPagination'
import { ListQuery } from '../../utils'

interface IProps {
  operators?: React.ReactNode
  freshen: () => Promise<any>
  defaultListQuery: ListQuery
  Panel: React.FunctionComponent<any>
  store: any
  searchPlaceholder?: string
}

@observer
export default class TabSuite extends React.Component<IProps> {
  wrapperRef = null
  resizeObserver = null

  constructor(props) {
    super(props)
    this.wrapperRef = React.createRef()
  }

  @observable height = 400
  @observable public width
  @observable public loading = false
  @observable public listQuery: ListQuery = this.props.defaultListQuery || {
    query: '',
    page: 1,
    pageSize: 10
  }
  @observable public total = 0

  public panelRef = null
  @action
  public updateWidth = width => (this.width = width)
  @action
  public updateLoading = loading => (this.loading = loading)
  @action
  public updateListQuery = newQuery => (this.listQuery = newQuery)
  @action
  public updateTotal = total => (this.total = total)

  public componentDidMount() {
    this.freshen()
    this.getPanelSize()

    this.resizeObserver = new ResizeObserver((entries) => {
      for (let entry of entries) {
          this.height = entry.contentRect.height
      }
    })
    
    this.resizeObserver.observe(this.wrapperRef.current)

    // hack: 处理Table首次加载 bug
    setTimeout(() => {
      this.wrapperRef.current.style.paddingRight = '1px'
    }, 3000)
  }

  componentWillUnmount() {
    this.resizeObserver && this.resizeObserver.disconnect()
  }

  public freshen = () => {
    const { freshen } = this.props
    this.updateLoading(true)
    freshen().finally(() => this.updateLoading(false))
  }

  public render() {
    const { width, loading, listQuery, updateListQuery, total, updateTotal } =
      this
    const { operators, store, Panel, searchPlaceholder } = this.props

    return (
      <SuiteWrapper ref={this.wrapperRef}>
        <Toolbar
          operators={operators}
          listQuery={listQuery}
          updateListQuery={updateListQuery}
          placeholder={searchPlaceholder || ''}
        />
        <div ref={ref => (this.panelRef = ref)}>
          <Panel
            {...{
              width,
              loading,
              store,
              listQuery,
              updateTotal,
              height: this.height + 20
            }}
          />
        </div>
        {total > 0 && (
          <div className='pagination'>
            <TabPagination
              total={total}
              listQuery={listQuery}
              updateListQuery={updateListQuery}
            />
          </div>
        )}
      </SuiteWrapper>
    )
  }

  private getPanelSize = () => {
    if (this.panelRef) {
      this.updateWidth(this.panelRef.clientWidth)
    }
  }
}
