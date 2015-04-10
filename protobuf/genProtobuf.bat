protoc MobileSuite.proto --go_out=../go_server/protocol/
protogen -i:MobileSuite.proto -o:../csharp_client/connectToGoServer/MobileSuiteModel.cs