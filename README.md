# Team

小团队协作平台（任务管理系统）

![预览](./Preview.png)

## 重构说明

1. 新版移除`Ant-Design`的依赖，改为使用自己实现的组件库。主要原因有：
    * Antd太大了，项目按需加载后，生成的app.js都2.5MB。
    * Antd的依赖过重
    * Form使用起来繁琐
    * Dropdown及Popover感觉有点反人类，为什么具体内容不是child而是attribute，而本可以做成label属性的反而是child
    * 想使用React HOOK实现一个通用的组件库，方便其他项目使用
2. 新版组件库的样式参考了`Ant-Design`与`layui`
3. 新版API不再使用PATCH方法，因为发现Edge等浏览器会有问题
4. 新版js使用es6，字体使用woff格式，IE等老旧浏览器不可用
5. 2.x与3.x数据兼容。

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

1. 对于不需要修改原码的用户，可直接从[发行版](https://gitee.com/love_linger/Team/releases)页面中下载编译好的可执行文件，放入`publish`目录中，然后直接第3步

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

    # 编译生成 publish/assets/app.js
    npm run build
    ```

    2.3 后端

    ```shell
    # 编译生成可执行文件
    cd backend
    go build -o ../publish/team.exe
    ```

3. 将 `publish` 目录下的文件拷贝到服务器部署路径

4. 运行：team。默认端口8080。【注】该版本已内置部署功能，初次访问会进行配置。







