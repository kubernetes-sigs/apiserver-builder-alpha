# Serializing resources in protobuf to backend storage

**TL;DR**: This document covers how to make your custom resources serialized 
into storages in protobuf instead of json format. To support protobuf serialization, 
we are making use of the standard [code-generator](https://github.com/kubernetes/code-generator/tree/master/cmd/go-to-protobuf) provided by kubernetes community
which helps converting field tags from the model to protobuf IDL source file with 
".proto" suffix and then generates marshalling extension methods easily from the protobuf IDL.
You are required to download the latest binaries which already includes the protobuf 
code-generator from the download page before walking through the following steps.

### Adding Protobuf Field Tags

> Note that you are supposed to have a working project with resources already 
initialized/implemented before reading on. 


We will be adding field tags with key `protobuf` onto the member fields of your resource 
definition so that the code-generator can recognize them and know how to squash them into valid
protobuf format. A valid format for the protobuf field tag will be `protobuf:"<type>,<index>,<modifier>,name=<name>"`
which is standing align w/ the kubernetes native resource definitions.


```go
type Foo struct {
	...
	Sample            SampleElem               `json:"sample,omitempty" protobuf:"bytes,2,opt,name=sample"`
	SamplePointer     *SamplePointerElem       `json:"sample_pointer,omitempty" protobuf:"bytes,3,opt,name=sample_pointer"`
	SampleList        []SampleListElem         `json:"sample_list,omitempty" protobuf:"bytes,4,rep,name=sample_list"`
	SamplePointerList []*SampleListPointerElem `json:"sample_pointer_list,omitempty" protobuf:"bytes,5,rep,name=sample_pointer_list"`
	SampleMap         map[string]SampleMapElem `json:"sample_map,omitempty" protobuf:"bytes,6,rep,name=sample_map"`
}
```

### Running Protobuf Code-Generation

As your resource already marked with protobuf field tags, we can run code-generation by:

```bash
apiserver-boot build generated --generator=protobuf
```

This command basically generates two files as a result: 

- `generated.proto`: Generated IDL from the protobuf field tags.
- `generated.pb.go`: Generated protobuf marshalling extension methods from the generated IDL.

If the code-generation completes happily, your resource model is now technically protobuf 
serializable. Also, we can verify by unit-testing at this point before moving on.


### Override Default Storage Media Type

Last but not least, we will have to override the content-type for the aggregated apiserver using the flag:

>  --storage-media-type string                               The media type to use to store objects in storage. Some resources or storage backends may only support a specific media type and will ignore this setting. (default "application/json")


Setting this flag to `application/vnd.kubernetes.protobuf` will make all the resource served by this aggregated
apiserver serialized into storage in protobuf format. To order to switch serialization format for a subset of your 
resources, please combine this document with custom REST implementation for certain resource.
