const PROTO_PATH = "./greetpb/greet.proto";

const grpc = require("grpc");
const protoLoader = require("@grpc/proto-loader");

const packageDefinition = protoLoader.loadSync(PROTO_PATH, {
  keepCase: true,
  longs: String,
  enums: String,
  defaults: true,
  oneofs: true,
});

const greetProto = grpc.loadPackageDefinition(packageDefinition).greet;

// console.log(helloProto);

function Greet(call, callback) {
  console.log("request", call.request);
  callback(null, {
    result: "Hello " + call.request.greeting.first_name,
  });
}

function main() {
  const server = new grpc.Server();
  server.addService(greetProto.GreetService.service, { Greet: Greet });
  server.bind("0.0.0.0:50051", grpc.ServerCredentials.createInsecure());
  console.log("binded");
  server.start();
  console.log("Started");
}

main();
