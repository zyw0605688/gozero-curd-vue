# gozero-curd-vue
gozero根据数据库模型，生成curd代码，生成vue增删改查页面

### 使用流程
1. 下载本项目，安装依赖
2. 编辑config.yaml文件，配置数据库连接信息，运行
3. 正常访问根路径即可看到hello world
4. post请求，生成代码（前提是你的go环境，goctl环境都已经好了）
```
curl --location --request POST 'http://localhost:8888/curd' \
--header 'Content-Type: application/json' \
--data-raw '{
    "model_name": "TdFirm",
    "only_gen_api": true
}'
```

### Curd是什么
一个低代码接口，只需定义model，即可自动生成增、删、改、详情、列表5个接口

### 工作流程
1. 反射model结构体，提取字段，根据字段及规则拼装字符串，生成.api文件
2. 将生成的.api文件名追加到goctl生成代码的入口index.api文件中
3. 调用goctl生成代码
4. 根据规则生成逻辑代码字符串，追加到logic代码末尾
5. 删除logic代码中多余的代码,将生成的vue文件和api文件放到前端项目里使用

### 注意事项
1. 定义model结构体时，主键字段放第一行
2. model结构体及名字，需要加入model.ModelList对象中，后续用来遍历反射结构体
3. 具体可查看curd logic源码
4. 前端项目，请留意修改反向代理
5. 可根据项目实际需求，修改模板和路径。

### 项目用到的第三方包，请提前安装
1. "gorm.io/gorm" 都知道
2. "github.com/jinzhu/copier" 用来拷贝参数
3. "github.com/go-cmd/cmd" 用来兼容执行shell命令

### 视频教程
后续补上，only_gen_api参数还未生效，后面有空再补充
如果有兴趣参与的小伙伴，也欢迎pr
走过路过的帮忙点个star，谢谢