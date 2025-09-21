![RTMP2HLS Header Image](https://static.pavece.com/public-files/rtmp2hls/rtmp2hls-logo.png)

# RTMP2HLS

This project features a partial RTMP implementation from the [spec](https://rtmp.veriskope.com/docs/spec/#status-of-this-document). It's able to ingest live video from sources like OBS, transcode that video to live ABR [HLS](https://en.wikipedia.org/wiki/HTTP_Live_Streaming), upload those segments to an object store for CDN serving to end audience. In short: this project implements a **basic version** of the system used by platforms like Twitch for live streaming

> [!WARNING]  
> This is a toy learning project. It is unfinished, under active development and not intended for production use (much better alternatives out there). Use at your own risk.

## Features

- Custom partial RTMP implementation
- RTMP video ingestion compatible with OBS and FFmpeg
- ABR HLS list generation with configurable renditions
- Automatic live HLS (segments and list) uploading to object storage (S3/MinIO)
- Automatic DVR (VOD playback) playlist generation and uploading
- Customizable callbacks for stream key validation and stream lifecycle management
- Multi streamer setups

## Demos

This video shows how you can stream directly from OBS to this server and perform live playback using HLS in VLC (or any other HLS player).

[![OBS Demo](https://static.pavece.com/public-files/rtmp2hls/quick-demo.jpg)](https://static.pavece.com/public-files/rtmp2hls/quick-demo.mp4)

⚠️ Note: This is a work-in-progress demo. Expect ~10s latency between OBS and HLS playback.

### DVR demos

These videos where streamed from OBS to this server. These clips are extracted from the automatically generated DVR HLS list.

#### Audio / Video synchronization

[![Audio video sync demo](https://static.pavece.com/public-files/rtmp2hls/sync.jpg)](https://static.pavece.com/public-files/rtmp2hls/sync.mp4)

#### Big buck bunny fragment

[![Big buck bunny demo](https://static.pavece.com/public-files/rtmp2hls/big-buck-bunny.jpg)](https://static.pavece.com/public-files/rtmp2hls/big-buck-bunny.mp4)

## Running with Docker

You can clone the repo and launch the docker compose to build the server's image and launch additional services.

```bash
docker compose build
docker compose up -d
```

If using minio remember to configure the CDN bucket for **anonymous readonly access**.

## Sources

- RTMP spec (veriskope HTML version): [https://rtmp.veriskope.com/docs/spec/#status-of-this-document](https://rtmp.veriskope.com/docs/spec/#status-of-this-document)
- FLV & F4V spec: [https://rtmp.veriskope.com/pdf/video_file_format_spec_v10.pdf](https://rtmp.veriskope.com/pdf/video_file_format_spec_v10.pdf)
- Action Message Format (AMF0) spec (not implemented in this project): [https://rtmp.veriskope.com/pdf/amf0-file-format-specification.pdf](https://rtmp.veriskope.com/pdf/amf0-file-format-specification.pdf)

## Attribution

This project uses **[FFmpeg](https://ffmpeg.org/)** to convert the RTMP output (FLV format) into HLS and to transcode video into multiple resolutions. FFmpeg is only responsible for the media processing steps, segmenting the video and generating playlists. FFmpeg is **not used for the RTMP streaming itself**, which is implemented entirely in Go. FFmpeg is licensed under the **LGPL/GPL**; see [FFmpeg License](https://ffmpeg.org/legal.html) for details.
