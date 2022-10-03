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

# Lecture #10.1

## Build Error

```
var laptopServer *service.LaptopServer
cannot use laptopServer (variable of type *service.LaptopServer) as pb.LaptopServiceServer value in argument to pb.RegisterLaptopServiceServer: missing method mustEmbedUnimplementedLaptopServiceServer

```

### Solution

- Add pb.Unimplemented\*Server to Server struct, which is generated by codegen plugin.
- e.g.

```
type LaptopServer struct {
	Store LaptopStore
	pb.UnimplementedLaptopServiceServer
}
```

# Lecture #10.2

## Preps

- Copied laptop_service.proto
- Find Java Annotation API. Click Gradle tab. Pick latest
- https://mvnrepository.com/artifact/javax.annotation/javax.annotation-api/1.3.2
- Update build.gradle

## Create Laptop service server code in Java.

    1. Add Java Annotation dependency to build.gradle.
    2. Copy laptop_service.proto from golang side.
    3. Codegen on laptop_service.proto.
    4. Create LaptopService class to implemet CreateLaptop.
    5. Create LaptopStore Interface and InMemoryLaptopStore class to store data.
    6. Create LaptopServer class to run gRPC server.
    7. Create AlreadyExistException class.

# Lecture #11.1

    1. Create laptop_filter_proto for the search filter.
    2. Implement LaptopStore.Search.
    3. Implement LaptopServer.SearchLaptop
    4. Add context code to detect cancel and deadline exceeds conditions.
    5. LaptopStore.Search uses a callback function as one of the arguments and LaptopServer.SearchLaptop passes a streaming function as the callback to LaptopStore.Search.

# Lecture #11.2

1. Add void Search(LaptopFilter filter, LaptopStream stream); to LaptopStore Interface.

```
public interface LaptopStore {
    // It could be a db, in memory store for now
    void Save(Laptop laptop) throws Exception;
    Laptop Find(String id);
    void Search(LaptopFilter filter, LaptopStream stream);
}
```

2. Create a new Interface LaptopStream with void Send(Laptop laptop);. This will be a callback and implemented by the caller of LaptopStore.Search().

```
public interface LaptopStream {
    void Send(Laptop laptop);
}
```

3. In InMemoryLaptopStore.Search(filter, stream),

- It iterates all entries in data.
- When the entry matches the filter, it calls stream.Send(laptop) to send out laptop entry.
- The stream is passed down by the caller.

```
@Override
    public void Search(LaptopFilter filter, LaptopStream stream) {
        for (Map.Entry<String, Laptop> entry: data.entrySet()) {
            Laptop laptop = entry.getValue();
            if (isQualified(filter, laptop)) {
                stream.Send(laptop.toBuilder().build());
            }
        }
    }

    private boolean isQualified(LaptopFilter filter, Laptop laptop) {
        if (laptop.getPriceUsd() > filter.getMaxPriceUsd()) {
            return false;
        }
        if (laptop.getCpu().getNumCores() < filter.getMinCpuCores()) {
            return false;
        }
        if (laptop.getCpu().getMinGhz() < filter.getMinCpuGhz()) {
            return false;
        }
        if (toBit(laptop.getRam()) < toBit(filter.getMinRam())) {
            return false;
        }

        return true;
    }

    private long toBit(Memory ram) {
        long value = ram.getValue();
        switch (ram.getUnit()) {
            case BIT:
                return  value;
            case BYTE:
                return value << 3; // 8 bits = 2^3 BIT
            case KILLOBYTE:
                return value << 13; // 8 * 1024 = 2^13 BIT
            case MEGABYTE:
                return value << 23; // 8 * 1024 * 1024 = 2^23 BIT
            case GIGABYTE:
                return value << 33; // 8 * 1024 * 1024 * 1024 = 2^33 BIT
            case TERABYTE:
                return value << 43; // 8 * 1024 * 1024 * 1024 * 1024 = 2^43 BIT
            default:
                return 0;
        }
    }
```

4. In LaptopService.searchLaptop(request, responseStreamObserver),

```
public void searchLaptop(SearchLaptopRequest request, StreamObserver<SearchLaptopResponse> responseStreamObserver) {
        LaptopFilter filter = request.getFilter();
        logger.info("get a search-laptop request with filter:\n" + filter);

        store.Search(filter, new LaptopStream() {
            @Override
            public void Send(Laptop laptop) {
                logger.info("found laptop with ID: " + laptop.getId());
                SearchLaptopResponse response = SearchLaptopResponse.newBuilder().setLaptop(laptop).build();
                responseStreamObserver.onNext(response);
            }
        });

        responseStreamObserver.onCompleted();
        logger.info("search laptop completed");
    }
```

5. In LaptopClient

- New methods

```
    private void SearchLaptop(LaptopFilter filter) {
        logger.info("search started");

        SearchLaptopRequest request = SearchLaptopRequest.newBuilder().setFilter(filter).build();
        Iterator<SearchLaptopResponse> responseIterator = blockingStub.searchLaptop(request);
        while (responseIterator.hasNext()) {
            SearchLaptopResponse response = responseIterator.next();
            Laptop laptop = response.getLaptop();
            logLaptop(laptop);
        }
        logger.info("search completed");
    }

    private void logLaptop(Laptop laptop) {
        logger.info("_ found: " + laptop.getId());
        // May log more info later
    }

```

- In main()

```
           for (int i = 0; i < 10; i++) {
                Laptop laptop = generator.NewLaptop();
                client.createLaptop(laptop);
            }

            Memory minRam = Memory.newBuilder()
                    .setValue(8)
                    .setUnit(Memory.Unit.GIGABYTE)
                    .build();

            LaptopFilter filter = LaptopFilter.newBuilder()
                    .setMaxPriceUsd(3000)
                    .setMinCpuCores(4)
                    .setMinCpuGhz(2.5)
                    .setMinRam(minRam)
                    .build();

            client.SearchLaptop(filter);
```
