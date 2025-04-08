/*!
 * Copyright (C) 2016-present, Yuansuan.cn
 */

import React, { useEffect } from 'react'

export const StatisticsComponent = function StatisticsComponent() {
  useEffect(() => {
    const s = document.createElement('script')
    s.text = `
var _hmt = _hmt || [];
(function() {
  var hm = document.createElement("script");
  hm.src = "https://hm.baidu.com/hm.js?773f4a984214f1949fa3a7dbc3fd0217";
  var s = document.getElementsByTagName("script")[0];
  s.parentNode.insertBefore(hm, s);
})();
`
    document.getElementsByTagName('head')[0].appendChild(s)
  }, [])

  return <></>
}
