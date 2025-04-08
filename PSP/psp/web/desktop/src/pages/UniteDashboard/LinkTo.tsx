import * as React from 'react'
import { useDispatch } from 'react-redux'
import { history } from '@/utils'

interface IProps {
  render: (props) => React.ReactNode,
  routerPath: string
}

export function LinkTo(props: IProps) {
    const { render, routerPath } = props
    const dispatch = useDispatch()
    
    const goTo = () => {
      history.push(routerPath)
      window.localStorage.setItem(
        'CURRENTROUTERPATH',
        routerPath
      )
      dispatch({
        type: 'ENTERPRISEMANAGE',
        payload: 'togg'
      })
    }

    return (
      <>
        {
          render({goTo})
        }
      </>
    )
}