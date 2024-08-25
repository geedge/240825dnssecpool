## 本代码对现有进行DNSSEC验证的解析器池进行探测
### 输入
所有活跃的Do53入口（四十多万个）
### 输出
支持DNSSEC验证的部分
### 步骤：
1.首先在自己权威服务器配置资源记录
    dnssec-true.ncache.site     A   1.2.3.4
    dnssec-false.ncache.site    A   1.2.3.4
2.对区文件签名生成db.ncache.site.signed
3.修改后者的RRSIG记录，随便改一个字母
4.restart服务，在服务器上跑着
5.在客户端（即该服务器）发送对两个域名的查询
    先对所有服务器发送一遍true的查询，看看有没有ad?必要吗
    筛掉所有没有ad的，输出列表
    再对所有剩下的服务器发送一遍false的查询，看看是否返回fail响应

涉及的技术：
多线程发起DNS查询