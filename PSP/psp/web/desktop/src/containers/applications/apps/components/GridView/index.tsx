import React, { useEffect, useMemo, useRef } from 'react'
import { FixedSizeGrid } from 'react-window'
import { EE, EE_CUSTOM_EVENT } from '@/utils/Event'
import { calculateCellDimensions } from '@/utils'
import Loading from '../Loading'
import ContentItem from '../ContentItem'

const GridView = ({
  items = [],
  rect,
  loading,
  selected,
  selectedName,
  selectedObj,
  handleClick,
  handleBlur,
  handleDouble,
  changeListView
}) => {
  const gridViewRef = useRef(null)
  const totalItems = items.length + (false ? 1 : 0)

  const width = rect?.width ?? 0
  const { cellHeight, cellWidth } = calculateCellDimensions(width)

  const itemsPerRow = Math.floor(width / cellWidth)
  const rowCount = Math.ceil(totalItems / itemsPerRow)

  function scrollToBottom(ref) {
    const viewInstance = ref.current

    if (viewInstance) {
      if (changeListView) {
        const rowCount = viewInstance.props.rowCount
        viewInstance.scrollToItem({
          columnIndex: 0,
          rowIndex: rowCount - 1
        })
      } else {
        viewInstance.scrollToItem(viewInstance.props.itemCount - 1)
      }
    }
  }

  useEffect(() => {
    EE.on(EE_CUSTOM_EVENT.SETGRIDLISTSCROLLTOEND, scrollend => {
      if (scrollend) {
        if (changeListView) {
          scrollToBottom(gridViewRef)
        }
      }
    })
  }, [items])

  const GridItemCell = ({
    columnIndex,
    rowIndex,
    style,
    data: { items, rowCount, itemsPerRow }
  }) => {
    const currentIndex = columnIndex + rowIndex * itemsPerRow
    const item = items[currentIndex]

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

  const generateItemKey = ({
    columnIndex,
    rowIndex,
    data: { items, itemsPerRow }
  }: {
    columnIndex: number
    rowIndex: number
    data
  }) => {
    const item = items[columnIndex + rowIndex * itemsPerRow]
    return item?.id ?? `${columnIndex}-${rowIndex}`
  }

  const MemoizedGridItemCell = useMemo(() => React.memo(GridItemCell), [changeListView])

  return loading ? (
    <Loading />
  ) : (
    <FixedSizeGrid
      ref={gridViewRef}
      itemData={{
        items: items,
        itemsPerRow,
        rowCount
      }}
      columnWidth={cellWidth}
      rowHeight={cellHeight}
      className='grid-view-body ys-filemanger-custom-list'
      height={rect.height}
      width={rect.width}
      columnCount={itemsPerRow}
      rowCount={rowCount}
      itemKey={generateItemKey}>
      {MemoizedGridItemCell}
    </FixedSizeGrid>
  )
}

export default GridView
