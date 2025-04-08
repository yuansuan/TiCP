## source

从 lightdesk/web-client ssh://vcssh@phabricator.intern.yuansuan.cn/source/lightdesk.git 拷贝至此

## node
using v18

建议使用nvm管理不同版本的node 

https://github.com/coreybutler/nvm-windows

```
nvm list available
nvm install 18.18.0
nvm use 18.18.0
```

## development
```
npm install
npm run dev

拉起浏览器窗口后，目前暂时需要将/index.html?room_id=xxx&signal=xxx手动复制到url后刷新
```