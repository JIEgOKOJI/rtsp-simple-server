package hls

import (
	"bytes"
	"encoding/hex"
	"io"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/aler9/gortsplib"
)

type muxerPrimaryPlaylist struct {
	fmp4       bool
	videoTrack *gortsplib.TrackH264
	audioTrack *gortsplib.TrackAAC
	streamName string
}

func newMuxerPrimaryPlaylist(
	fmp4 bool,
	videoTrack *gortsplib.TrackH264,
	audioTrack *gortsplib.TrackAAC,
	streamName string,
) *muxerPrimaryPlaylist {
	return &muxerPrimaryPlaylist{
		fmp4:       fmp4,
		videoTrack: videoTrack,
		audioTrack: audioTrack,
		streamName: streamName,
	}
}
func (p *muxerPrimaryPlaylist) getPrem(streamname string) (int, error) {
	var Client = &http.Client{Timeout: 5 * time.Second}
	r, err := Client.Get("https://goodgame.ru/api/4/transcoding/" + streamname)
	if err != nil {
		return 0, err
	}
	defer r.Body.Close()
	api, _ := ioutil.ReadAll(r.Body)
	if string(api) == "true" {
		return 1, nil
	} else {
		return 0, nil
	}
}
func (p *muxerPrimaryPlaylist) file() *MuxerFileResponse {
	return &MuxerFileResponse{
		Status: http.StatusOK,
		Header: map[string]string{
			"Content-Type": `audio/mpegURL`,
		},
		Body: func() io.Reader {
			var codecs []string

			if p.videoTrack != nil {
				sps := p.videoTrack.SPS()
				if len(sps) >= 4 {
					codecs = append(codecs, "avc1."+hex.EncodeToString(sps[1:4]))
				}
			}

			// https://developer.mozilla.org/en-US/docs/Web/Media/Formats/codecs_parameter
			if p.audioTrack != nil {
				codecs = append(codecs, "mp4a.40."+strconv.FormatInt(int64(p.audioTrack.Type()), 10))
			}

			switch {
			case !p.fmp4:
				return bytes.NewReader([]byte("#EXTM3U\n" +
					"#EXT-X-VERSION:3\n" +
					"\n" +
					"#EXT-X-STREAM-INF:BANDWIDTH=200000,CODECS=\"" + strings.Join(codecs, ",") + "\"\n" +
					"stream.m3u8\n"))

			default:
				status, err := p.getPrem(p.streamName)
				if err != nil {
					return bytes.NewReader([]byte("#EXTM3U\n" +
						"#EXT-X-VERSION:9\n" +
						"#EXT-X-INDEPENDENT-SEGMENTS\n" +
						"\n" +
						"#EXT-X-STREAM-INF:BANDWIDTH=200000,CODECS=\"" + strings.Join(codecs, ",") + "\"\n" +
						"/lhls/" + p.streamName + "/stream.m3u8\n"))
				}
				if status == 1 {
					return bytes.NewReader([]byte("#EXTM3U\r\n" +
						"#EXT-X-VERSION:9\r\n" +
						"#EXT-X-INDEPENDENT-SEGMENTS\n" +
						"\n" +
						"#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=7155000,RESOLUTION=1920x1080\r\n" +
						"/lhls/" + p.streamName + "_prem/stream.m3u8\r\n" +
						"#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=2000000,RESOLUTION=1280x720\r\n" +
						"/lhls/" + p.streamName + "_720/stream.m3u8\r\n" +
						"#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=1320000,RESOLUTION=800x450\r\n" +
						"/lhls/" + p.streamName + "_480/stream.m3u8\r\n" +
						"#EXT-X-STREAM-INF:PROGRAM-ID=1,BANDWIDTH=296100,RESOLUTION=480x270\r\n" +
						"/lhls/" + p.streamName + "_240/stream.m3u8"))
				}
				return bytes.NewReader([]byte("#EXTM3U\n" +
					"#EXT-X-VERSION:7\n" +
					"\n" +
					"#EXT-X-STREAM-INF:BANDWIDTH=200000,CODECS=\"" + strings.Join(codecs, ",") + "\"\n" +
					"stream.m3u8\n" +
					"\n"))
			}
		}(),
	}
}
