import React, { useState, useEffect } from 'react'
import { StyleLayout } from './style'
import { message, Select, Button } from 'antd'
import { ValidInput } from '@/components'
import { PlusOutlined, DeleteOutlined } from '@ant-design/icons'
function TodoList({ zoneList = [], onChange, defaultValues = [] }) {
  const [todos, setTodos] = useState(defaultValues)
  const [newKey, setNewKey] = useState('')
  const [newValue, setNewValue] = useState('')
  const [errorMsg, setErrorMsg] = useState('')

  useEffect(() => {
    if (onChange) {
      const newTodo = { key: newKey, value: newValue }
      const keyExists = todos.some(todo => todo.key === newKey)
      if (keyExists) {
        message.error(`「${newKey}」已经存在，请重新添加！`)
        return
      }
      const duplicateTodo = todos
        ?.concat(newTodo)
        ?.filter(todo => todo.key && todo.value)
      onChange(duplicateTodo)
    }
  }, [todos, newKey, newValue])
  const handleKeyChange = (newKey, index) => {
    const currentTodo = todos[index]
    const keyExists = todos.some(todo => todo.key === newKey)
    if (keyExists) {
      message.error(`「${newKey}」已经存在，请重新添加！`)
      return
    }
    if (newKey.length + currentTodo.value.length <= 255) {
      updateTodoAtIndex(index, { key: newKey })
      setErrorMsg('')
    } else {
      setErrorMsg('键和值的总长度不应超过 255 个字符！')
    }
  }

  const handleValueChange = (event, index) => {
    const newValue = event.target.value
    const currentTodo = todos[index]
    if (currentTodo.key.length + newValue.length <= 255) {
      updateTodoAtIndex(index, { value: newValue })
      setErrorMsg('')
    } else {
      setErrorMsg('键和值的总长度不应超过 255 个字符！')
    }
  }

  const handleAddTodo = () => {
    if (newKey && newValue) {
      const keyExists = todos.some(todo => todo.key === newKey)
      if (keyExists) {
        message.error(`「${newKey}」已经存在，请重新添加！`)
      } else {
        const newTodo = { key: newKey, value: newValue }
        setTodos([...todos, newTodo])
        setNewKey('')
        setNewValue('')
      }
    }
  }

  const updateTodoAtIndex = (index, updatedFields) => {
    const updatedTodos = todos.map((todo, i) =>
      i === index ? { ...todo, ...updatedFields } : todo
    )
    setTodos(updatedTodos)
  }

  const handleDeleteTodo = index => {
    const updatedTodos = todos.filter((_, i) => i !== index)
    setTodos(updatedTodos)
  }

  return (
    <StyleLayout>
      {todos?.length > 0 &&
        todos?.map((todo, index) => (
          <div className='todo-item' key={index}>
            <Select
              key={index}
              showSearch
              value={todo.key}
              placeholder='全部'
              onChange={value => handleKeyChange(value, index)}>
              {zoneList?.map(item => {
                return (
                  <Select.Option title={item} key={item} value={item}>
                    {item}
                  </Select.Option>
                )
              })}
            </Select>

            <ValidInput
              placeholder='请输入应用可执行文件绝对路径'
              value={todo.value}
              style={{ width: '250px' }}
              onChange={e => handleValueChange(e, index)}
            />
            <Button
              className='delete-button'
              type='primary'
              danger
              onClick={() => handleDeleteTodo(index)}>
              <DeleteOutlined />
            </Button>
          </div>
        ))}
      <div className='add-todo'>
        <Select
          showSearch
          placeholder='全部'
          value={newKey}
          onChange={value => setNewKey(value)}>
          {zoneList?.map(item => {
            return (
              <Select.Option title={item} key={item} value={item}>
                {item}
              </Select.Option>
            )
          })}
        </Select>
        <ValidInput
          placeholder={'请输入应用可执行文件绝对路径'}
          value={newValue}
          style={{ width: '250px' }}
          onChange={e => setNewValue(e.target.value)}
        />
      </div>
      {newKey && newValue && (
        <div className='add-item-action' onClick={handleAddTodo}>
          <PlusOutlined /> 添加执行路径
        </div>
      )}
    </StyleLayout>
  )
}

export default TodoList
