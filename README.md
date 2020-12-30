# clang-formatter
clang-formatter是一个通过调用clang-format进行代码批量格式化的工具，使用时根据配置项生成clang-format命令并执行。<br>
使用该工具前需要先安装clang-format并将其加入环境变量。<br>
LLVM下载地址: https://releases.llvm.org/ <br><br>
**注意：不同版本的clang-format对于同一style在某些地方也会生成不同风格！团队成员间请使用同样的clang-format版本，或者使用.clang-format文件限定代码风格！**

## 配置文件说明
```
{
  "style": "WebKit",
  "filter": [
    "*.c",
    "*.cc",
    "*.cpp",
    "*.h",
    "*.hh",
    "*.hpp"
  ],
  "dirs": [
    "./"
  ]
}
```
`style` 用于指定format风格，可选值有 `LLVM, GNU, Google, Chromium, Microsoft, Mozilla, WebKit, file` ，当该参数设置为 `file` 时，将会读取当前目录下的 `.clang-format` 文件作为代码风格配置，该文件可使用clang-format生成。<br>
`filter` 填入需要格式化的文件类型列表，符合条件的文件将会被clang-format格式化。<br>
`dirs` 填入需要格式化的目录列表，程序运行时会遍历递归遍历子目录，索引出符合 `filter` 列表中的文件。

## 使用说明
直接运行程序将会在当前目录下生成一个 `format.json` 文件

```
$ ./clang-formatter
param config is empty, create a example config...
example config is created!
```

根据上一节的内容编写好配置文件后，运行程序并使用 `-c` 参数引入配置文件，即可批量格式化代码。

```
$ ./clang-formatter -c ./format.json
```