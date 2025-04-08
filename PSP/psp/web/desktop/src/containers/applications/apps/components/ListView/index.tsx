import React, { useEffect, useMemo, useRef } from 'react'
import { FixedSizeList } from 'react-window'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
import { rootFontSize } from '@/utils/calculateCellDimensions'
import Loading from '../Loading'
import ContentItem from '../ContentItem'

const ListView = ({
  items = [],
  loading,
  rect,
  selected,
  selectedName,
  selectedObj,
  handleClick,
  handleBlur,
  handleDouble,
  changeListView
}) => {
  const listViewRef = useRef(null)
  const itemCount = items.length + 1
  const itemHeight = rootFontSize() * 2.18



  const ListItemRow = ({ index, style, data: { items } }) => {
    const item = items[index]

    return item ? (
      <ContentItem
        id={`${item?.id}`}
        item={item}
        key={item?.id}
        selected={selected}
        selectedName={selectedName}
        selectedObj={selectedObj}
        handleClick={handleClick}
        handleDouble={handleDouble}
        handleBlur={handleBlur}
        changeListView={changeListView}
        style={style}
      />
    ) : null
  }
  function showMenu(selected) {
    if (selected?.length === 0) {
      return 'noFile'
    } else if (selected?.length === 1) {
      return 'singleFile'
    } else {
      return 'multipleFiles'
    }
  }

  const MemoizedListItemCell = useMemo(() => React.memo(ListItemRow), [changeListView])

  useEffect(() => {
    EE.on(EE_CUSTOM_EVENT.SETGRIDLISTSCROLLTOEND, scrollend => {
      if (scrollend) {
        listViewRef.current?.scrollToItem(items.length, 'end')
      }
    })
  }, [listViewRef, items])

  return (
    <div
      className='list-view'
      data-parentcontent='parentcontent'
      data-selected={selected.join(',')}
      data-selectedname={selectedName.join(',')}
      data-selectedobj={JSON.stringify(selectedObj)}
      data-menu={showMenu(selected)}>
      <div className='list-header' style={{width: rect.width}}>
        <div className='header-name'>名称</div>
        <div className='header-type'>大小</div>
        <div className='header-size'>文件类型</div>
        <div className='header-date'>修改日期</div>
      </div>
      {loading ? (
        <Loading />
      ) : (
        <FixedSizeList
          ref={listViewRef}
          itemCount={itemCount}
          itemSize={itemHeight}
          className='list-view-body ys-filemanger-custom-list'
          itemData={{
            items,
            itemCount
          }}
          width={rect.width}
          height={rect.height}
          itemKey={(index, data) =>
            index === itemCount - 1 ? 'loader' : `${data?.items[index]?.id}`
          }>
          {MemoizedListItemCell}
        </FixedSizeList>
      )}
    </div>
  )
}

export default ListView
