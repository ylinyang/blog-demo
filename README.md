# blog-demo



1. 用户模块
   生成用户ID，采用雪花算法 可按照时间排序检索
   参数校验，使用validator库进行校验







# gin框架
1. 获取入参的方式 
- querystring，的是URL中?后面携带的参数，例如：/user/search?username=小王子&address=沙河。 获取请求的querystring参数
- 获取form参数，当前端请求的数据通过form表单提交时，例如向/user/search发送一个POST请求
- 获取json参数，当前端请求的数据通过JSON提交时
- 参数绑定，为能够更方便的获取请求的相关参数，可以基于请求的Content-Type识别请求数据类型并利用反射机制自动提取请求中QueryString、form表单、JSON、XML，并把值绑定到指定的结构体对象

# 数据库 
1. 使用sqlx1
