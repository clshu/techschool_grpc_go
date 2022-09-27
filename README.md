# Lecture #7

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

# Lecture #8

## Get Java gradle plugin

- Search "protobuf gradle plugin"
- https://github.com/google/protobuf-gradle-plugin
- Find "Using the Gradle plugin DSL

```
plugins {
  id "com.google.protobuf" version "0.8.19"
  id "java"
}
```

- Copy 1st line in plugins
- In Projcet of IntelliJ, pcbook > src > build.gradle
- Append the copied line to plugins, order does not matter

## Get Maven protobug plugin

- Search "Mavern protobuf java"
- https://mvnrepository.com/artifact/com.google.protobuf/protobuf-java
- Click on the latest stable
- https://mvnrepository.com/artifact/com.google.protobuf/protobuf-java/3.21.6
- Click on Gradle tab
- Copy the setting

```
// https://mvnrepository.com/artifact/com.google.protobuf/protobuf-java
implementation group: 'com.google.protobuf', name: 'protobuf-java', version: '3.21.6'

```

- Add the setting to ependenciies of build.gradle

## Get Maven grpc-all plugin

- In MVN Repository search box, type "grpc-all"
- https://mvnrepository.com/artifact/io.grpc/grpc-all/1.49.1

```
// https://mvnrepository.com/artifact/io.grpc/grpc-all
implementation group: 'io.grpc', name: 'grpc-all', version: '1.49.1'
```

- Add setting to the dependencies of build.gradle

## Get protobuf compiler

### Get Compiler version

- In MVN Repository search box, type "protobuf compiler"
- Find https://mvnrepository.com/artifact/com.google.protobuf/protoc
- Remeber the latest version (3.21.6 at this time)

### Get protobuf executable setting

- Go to https://github.com/google/protobuf-gradle-plugin
- Find "Customize Protobuf compilation"

```
protobuf {
  // Configure the protoc executable
  protoc {
    // Download from repositories
    artifact = 'com.google.protobuf:protoc:3.21.6'
  }
}
```

- Copy and paste the whole block to build.gradle
- Replace the version with 3.21.6

## Tell protoc to use Java codegen plugin

- Go to https://github.com/google/protobuf-gradle-plugin
- Find "protobuf" block again under previous text

```
protobuf {
// Locate the codegen plugins
  plugins {
    // Locate a plugin with name 'grpc'. This step is optional.
    // If you don't locate it, protoc will try to use "protoc-gen-grpc" from
    // system search path.
    grpc {
      artifact = 'io.grpc:protoc-gen-grpc-java:1.49.1'
      // or
      // path = 'tools/protoc-gen-grpc-java'
    }
    // Any other plugins
  }
}
```

- Find "io.grpc:protoc-gen-grpc-java" in MVN repository. Replace the version.
- Add the previous plugins block to probuf block.

## Customize codegen task

- Add the following block to protobuf block.

```
generateProtoTasks {
  all*.plugins {
    grpc{}
  }
}
```

## Tell IntelliJ where to generate the code

- Create a top level block sourceSets {} in build.gradle

```
sourceSets {
  main {
    java {
      srcDirs 'build/generated/source/proto/main/grpc'
      srcDirs 'build/generated/source/proto/main/java'
    }
  }
}

```

# Lecture #9.1

## Error

- protoreflect.ProtoMessage does not implement protoiface.MessageV1 (missing ProtoMessage method)

## Solution

- https://github.com/golang/protobuf/issues/1133

- From dsnt

- There's a mismatch going on, since the code is trying to pass a "google.golang.org/protobuf/proto".Message to "github.com/golang/protobuf/jsonpb".Marshaler.MarshalToString, which expects a "github.com/golang/protobuf/proto".Message.

- This can be resolved by either:

1. Using the "google.golang.org/protobuf/encoding/protojson" package instead, which accepts the newer Message interface type you currently have on hand, OR
2. Use the "github.com/golang/protobuf/proto".MessageV1 adaptor function to convert the newer Message interface type you have to the legacy one.
   I recommend the first solution:

```
package serializer

import (
    "google.golang.org/protobuf/encoding/protojson"
    "google.golang.org/protobuf/proto"
)

func ProtobufToJSON(message proto.Message) (string, error) {
    b, err := protojson.MarshalOptions{
        Indent: true,
        UseProtoNames: true,
        EmitUnpopulated: true,
    }
    return string(b), err
}
```

### Reasons

- The lecture's protoc codegen is different from the latest protoc codegen.
- The latest protoc codegen generate the code importing from google.golang.org/protobuf/proto

- e.g.

```
import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)
```

# Lecture #10.1

## protoc works differently

- Need to add "--go-grpc_out=." to protoc arguments
- LaptopService is created in a different file
- 2 files created. laptop_service.pb.go and laptop_service_grpc.pb.go
- Install google.golang.org/grpc
