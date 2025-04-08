import React from 'react'
import ReactDOM from 'react-dom'
import { BrowserRouter as Router } from 'react-router-dom'
import { createBrowserHistory } from 'history'
import { applyMiddleware, createStore, compose } from 'redux'
import { Provider } from 'react-redux'
import createSagaMiddleware from 'redux-saga'
import rootSaga from '@/saga'
import AppReducer from '@/reducers/AppReducer'
import ScrollToTop from '@/views/ScrollToTop'
import Layout from '@/views/Layout'
import '@/views/app.scss'
import 'bulma/css/bulma.css'

const history = createBrowserHistory()
const sagaMiddleware = createSagaMiddleware()
const composeEnhancers = window.__REDUX_DEVTOOLS_EXTENSION_COMPOSE__ || compose
const store = createStore(
  AppReducer,
  composeEnhancers(applyMiddleware(sagaMiddleware))
)

sagaMiddleware.run(rootSaga)

// work from Chrome 68
if (navigator.keyboard && navigator.keyboard.lock) {
  navigator.keyboard.lock()
  document.addEventListener('keydown', e => {
    if (e.keyCode == 27) {
      // escape
      e.preventDefault()
    }
  })
  document.addEventListener('keyup', e => {
    if (e.keyCode == 27) {
      e.preventDefault()
    }
  })
}

window.$getStore = function() {
  return store.getState()
}
const search = new URLSearchParams(window.location.search)
const access_token = search.get('access_token')
if (access_token) {
  sessionStorage.setItem('access_token', access_token)
}

class App extends React.Component {
  render() {
    return (
      <Provider store={store}>
        <Router basename="/rdp" history={history}>
          <ScrollToTop>
            <Layout />
          </ScrollToTop>
        </Router>
      </Provider>
    )
  }
}

ReactDOM.render(<App />, document.getElementById('react-app'))
