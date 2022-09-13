package com.example.grpcclientjava.controller;

import com.example.grpcclientjava.GrpcClientService;
import lombok.RequiredArgsConstructor;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

@RestController
@RequiredArgsConstructor
public class GrpcClientController {

    private final GrpcClientService grpcClientService;

//    @RequestMapping("/test")
//    @GetMapping
//    public String printMessage() {
//        return grpcClientService.sendMessage("test");
//    }

    @RequestMapping("/moct")
    @GetMapping
    public String printMoctRsp() {
        return grpcClientService.getMoctBounds();
    }

}
