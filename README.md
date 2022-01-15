## For what
我的想法：
1. 部署一个私有的tg api server到国外的vps
2. 通过 sshfs(sftp) 挂载我国内的nas到vps上
3. 在tg上转发视频到我的bot
4. 通过此程序handle此转发
5. 透过私有api server会下载视频到我的vps上
6. 将视频移动到挂载目录，自动同步到我的nas
7. done

以上愿景目前已基本实现，同时此仓库还保留了bot的扩展性，以后有其他想法如：上传图片返回链接(当图床用)、上传网页生成pdf等，可以直接在 bot目录下添加扩展

## 使用方法
1. 找botFarther申请bot，获取token
2. 将token放入config.toml中
3. [去这里](https://core.telegram.org/api/obtaining_api_id) 申请自己的api server，然后使用[api服务器](https://github.com/go-telegram-bot-api/telegram-bot-api)搭建
4. 修改相关的配置文件
5. 部署程序到可达tg的服务器
6. 将config.example.toml改为config.toml，然后运行程序
7. 开始转发视频即可

## TODO
1. 编写一个build脚本，自动生成run包，一键部署到vps
2. 尝试与chevereto图床联动

## Thanks
1. 生成uuid的库：https://github.com/google/uuid
2. 打印日志的库：https://github.com/sirupsen/logrus
3. 美化日志的库：https://github.com/antonfisher/nested-logrus-formatter
4. 处理配置文件的库：https://github.com/gookit/config
5. go操作telegram的api库：https://github.com/go-telegram-bot-api/telegram-bot-api
6. 协程池：https://github.com/panjf2000/ants


> **另外，由于官方api server的限制，这里我们需要搭建自己的[api服务器](https://github.com/go-telegram-bot-api/telegram-bot-api)**  
> 使用私有服务器能绕过官方服务器的下载限制  
