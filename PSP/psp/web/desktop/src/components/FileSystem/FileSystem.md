# FileSystem

FileSystem 定义了 PaaS_v3.0 的文件管理系统。相关部件定义如下：

## Actions

文件操作相关命令。Actions 使用 HOC 复用文件操作逻辑。

- **ChooseFile**: 选择文件 action，可以用于文件移动到、复制到等操作；
- **Compress**: 文件压缩；
- **Delete**: 文件删除；
- **Download**: 文件下载；
- **Edit**: 文件编辑；
- ...

## Editor

文件编辑器，接收如下 props：

- **win**: 当前窗口实例；
- **path**: 远程文件 path；
- **readOnly**: 只读模式；
- **viewContent**: 获取文件内容 api；
- **saveContent**: 文件保存 api；

## List

文件列表。接收如下 props：

- **win**: 当前窗口实例；
- **point**: 文件管理节点，提供了文件操作 api；
- **fetching**: 文件获取中标志；
- **fileList**: 文件列表；
- **RxJS Stream**: 通过 Stream 的监听和传值与外界交互
  - **selectedKeys\$**: 当前选中的文件键值；
  - **keyword\$**: 文件搜索关键字；
  - **resize\$**: 窗口大小调整流；
  - **rename\$**: 文件重命名；

## Menu

文件目录树菜单，用于展示文件系统结构和定位文件。接收如下 props：

- **win**: 当前窗口实例；
- **history**: 路径历史（@/domain/FileSystem/PathHistory）；
- **points**: 需要展示的文件节点；
- **favoritePoint**: 是否配置收藏夹节点；
- **newDirectory\$**: 创建文件命令流；
- **selectedKeys**: 当前选中的菜单；

## TimeMachine

文件历史管理器。用于记录文件跳转历史，以及前进和后退的功能。接收如下 props：

- **resize\$**: 窗口 resize 事件流，用于历史路径的自适应显示；
- **favoritePoint**: 是否启用收藏夹功能；
- **history**: 路径历史（@/domain/FileSystem/PathHistory）；
- **path**: 文件路径；
- **points**: 启用的文件节点；

## Toolbar

文件操作工具栏。提供了文件操作、文件搜索和自定义 action 的功能。接收如下 props：

- **win**: 当前窗口实例；
- **RxJS Stream**: 通过 Stream 的监听和传值与外界交互
  - **selectedKeys\$**: 当前选中的文件键值；
  - **keyword\$**: 文件搜索关键字；
  - **rename\$**: 文件重命名；
  - **newDirectory\$**: 新建文件夹；
- **point**: 文件管理节点；
- **parentPath**: 文件父路径；
- **config**: 文件工具栏配置，可以自定义需要禁用和启用的功能；
  - ActionType.upload
  - ActionType.download
  - ActionType.newFolder
  - ActionType.edit
  - ActionType.delete
  - ActionType.moveTo
  - ActionType.copyTo
  - ActionType.rename
  - ActionType.compress
  - ActionType.decompress
  - ActionType.search
  - ActionType.customActions

## Suite

Suite 是上述基本部件的集成，提供了完整的文件系统功能。你也可以根据上述部件自定义 Suite。接收如下 props：

- **win**: 当前窗口实例；
- **points**: 需要启用的文件管理节点；
- **defaultPath**: 默认显示文件，可以用于其他组件打开指定文件夹的 case；
- **defaultPoint**: 默认选择的文件节点；
- **showMenu**: 是否显示文件菜单；
- **showFavorite**: 是否启用收藏夹功能；
- **toolbar**: 自定义工具栏；
- **selectedKeys\$**: 当前选中的文件键值；
