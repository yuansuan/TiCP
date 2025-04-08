import React, { useRef } from 'react'
import { Image } from '@/utils/general'
import { Input, Tooltip } from 'antd'
import '../assets/fileexpo.css'
import { formatByte } from '@/utils/Validator'
import { useSelector } from 'react-redux'

export default function ContentItem({
  id,
  item,
  selected,
  selectedName,
  selectedObj,
  handleClick,
  handleDouble,
  handleBlur,
  changeListView,
  style,
}) {
  const files = useSelector(state => state.files)
  const fdata = files.data.getId(files.cdir)

  const inputRef = useRef(null)

  const fileType = item?.isFile
    ? `${item?.name?.split('/').pop()?.split('.').slice(1).pop() || ''}文件`
    : '文件夹'

  return (
    <div
      id={id}
      key={item?.id}
      style={style}
      data-menu={'singleFile'}
      className={`conticon hvtheme flex  items-center prtclk ${
        selected.includes(item?.id) ? 'selected' : ''
      } ${changeListView ? 'flex-col' : 'list-item'}`}
      data-id={item?.id}
      data-originid={item?.originId}
      data-type={item?.type}
      data-path={item?.path}
      data-size={item?.size}
      data-name={item?.name}
      data-isfile={item?.isFile}
      data-onlyread={item?.readOnly}
      data-focus={selected?.includes(item?.id)}
      onClick={(e) =>handleClick(e)}
      onDoubleClick={(e) => handleDouble(e)}>
      <div className={changeListView ? 'card-name' : 'list-item-name'}>
        {changeListView ? (
          <Image src={`icon/win/${item?.info?.icon}`} w={80} h={80} />
        ) : null}
        {item?.editFlag ? (
          <Input
            ref={inputRef}
            size='small'
            autoFocus
            maxLength={64}
            key={item?.id}
            onFocus={event => event.target?.focus()}
            onPressEnter={e =>
              handleBlur(e, item?.path, item?.id, item?.name, item?.isMkdir,files?.cdir,fdata?.path)
            }
            defaultValue={item?.name}
            onBlur={e =>
              handleBlur(e, item?.path, item?.id, item?.name, item?.isMkdir,files?.cdir,fdata?.path)
            }
          />
        ) : (
          <div className='item-name-inner'>
            {!changeListView && (
              <Image src={`icon/win/${item?.info?.icon}`} w={25} h={25} />
            )}
            {
              <Tooltip title={item?.name}>
                <span>{item?.name}</span>
              </Tooltip>
            }
          </div>
        )}
      </div>

      {!changeListView ? (
        <>
          <div className='item-size'>{formatByte(item?.size)}</div>
          <div className='item-type'>{fileType}</div>
          <div className='item-date'>{item?.mtime}</div>
        </>
      ) : null}
    </div>
  )
}
