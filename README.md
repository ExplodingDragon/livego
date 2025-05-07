
项目克隆自 https://github.com/gwuhaolin/livego , 并经过部分更改，仅供自用

<p align='center'>
    <img src='./logo.png' width='200px' height='80px'/>
</p>

简单高效的直播服务器：
- 安装和使用非常简单；
- 纯 Golang 编写，性能高，跨平台；
- 支持常用的传输协议、文件格式、编码格式；

#### 支持的传输协议
- RTMP
- AMF
- HLS
- HTTP-FLV

#### 支持的容器格式
- FLV
- TS

#### 支持的编码格式
- H264
- AAC
- MP3


## 使用
1. 启动服务：执行 `livego` 二进制文件启动 livego 服务, 在配置文件中指定房间和推流码
2. 推流: 通过`RTMP`协议推送视频流到地址 `rtmp://localhost:1935/{appname}/{channelkey}`, 例如： 使用 `ffmpeg -re -i demo.flv -c copy -f flv rtmp://localhost:1935/{appname}/{channelkey}` 推流([下载demo flv](https://s3plus.meituan.net/v1/mss_7e425c4d9dcb4bb4918bbfa2779e6de1/mpack/default/demo.flv));
3. 播放: 支持多种播放协议，播放地址如下:
    - `RTMP`:`rtmp://localhost:1935/{appname}/live`
    - `FLV`:`http://127.0.0.1:7001/flv/{appname}/live.flv`
    - `HLS`:`http://127.0.0.1:7002/hls/{appname}/live.m3u8`

所有配置项: 
```bash
Usage of ./livego:
  -conf string
        config path (default "livego.yaml")
```
