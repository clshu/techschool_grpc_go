## Fix vscode-proto3 can find import problem.

- Go to vscode-proto3
- Click on example/.vscode link https://github.com/zxh0/vscode-proto3/blob/master/example/.vscode/settings.json

- Or copy/paste the example.
- Copy .vscode and modify it.

```
"protoc": {
    "path": "/usr/local/bin/protoc",
    "options": [
      "--proto_path=proto"
    ]
  }
```

- Preferences -> Setting
- Search "proto3"
- Edit setting.json
- Add json code above to "portoc" to setting.json.

## Install clang-format to format \*.proto

- brew install clang-format
