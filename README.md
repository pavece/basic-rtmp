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


## Sources

- RTMP spec (veriskope HTML version): [https://rtmp.veriskope.com/docs/spec/#status-of-this-document](https://rtmp.veriskope.com/docs/spec/#status-of-this-document)
- FLV & F4V spec: [https://rtmp.veriskope.com/pdf/video_file_format_spec_v10.pdf](https://rtmp.veriskope.com/pdf/video_file_format_spec_v10.pdf)
- Action Message Format (AMF0) spec (not implemented in this project): [https://rtmp.veriskope.com/pdf/amf0-file-format-specification.pdf](https://rtmp.veriskope.com/pdf/amf0-file-format-specification.pdf)

## Attribution

This project uses **[FFmpeg](https://ffmpeg.org/)** to convert the RTMP output (FLV format) into HLS and to transcode video into multiple resolutions. FFmpeg is only responsible for the media processing steps, segmenting the video and generating playlists. FFmpeg is **not used for the RTMP streaming itself**, which is implemented entirely in Go. FFmpeg is licensed under the **LGPL/GPL**; see [FFmpeg License](https://ffmpeg.org/legal.html) for details.
