protoc MobileSuite.proto --go_out=../chat_server/protocol/
protogen -i:MobileSuite.proto -o:../csharp_client/connectToGoServer/MobileSuiteModel.cs