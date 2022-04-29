## GO 爬虫

### 简单的GO爬虫脚本，仅供学习参考

#### 依赖服务
- 爬虫框架:colly (http://go-colly.org/docs/)
- redis
- GORM

#### 开发步骤
1. 分析目标站点数据
2. 编写采集策略
3. 采集数据，写入数据库

#### 涉及到的知识点
1. GO 基础用法
2. 文件和目录操作
3. colly 基础用法及HTML解析规则
4. 并发操作
5. 数据库GORM
6. Redis

#### 支持的功能
- [x] GameSky 动态图和囧图下载
- [ ] Go语言中文的文章采集 

#### 如何使用

```shell
go run main.go
```