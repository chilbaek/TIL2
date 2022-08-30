package com.example.grpcserverjava;

import io.grpc.stub.StreamObserver;
import net.devh.boot.grpc.server.service.GrpcService;
import org.chb.examples.lib.HelloReply;
import org.chb.examples.lib.HelloRequest;
import org.chb.examples.lib.SimpleGrpc;


@GrpcService
public class GrpcServerService extends SimpleGrpc.SimpleImplBase{

    @Override
    public void sayHello(HelloRequest request, StreamObserver<HelloReply> responseObserver) {
        HelloReply reply = HelloReply.newBuilder()
                .setMessage("Hello ==> " + request.getName())
                .build();
        responseObserver.onNext(reply);
        responseObserver.onCompleted();
    }

}
