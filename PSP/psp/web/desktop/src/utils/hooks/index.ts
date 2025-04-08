import { useEffect, useRef } from 'react'

export { useAsync } from './useAsync'
export { useInterval } from './useInterval'
export { useElementRect } from './useElementRect'
export { useLayoutRect } from './useLayoutRect'

export function useDidUpdate(fn, inputs?: any[]) {
    const didMountRef = useRef(false)
  
    useEffect(() => {
      if (didMountRef.current) {
        fn()
      } else {
        didMountRef.current = true
      }
    }, inputs)
  }