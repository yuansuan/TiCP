import * as React from 'react'
import debounce from 'lodash.debounce'
import { Input } from 'antd'
import styled from 'styled-components'

interface InputWrapperProps {
  mode?: 'dark' | 'normal'
}

const InputWrapper = styled.div<InputWrapperProps>`
  .ant-input-affix-wrapper {
    padding: 0 8px;
    .ant-input:not(:last-child) {
      padding-right: 48px;
      height: 30px;
      outline: none;
    }
  }

  ${props =>
    (props as any).mode === 'dark'
      ? `.ant-input {
    background: rgb(94 46 126/ 80%);
    color: #fff;
    border: 1px solid rgb(94 46 126);
  }
  .ant-input-search-icon {
    color: #fff;
  }
  .ant-input-clear-icon {
    color: #fff;
  }
  `
      : ``}
`

interface SearchProps {
  onSearch: (searchValue: string) => void
  debounceWait?: number
  placeholder?: string
  className?: string
  style?: Object
  mode?: 'dark' | 'normal'
  defaultValue?: string
}

interface SearchState {
  value: string
}

export default class Search extends React.Component<SearchProps, SearchState> {
  static defaultProps = {
    debounceWait: 0
  }

  state = {
    value: ''
  }

  debouncedSearch = debounce(() => {
    this.props.onSearch(this.state.value)
  }, this.props.debounceWait)

  handleChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    this.setState(
      {
        value: e.target.value
      },
      this.debouncedSearch
    )
  }

  componentDidMount() {
    this.setState(
      {
        value: this.props.defaultValue || ''
      },
      () => {
        if (this.props.defaultValue) {
          this.props.onSearch && this.props.onSearch(this.props.defaultValue)
        }
      }
    )
  }

  render() {
    const { onSearch, debounceWait, placeholder, className, ...rest } =
      this.props

    return (
      <InputWrapper mode={this.props.mode || 'normal'}>
        <Input.Search
          autoComplete='off'
          style={{ width: 200, height: 32 }}
          value={this.state.value}
          maxLength={64}
          onChange={this.handleChange}
          placeholder={placeholder}
          className={className}
          allowClear
          {...rest}
        />
      </InputWrapper>
    )
  }
}
