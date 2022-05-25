# Scan UDP扫描器

> 原理：通过抓取ICMP回包，判断UDP端口是否可达

> 使用

```shell
> go get github.com/comeonjy/scan
> go install github.com/comeonjy/scan@main
> scan u 123.123.312.321:30001
```