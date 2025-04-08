import { createStore } from '@/utils/reducer'
import { useLocalStore } from 'mobx-react-lite'
import { vis } from '@/domain'

export function useModel() {
  return useLocalStore(() => ({
    vis
  }))
}

const store = createStore(useModel)

export const Provider = store.Provider
export const Context = store.Context
export const useStore = store.useStore
