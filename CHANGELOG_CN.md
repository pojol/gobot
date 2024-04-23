# v0.4.4
> 里程碑版本, 重构代码使其更易阅读
* Feature
    - 目录移动，将 driver 相关的代码移动到 driver 目录，根目录只保留 editor & driver 两个目录
    - 目录移动，将原本 sample 目录中的机器人运行文件移动到 /mock/bot_sample_files 中
    - 修复 editor 中的选中节点的问题，目前点击空白处会正常将选中节点重制
    - 节点偏移精细控制，当节点处于选中状态，按键 [up, down, left, right] 会将该节点偏移一像素
    - 连线孔设置，当前鼠标移动到节点时，会放大连线孔的大小，当鼠标单击空白区域时会还原回去
    - 启动命令调整，新的 -h 更加易读，并提供了更改 port 的能力
    - 更详细的文档
    - 引入日志库，完善程序中的日志输出
    - 将脚本层原 meta 结构更改为 bot（更为直观的语意
    - debug逻辑调整，当前提供更多的运行时调试能力

# v0.4.3
* Features
    - 将 meta 面板的命名调整为 blackboard，更贴近行为树的语意 #19
    - 将延迟显示label的长度固定，避免在不同的延迟下影响其他控件的排列显示 #15
    - 将原有的memory类型sqlite，修改为文件型的sqlite（避免重启丢失机器人文件 #17
* Fixes
    - 修复心跳检查延迟的重复构建错误（会导致过快的刷新 #18

# v0.4.2
* Features 
    - 提供了 message 模块，用户现在可以在脚本层自行对 stream byte data 进行处理（拆包，封包
    - 修改了 report 的概念，不再提供请求耗时等信息（交由后台去统计更合理）report 目前仅提供，req,res,ntf等次数维度的统计

# v0.4.1
* Features
    - 添加了 websocket 模块
* Fixes
    - 修复丢失的 banner 打印

# v0.4.0
* Features
    - 添加了集群部署模式

# v0.3.6 (pre
* Features 
    - 将 report 中的预览方式从通过点击 tag 的形式换为直接显示在下方，通过tab 进行图表的切换（更直观
    - 替换了 codemirror 的实现库，使代码的编写体验更好
    - 添加 share 功能，在 bots 面板中选中 bot 点击 share 能将 bot 的地址复制在剪贴板中，别人可以直接访问这个地址打开 bot 的编辑视窗
    - 添加了 running 的自动刷新（默认 10s
    - 更换了 batch 的存储实现，现在将存储在 db 中，以便于在异常中断后还能继续执行

* Fixes
    - 如果数据库无法连接则会直接panic（遇到错误应该立即终止
    - share 按钮，剪贴板的实现替换为更早的 api（能适配更多的浏览器
    - 解决的 report 没有按时间排序的问题
    - 解决了 bots 中控件点击事件错乱的问题

# v0.3.5 (pre
* Features
    - 添加了一个入队延迟的配置，用于控制机器人的调度频率
    - 优化了 sideplane 中节点的 css 实现
    - 添加了 http query params 形式的输入
* Fixes
    - 代码输入框在切换输入法后输入逻辑错乱
    - bots 中点击 inputnumber 会丢失焦点的问题
    - 缩放，调整视窗时编辑器窗口没有等比放大缩小


# v0.3.1 (pre
* Features
    - 添加擦除行为树的按钮
    - graph 的绘制现在完全居于 model/tree 中的数据进行
* Fixes
    - 点击过快导致当前节点绘制出现错误
    - 调试窗口的一些跳转进行修复
    - 缩放，调整视窗时编辑器窗口没有等比放大缩小

# v0.3.0 (pre
* Features
    - 重构整个editor使用umi框架和ts进行编写（类型安全，支持暗黑模式切换
    - 将组件更换为function，全面使用 hooks + redux 的方式编写代码（无状态的模式（优化了载入时间，绘制效率
    - 添加了新的bot加载方式，可以通过访问url的形式加载某个bot（更容易传播
    - 引入了内存数据库 sqlite (方便试用，可以快速的进行本地部署
* Fixes
    - 修复了尾部节点信息丢失
    - 修复了在压力测试过程中批次不能准确退出的问题

# v0.2.5
* Features
    - 重写侧边栏提供更好的筛选方式
    - 将 prefab 单独作为一个页，并提供搜索和编辑功能
    - 优化连接点（当没有鼠标 movein 的时候缩小
    - 添加 report 页的时间排序

## v0.2.1
* Features
    - 添加了新的并行节点
    - 删除原有的 assert 节点类型
    - 新增了 runtime err 栏用于输出运行时错误信息
    - 重构了 bot 的运行时逻辑
    - response 栏中引入了展示 thread 信息的逻辑（并行节点将创建新的 thread）
    - 运行到节点时添加了个小动画（提示优化

## v0.1.17
* Features
* Fixes 
    - 修复异步加载行为树的逻辑错误

## v0.1.16
* Features 
    - 为lua代码提供fmt功能
    - 为 prefab 留出足够的空间，将change移动到和meta窗口重叠显示
    - 添加 step 的快捷键 【 F10 】
* Fixes
    - 删除配置文件不能成功的错误
    - root节点位置修正
    - 修复初始化服务器地址后未进行健康检查

## v0.1.15
* Features
    - 去除 edit 界面中的 debug 按钮，在点击 step 的时候自动进行创建
    - 在最后一个节点上为 step 设置延迟（防止连续连击
    - 添加 reset 按钮防止用户中途不想执行下去时，必须要走到底
    - 去除 script module 中的 utils (改为独立的 uuid 和 random 接口显示在一级目录方便查阅
    - 添加节点`预制功能`（现在用户可以在config面板定义和复用自己的脚本节点
    - 添加连接状态提示
* Fixes 
    - step api 没有返回正确的错误信息

## v0.1.14
* Features
* Fixes
    - step 没有锁定，导致多次点击会加快播放速度
    - unlink 的时候传入的应该是 node id 而不是 cell id（会导致链路没有正真的断开
    - 脚本节点不应该可以挂接多个节点

## v0.1.13
* Fixes 
    - 位置信息没有被正确同步（原先只同步了当前节点

## v0.1.12
* Features
    - 添加了 botid 到 meta 数据结构中（每批次唯一依次递增 1,2,3
    - 添加了 mongodb 模块
    - 添加了 无数据库运行模式 （方便用户体验
    - 去掉了使用 tabs 来控制输入的切换（带来更大的可编辑区域

* Fixes
    - 代码框宽度的自适应
    - 改善了http的超时机制
    - 将一些无效的操作屏蔽（比如连接到根节点
    - 将react-split-pane替换成react-reflex，前者在resize事件触发时无法保持切分的比率（也没有接口可以重新计算出来

## v0.1.9
* Features
    - 添加了代码窗口的风格选择
    - 添加了step可选次数（现在可以在调试阶段更容易的定位到目标节点
    - 支持了 edit 和 report 界面中元素的中文显示
    - 报告将被保存在db中，而非内存
    - 报告中新增了后缀信息（标注这是 ms 或者 times
    - 为 debug 信息新增了行号
    - change 窗口更换为 #medium-editor 为了更好的支持文字作色
    - 添加了 loading 页
* Fixes
    - 

## v0.1.7
* Features
    - 将配置信息交给后台存储
    - 将配置界面和内部数据调整为可动态添加节点模版的表现方式
    - 一些界面逻辑的优化
* Fixes
    - 

## v0.1.6
* Features
    - 添加了 electron 打包模块
    - 修改的原来的节点返回值视图（现在可以捕捉到 lua vm 的错误
* Fixes
    - 修复了在运行批次机器人遇到错误时，不能继续往下执行的错误