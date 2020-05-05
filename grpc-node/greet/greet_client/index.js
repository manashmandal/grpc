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

function main() {
  const client = new greetProto.GreetService(
    "0.0.0.0:50051",
    grpc.credentials.createInsecure()
  );

  client.Greet({ greeting: { first_name: "manash" } }, function (
    err,
    response
  ) {
    console.log("Greeting ", response.result);
    console.error(err);
  });
}

main();
