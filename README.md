# Team

小团队协作平台（任务管理系统）。**IE等旧浏览器不支持。推荐Chrome**  
如果有问题，请提ISSUE。因为这个是业余时间在搞，可能不会及时回复。但我是无法忍受有未解决的ISSUE的！

![预览](./screenshot.png)

## 特色

1. 部署简单，只需要下载可执行文件及提供一个可用的MySQL服务。提供可视化参数配置。
2. 集成任务流管理、团队文档、文件分享等常用协作功能。
3. 代码简单，外部依赖少，二次开发门槛低。

## 旧版本升级说明

4.x版本由于重构了后台，并增加了里程碑模块，导致数据库与3.x版本不兼容。  
**请不要直接使用3.x的数据库！！！**

## 使用说明

1. [发行版](https://gitee.com/love_linger/Team/releases)中提供Windows与Linux的可执行文件。MacOS用户需要按2说明，自行编译

2. 自行编译说明。  

    2.1 环境

    * Go 1.12+  
    * Node.js
    * Git  

    2.2 编译生成可执行文件

    ```shell
    # 第一步生成前端JS代码
    cd view
    npm install
    npm run build

    # 第二步生成可执行文件
    cd ..
    go build

    # 第三步使用Go.Rice将资源文件打包入可执行文件中，如果不打入包中，需要将view/dist/目录也放入部署环境
    # 【注1】Go.Rice的安装方式`go get github.com/GeertJohan/go.rice/rice`
    # 【注2】windows下`--exec`后面的参数需要加上.exe后缀
    rice append --exec team
    ```

3. 运行team可执行文件，访问 http://localhost:8080 进行配置







