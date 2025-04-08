import * as React from 'react'
import { observer } from 'mobx-react'
import { Subject } from 'rxjs'

interface IProps {
  rename$: Subject<any>
  onClick?: any
  target: string
}

@observer
export default class Rename extends React.Component<IProps> {
  render() {
    const { onClick, rename$, target } = this.props

    const children = React.Children.map(this.props.children, (child) =>
      React.cloneElement(child as React.ReactElement, {
        onClick: () => {
          onClick && onClick()
          rename$.next(target)
        }
      })
    )

    return <>{children}</>
  }
}
