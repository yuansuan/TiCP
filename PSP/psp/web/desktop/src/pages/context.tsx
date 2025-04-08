import hoistNonReactStatic from 'hoist-non-react-statics'
import * as React from 'react'

const GlobalContext = React.createContext({})

export default GlobalContext

export function inject(selector?) {
  return function Wrapper(WrappedComponent): any {
    class ContextComponent extends React.Component<any> {
      public render() {
        const { forwardedRef, ...rest } = this.props

        return (
          <GlobalContext.Consumer>
            {data => {
              return (
                <WrappedComponent
                  ref={forwardedRef}
                  {...rest}
                  {...(selector ? selector(data) : { context: data })}
                />
              )
            }}
          </GlobalContext.Consumer>
        )
      }
    }

    hoistNonReactStatic(ContextComponent, WrappedComponent)

    const forwardRef = (props, ref) => (
      <ContextComponent {...props} forwardedRef={ref} />
    )
    forwardRef.displayName = `inject(${
      WrappedComponent.displayName || WrappedComponent.name
    })`

    return React.forwardRef(forwardRef)
  }
}
