package hls

import (
	"context"
	"time"

	"github.com/aler9/gortsplib"
	"github.com/aler9/gortsplib/pkg/h264"
	"github.com/asticode/go-astits"
)

const (
	mpegtsSegmentMinAUCount = 100
)

type writerFunc func(p []byte) (int, error)

func (f writerFunc) Write(p []byte) (int, error) {
	return f(p)
}

type muxerVariantMPEGTSSegmenter struct {
	segmentDuration time.Duration
	segmentMaxSize  uint64
	videoTrack      *gortsplib.TrackH264
	audioTrack      *gortsplib.TrackAAC
	onSegmentReady  func(*muxerVariantMPEGTSSegment)

<<<<<<< HEAD
	writer            *astits.Muxer
	currentSegment    *muxerVariantMPEGTSSegment
	videoDTSExtractor *h264.DTSExtractor
	startPCR          time.Time
	startDTS          time.Duration
=======
	writer         *astits.Muxer
	currentSegment *muxerVariantMPEGTSSegment
	videoDTSEst    *h264.DTSEstimator
	startPCR       time.Time
	startPTS       time.Duration
>>>>>>> dahua
}

func newMuxerVariantMPEGTSSegmenter(
	segmentDuration time.Duration,
	segmentMaxSize uint64,
	videoTrack *gortsplib.TrackH264,
	audioTrack *gortsplib.TrackAAC,
	onSegmentReady func(*muxerVariantMPEGTSSegment),
) *muxerVariantMPEGTSSegmenter {
	m := &muxerVariantMPEGTSSegmenter{
		segmentDuration: segmentDuration,
		segmentMaxSize:  segmentMaxSize,
		videoTrack:      videoTrack,
		audioTrack:      audioTrack,
		onSegmentReady:  onSegmentReady,
	}

	m.writer = astits.NewMuxer(
		context.Background(),
		writerFunc(func(p []byte) (int, error) {
			return m.currentSegment.write(p)
		}))

	if videoTrack != nil {
		m.writer.AddElementaryStream(astits.PMTElementaryStream{
			ElementaryPID: 256,
			StreamType:    astits.StreamTypeH264Video,
		})
	}

	if audioTrack != nil {
		m.writer.AddElementaryStream(astits.PMTElementaryStream{
			ElementaryPID: 257,
			StreamType:    astits.StreamTypeAACAudio,
		})
	}

	if videoTrack != nil {
		m.writer.SetPCRPID(256)
	} else {
		m.writer.SetPCRPID(257)
	}

	return m
}

func (m *muxerVariantMPEGTSSegmenter) writeH264(pts time.Duration, nalus [][]byte) error {
	now := time.Now()
	idrPresent := h264.IDRPresent(nalus)
<<<<<<< HEAD
	var dts time.Duration
=======
>>>>>>> dahua

	if m.currentSegment == nil {
		// skip groups silently until we find one with a IDR
		if !idrPresent {
			return nil
		}

<<<<<<< HEAD
		m.videoDTSExtractor = h264.NewDTSExtractor()

		var err error
		dts, err = m.videoDTSExtractor.Extract(nalus, pts)
		if err != nil {
			return err
		}

		m.startPCR = now
		m.startDTS = dts
		dts = 0
		pts -= m.startDTS

		// create first segment
		m.currentSegment = newMuxerVariantMPEGTSSegment(now, m.segmentMaxSize,
			m.videoTrack, m.audioTrack, m.writer.WriteData)
	} else {
		var err error
		dts, err = m.videoDTSExtractor.Extract(nalus, pts)
		if err != nil {
			return err
		}

		dts -= m.startDTS
		pts -= m.startDTS

		// switch segment
		if idrPresent &&
			(dts-*m.currentSegment.startDTS) >= m.segmentDuration {
			m.currentSegment.endDTS = dts
=======
		// create first segment
		m.currentSegment = newMuxerVariantMPEGTSSegment(now, m.segmentMaxSize,
			m.videoTrack, m.audioTrack, m.writer.WriteData)
		m.startPCR = now
		m.videoDTSEst = h264.NewDTSEstimator()
		m.startPTS = pts
		pts = 0
	} else {
		pts -= m.startPTS

		// switch segment
		if idrPresent &&
			m.currentSegment.startPTS != nil &&
			(pts-*m.currentSegment.startPTS) >= m.segmentDuration {
			m.currentSegment.endPTS = pts
>>>>>>> dahua
			m.onSegmentReady(m.currentSegment)
			m.currentSegment = newMuxerVariantMPEGTSSegment(now, m.segmentMaxSize,
				m.videoTrack, m.audioTrack, m.writer.WriteData)
		}
	}

<<<<<<< HEAD
	err := m.currentSegment.writeH264(
		now.Sub(m.startPCR),
		dts,
		pts,
		idrPresent,
		nalus)
	if err != nil {
=======
	dts := m.videoDTSEst.Feed(pts)

	err := m.currentSegment.writeH264(now.Sub(m.startPCR), dts,
		pts, idrPresent, nalus)
	if err != nil {
		if m.currentSegment.buf.Len() > 0 {
			m.onSegmentReady(m.currentSegment)
		}
		m.currentSegment = nil
>>>>>>> dahua
		return err
	}

	return nil
}

func (m *muxerVariantMPEGTSSegmenter) writeAAC(pts time.Duration, aus [][]byte) error {
	now := time.Now()

	if m.videoTrack == nil {
		if m.currentSegment == nil {
<<<<<<< HEAD
			m.startPCR = now
			m.startDTS = pts
			pts = 0

			// create first segment
			m.currentSegment = newMuxerVariantMPEGTSSegment(now, m.segmentMaxSize,
				m.videoTrack, m.audioTrack, m.writer.WriteData)
		} else {
			pts -= m.startDTS

			// switch segment
			if m.currentSegment.audioAUCount >= mpegtsSegmentMinAUCount &&
				(pts-*m.currentSegment.startDTS) >= m.segmentDuration {
				m.currentSegment.endDTS = pts
=======
			// create first segment
			m.currentSegment = newMuxerVariantMPEGTSSegment(now, m.segmentMaxSize,
				m.videoTrack, m.audioTrack, m.writer.WriteData)
			m.startPCR = now
			m.startPTS = pts
			pts = 0
		} else {
			pts -= m.startPTS

			// switch segment
			if m.currentSegment.audioAUCount >= mpegtsSegmentMinAUCount &&
				m.currentSegment.startPTS != nil &&
				(pts-*m.currentSegment.startPTS) >= m.segmentDuration {
				m.currentSegment.endPTS = pts
>>>>>>> dahua
				m.onSegmentReady(m.currentSegment)
				m.currentSegment = newMuxerVariantMPEGTSSegment(now, m.segmentMaxSize,
					m.videoTrack, m.audioTrack, m.writer.WriteData)
			}
		}
	} else {
		// wait for the video track
		if m.currentSegment == nil {
			return nil
		}

<<<<<<< HEAD
		pts -= m.startDTS
=======
		pts -= m.startPTS
>>>>>>> dahua
	}

	err := m.currentSegment.writeAAC(now.Sub(m.startPCR), pts, aus)
	if err != nil {
<<<<<<< HEAD
=======
		if m.currentSegment.buf.Len() > 0 {
			m.onSegmentReady(m.currentSegment)
		}
		m.currentSegment = nil
>>>>>>> dahua
		return err
	}

	return nil
}
