package com.example.grpcclientjava;

import com.google.protobuf.ProtocolStringList;
import io.grpc.Channel;
import io.grpc.ManagedChannelBuilder;
import io.grpc.StatusRuntimeException;
import net.devh.boot.grpc.client.inject.GrpcClient;
import org.springframework.stereotype.Service;

@Service
public class GrpcClientService {

    @GrpcClient("test")
//    private SimpleGrpc.SimpleBlockingStub simpleStub;

//    public String sendMessage(final String name) {
//        try {
//            HelloReply response = this.simpleStub.sayHello(HelloRequest.newBuilder().setName(name).build());
//            return response.getMessage();
//        } catch (StatusRuntimeException e) {
//            return "Failed with " + e.getStatus().getCode().name();
//        }
//    }

    Channel channel = ManagedChannelBuilder.forAddress("localhost", 1234).usePlaintext().build();
    MoctSvrGrpc.MoctSvrBlockingStub moctSimpleStub = MoctSvrGrpc.newBlockingStub(channel);

    public String getMoctBounds() {
        try{
            Moct.Bounds.Builder bounds = Moct.Bounds.newBuilder()
                    .setMin(Moct.Point.newBuilder().setLat(37.56786).setLon(126.98575))
                    .setMax(Moct.Point.newBuilder().setLat(37.56264).setLon(126.99570));
            Moct.IntersectsResponse rsp = this.moctSimpleStub.intersectsBounds(Moct.IntersectsBoundsRequest.newBuilder()
                    .setKey("link").setCursor(0).setLimit(0).setBounds(bounds).build());
            ProtocolStringList geoJsonObjectList = rsp.getObjectsList();
            return geoJsonObjectList.toString();
        } catch (StatusRuntimeException e) {
            return "Failed with " + e.getStatus().getCode().name();
        }
    }
}
