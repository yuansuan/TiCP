
let rootFontSizeCache: number | undefined = undefined;

const getRootFontSize = () => {
    return parseFloat(window.getComputedStyle(document.querySelector('html') as Element).getPropertyValue('font-size'));
};

export const rootFontSize = (reset?: boolean) => {
    if (rootFontSizeCache === undefined || reset === true) {
        rootFontSizeCache = getRootFontSize();
    }
    return rootFontSizeCache;
};

export const calculateCellDimensions = (areaWidth: number) => {
  const itemWidth = 6.25 * rootFontSize(); 
  const itemHeight = 8.25 * rootFontSize(); 

  const rowItemCount = Math.floor(areaWidth / itemWidth);
  const expandedItemWidth = areaWidth / rowItemCount;
  const squishedItemWidth = areaWidth / (rowItemCount + 1);
  const oversizing = expandedItemWidth - itemWidth;
  const oversquishing = itemWidth - squishedItemWidth;
  const ratio = itemHeight / itemWidth;

  // If expanded width is less imperfect than squished width
  if (oversizing <= oversquishing) {
      return {
          cellWidth: expandedItemWidth,
          cellHeight: expandedItemWidth * ratio,
      };
  }

  return {
      cellWidth: squishedItemWidth,
      cellHeight: squishedItemWidth * ratio,
  };
};
