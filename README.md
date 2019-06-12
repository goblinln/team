# Team

小团队协作平台（任务管理系统）

![预览](./Preview.png)

## 重要声明

1. **2.0后的版本为完全重构版，与之前lua服务器版不再兼容。不要尝试直接升级！！！**
2. 原`lua`版不再维护。

## 更新日志

* v2.0

  1. 完成通知功能
  2. 修复多处BUG
  3. 多处显示布局修改，增加【加载中】提示
  4. 服务器可配置监听端口，但该操作不支持可视化配置且需要重启

* v2.0b1
  
  1. 前端使用`React`重构，解决之前项目代码混乱的问题
  2. 后台使用`golang`重构，解决跨平台编译问题
  3. 彻底分离前后端，不再使用服务器渲染模式

## 实现功能

+ [x] 可视化配置部署
+ [x] 系统管理
    - [x] 帐号管理
    - [x] 项目管理
+ [x] 个人信息
    - [x] 修改
    - [x] 通知
+ [x] 任务管理
    - [x] 发布任务
    - [x] 任务流
    - [x] 看板
    - [x] 甘特图
    - [x] 过滤
    - [x] 评论
    - [x] 事件回顾
+ [x] 项目管理
    - [x] 人员配置
    - [x] 分支
    - [x] 周报
    - [x] 项目任务
+ [x] 文档
+ [x] 文件分享

## 使用说明

1. 对于不需要修改原码的用户，可直接从[发行版](https://gitee.com/love_linger/Team/releases)页面中下载编译好的可执行文件

2. 对于有需求修改原码的用户，修改完后可按下面的步骤自行编译。  

    2.1 环境

    * Go 1.12+  
    * Node.js
    * Git  

    2.2 前端

    ```shell
    # 拉取依赖
    cd frontend
    npm install

    # 编译生成 app.bundle.js
    npm run build
    ```

    2.3 后端

    ```shell
    # 编译生成可执行文件
    cd backend
    go build

    # 如果未安装过go.rice请先执行下面这一步，并确保GOPATH/bin加入PATH环境变量中
    go get github.com/GeertJohan/go.rice/rice

    # 将 frontend/dist/ 中的前端资源打包入可执行文件中，windows版中需要加.exe后缀
    rice append --exec team.exe
    ```

3. 服务器运行：team。默认端口8080。【注】该版本已内置部署功能，初次访问会进行配置。







