import { useState, useCallback } from 'react'

export function useAsync(asyncFun: any, keepRes?: boolean) {
  const [data, setData] = useState(null)
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState(null)

  const execute = useCallback(() => {
    setLoading(true)
    setData(null)
    setError(null)

    return asyncFun()
      .then(res => {
        setLoading(false)
        setData(keepRes ? res : res.data)
        return res
      })
      .catch(err => {
        setLoading(false)
        setError(err)
      })
  }, [asyncFun])

  return { execute, loading, data, error }
}
