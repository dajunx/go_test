## 项目说明：

****

##### 用途
用作自己学习和测试go语言使用

##### 项目依赖管理工具
- 工具名称： **gvt**
- gvt 使用例子：
  - 1.gvt fetch github.com/xxxx  拉取依赖版本
  - 2.gvt delete  github.com/xxxx  delete a local dependency
  - 3.gvt list list dependencies one per line
- gvt github地址：https://github.com/FiloSottile/gvt
- gvt 程序已经放在了自己github仓库里面：https://github.com/dajunx/go_test/tree/master/bin
- go 守跨网原因，拉取不到的依赖 golang.org，已放置在 https://github.com/dajunx/3rd_lib 中，布置环境的时候，从该地方下载放在 $(GOPATH)/src/ 下面
