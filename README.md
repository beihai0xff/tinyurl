## tinyUrl

tinyUrl 是一个使用 Go 语言与 Base36 编码实现的短链接服务，并且支持 URL 替换功能，即可以将`https://www.amazon.com/%E6%9C%9D%E8%8A%B1%E5%A4%95%E6%8B%BE-%E9%B2%81%E8%BF%85/dp/7519015432`这样的长链接转换为类似于`http://localhost/t/6c7f`的短链接。

当短链接映射的网址发生改变，例如域名或协议发生变化时，我们可以更新短链接映射的网址，从而实现短链接的长期有效。

tinyUrl 默认数据一写入就落盘存储，相关内容在 [storage](https://github.com/wingsxdu/tinyurl/tree/master/storage) 模块中实现。

> 这还是一个早期版本，目前已在作者的博客中搭建并提供基础功能。

#### 示例

作者编写了一个使用 echo 框架进行演示的示例 [example/echo](https://github.com/wingsxdu/tinyurl/tree/master/example/echo)