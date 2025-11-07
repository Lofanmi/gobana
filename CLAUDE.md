# CLAUDE.md

## 前言

- 请你使用简体中文，不要使用其它语言
- 在我未同意之前, 你不能提交代码，不能推送到远程仓库
- 不需要生成单元测试、实例、README.md，除非我要求你生成
- 我要求生成单元测试，一般不需要生成示例和文档，除非我要求你生成
- 项目存在 Makefile，编译项目使用 `make build`，而不是 `go build`；依赖更新使用 `make wire`。

## 特殊的代码风格和提交规范

#### 注释需要包含标识符，和Go标准库的文档一样

如 `// String xxx...` 对应方法 `String() string`。

#### 函数返回时，尽量减少显式的结构体初始化，更加简洁，并注意错误处理。

```go
type Data struct{
	A string
}
func A() (res Data, err error) {
	return // 这样写比较简洁
	// 不要这样写
	// return Data{}, errors.New("error")
	// return Data{}, nil
}
```

#### git commit 格式： feat|fix|docs(xxx): message

1. feat(功能点或者是模块点，英文): 功能点是什么
2. fix(修复点或者是模块点，英文): 修复了什么
3. feat(refactor): 重构了什么
注：feat|fix|docs，必须是这三者之一，必须带括号说明改动点。