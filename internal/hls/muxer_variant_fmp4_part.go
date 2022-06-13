package hls

import (
	"bytes"
	"io"
	"math"
	"strconv"
	"time"

<<<<<<< HEAD
	gomp4 "github.com/abema/go-mp4"
	"github.com/aler9/gortsplib"
	"github.com/aler9/gortsplib/pkg/aac"

	"github.com/aler9/rtsp-simple-server/internal/mp4"
=======
	"github.com/abema/go-mp4"
	"github.com/aler9/gortsplib"
	"github.com/aler9/gortsplib/pkg/aac"
>>>>>>> dahua
)

func durationGoToMp4(v time.Duration, timescale time.Duration) int64 {
	return int64(math.Round(float64(v*timescale) / float64(time.Second)))
}

func mp4PartGenerateVideoTraf(
<<<<<<< HEAD
	w *mp4.Writer,
	trackID int,
	videoSamples []*fmp4VideoSample,
) (*gomp4.Trun, int, error) {
=======
	w *mp4Writer,
	trackID int,
	videoSamples []*fmp4VideoSample,
	startDTS time.Duration,
) (*mp4.Trun, int, error) {
>>>>>>> dahua
	/*
		traf
		- tfhd
		- tfdt
		- trun
	*/

<<<<<<< HEAD
	_, err := w.WriteBoxStart(&gomp4.Traf{}) // <traf>
=======
	_, err := w.writeBoxStart(&mp4.Traf{}) // <traf>
>>>>>>> dahua
	if err != nil {
		return nil, 0, err
	}

	flags := 0

<<<<<<< HEAD
	_, err = w.WriteBox(&gomp4.Tfhd{ // <tfhd/>
		FullBox: gomp4.FullBox{
=======
	_, err = w.writeBox(&mp4.Tfhd{ // <tfhd/>
		FullBox: mp4.FullBox{
>>>>>>> dahua
			Flags: [3]byte{2, byte(flags >> 8), byte(flags)},
		},
		TrackID: uint32(trackID),
	})
	if err != nil {
		return nil, 0, err
	}

<<<<<<< HEAD
	_, err = w.WriteBox(&gomp4.Tfdt{ // <tfdt/>
		FullBox: gomp4.FullBox{
			Version: 1,
		},
		// sum of decode durations of all earlier samples
		BaseMediaDecodeTimeV1: uint64(durationGoToMp4(videoSamples[0].dts, fmp4VideoTimescale)),
=======
	_, err = w.writeBox(&mp4.Tfdt{ // <tfdt/>
		FullBox: mp4.FullBox{
			Version: 1,
		},
		// sum of decode durations of all earlier samples
		BaseMediaDecodeTimeV1: uint64(durationGoToMp4(startDTS, fmp4VideoTimescale)),
>>>>>>> dahua
	})
	if err != nil {
		return nil, 0, err
	}

	flags = 0
	flags |= 0x01  // data offset present
	flags |= 0x100 // sample duration present
	flags |= 0x200 // sample size present
	flags |= 0x400 // sample flags present
	flags |= 0x800 // sample composition time offset present or v1

<<<<<<< HEAD
	trun := &gomp4.Trun{ // <trun/>
		FullBox: gomp4.FullBox{
=======
	trun := &mp4.Trun{ // <trun/>
		FullBox: mp4.FullBox{
>>>>>>> dahua
			Version: 1,
			Flags:   [3]byte{0, byte(flags >> 8), byte(flags)},
		},
		SampleCount: uint32(len(videoSamples)),
	}

	for _, e := range videoSamples {
		off := e.pts - e.dts

		flags := uint32(0)
		if !e.idrPresent {
			flags |= 1 << 16 // sample_is_non_sync_sample
		}

<<<<<<< HEAD
		trun.Entries = append(trun.Entries, gomp4.TrunEntry{
=======
		trun.Entries = append(trun.Entries, mp4.TrunEntry{
>>>>>>> dahua
			SampleDuration:                uint32(durationGoToMp4(e.duration(), fmp4VideoTimescale)),
			SampleSize:                    uint32(len(e.avcc)),
			SampleFlags:                   flags,
			SampleCompositionTimeOffsetV1: int32(durationGoToMp4(off, fmp4VideoTimescale)),
		})
	}

<<<<<<< HEAD
	trunOffset, err := w.WriteBox(trun)
=======
	trunOffset, err := w.writeBox(trun)
>>>>>>> dahua
	if err != nil {
		return nil, 0, err
	}

<<<<<<< HEAD
	err = w.WriteBoxEnd() // </traf>
=======
	err = w.writeBoxEnd() // </traf>
>>>>>>> dahua
	if err != nil {
		return nil, 0, err
	}

	return trun, trunOffset, nil
}

func mp4PartGenerateAudioTraf(
<<<<<<< HEAD
	w *mp4.Writer,
	trackID int,
	audioTrack *gortsplib.TrackAAC,
	audioSamples []*fmp4AudioSample,
) (*gomp4.Trun, int, error) {
=======
	w *mp4Writer,
	trackID int,
	audioTrack *gortsplib.TrackAAC,
	audioSamples []*fmp4AudioSample,
) (*mp4.Trun, int, error) {
>>>>>>> dahua
	/*
		traf
		- tfhd
		- tfdt
		- trun
	*/

	if len(audioSamples) == 0 {
		return nil, 0, nil
	}

<<<<<<< HEAD
	_, err := w.WriteBoxStart(&gomp4.Traf{}) // <traf>
=======
	_, err := w.writeBoxStart(&mp4.Traf{}) // <traf>
>>>>>>> dahua
	if err != nil {
		return nil, 0, err
	}

	flags := 0

<<<<<<< HEAD
	_, err = w.WriteBox(&gomp4.Tfhd{ // <tfhd/>
		FullBox: gomp4.FullBox{
=======
	_, err = w.writeBox(&mp4.Tfhd{ // <tfhd/>
		FullBox: mp4.FullBox{
>>>>>>> dahua
			Flags: [3]byte{2, byte(flags >> 8), byte(flags)},
		},
		TrackID: uint32(trackID),
	})
	if err != nil {
		return nil, 0, err
	}

<<<<<<< HEAD
	_, err = w.WriteBox(&gomp4.Tfdt{ // <tfdt/>
		FullBox: gomp4.FullBox{
=======
	_, err = w.writeBox(&mp4.Tfdt{ // <tfdt/>
		FullBox: mp4.FullBox{
>>>>>>> dahua
			Version: 1,
		},
		// sum of decode durations of all earlier samples
		BaseMediaDecodeTimeV1: uint64(durationGoToMp4(audioSamples[0].pts, time.Duration(audioTrack.ClockRate()))),
	})
	if err != nil {
		return nil, 0, err
	}

	flags = 0
	flags |= 0x01  // data offset present
	flags |= 0x100 // sample duration present
	flags |= 0x200 // sample size present

<<<<<<< HEAD
	trun := &gomp4.Trun{ // <trun/>
		FullBox: gomp4.FullBox{
=======
	trun := &mp4.Trun{ // <trun/>
		FullBox: mp4.FullBox{
>>>>>>> dahua
			Version: 0,
			Flags:   [3]byte{0, byte(flags >> 8), byte(flags)},
		},
		SampleCount: uint32(len(audioSamples)),
	}

	for _, e := range audioSamples {
<<<<<<< HEAD
		trun.Entries = append(trun.Entries, gomp4.TrunEntry{
=======
		trun.Entries = append(trun.Entries, mp4.TrunEntry{
>>>>>>> dahua
			SampleDuration: uint32(durationGoToMp4(e.duration(), time.Duration(audioTrack.ClockRate()))),
			SampleSize:     uint32(len(e.au)),
		})
	}

<<<<<<< HEAD
	trunOffset, err := w.WriteBox(trun)
=======
	trunOffset, err := w.writeBox(trun)
>>>>>>> dahua
	if err != nil {
		return nil, 0, err
	}

<<<<<<< HEAD
	err = w.WriteBoxEnd() // </traf>
=======
	err = w.writeBoxEnd() // </traf>
>>>>>>> dahua
	if err != nil {
		return nil, 0, err
	}

	return trun, trunOffset, nil
}

func mp4PartGenerate(
	videoTrack *gortsplib.TrackH264,
	audioTrack *gortsplib.TrackAAC,
	videoSamples []*fmp4VideoSample,
	audioSamples []*fmp4AudioSample,
<<<<<<< HEAD
=======
	startDTS time.Duration,
>>>>>>> dahua
) ([]byte, error) {
	/*
		moof
		- mfhd
		- traf (video)
		- traf (audio)
		mdat
	*/

<<<<<<< HEAD
	w := mp4.NewWriter()

	moofOffset, err := w.WriteBoxStart(&gomp4.Moof{}) // <moof>
=======
	w := newMP4Writer()

	moofOffset, err := w.writeBoxStart(&mp4.Moof{}) // <moof>
>>>>>>> dahua
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	_, err = w.WriteBox(&gomp4.Mfhd{ // <mfhd/>
=======
	_, err = w.writeBox(&mp4.Mfhd{ // <mfhd/>
>>>>>>> dahua
		SequenceNumber: 0,
	})
	if err != nil {
		return nil, err
	}

	trackID := 1

<<<<<<< HEAD
	var videoTrun *gomp4.Trun
=======
	var videoTrun *mp4.Trun
>>>>>>> dahua
	var videoTrunOffset int
	if videoTrack != nil {
		var err error
		videoTrun, videoTrunOffset, err = mp4PartGenerateVideoTraf(
<<<<<<< HEAD
			w, trackID, videoSamples)
=======
			w, trackID, videoSamples, startDTS)
>>>>>>> dahua
		if err != nil {
			return nil, err
		}

		trackID++
	}

<<<<<<< HEAD
	var audioTrun *gomp4.Trun
=======
	var audioTrun *mp4.Trun
>>>>>>> dahua
	var audioTrunOffset int
	if audioTrack != nil {
		var err error
		audioTrun, audioTrunOffset, err = mp4PartGenerateAudioTraf(w, trackID, audioTrack, audioSamples)
		if err != nil {
			return nil, err
		}
	}

<<<<<<< HEAD
	err = w.WriteBoxEnd() // </moof>
=======
	err = w.writeBoxEnd() // </moof>
>>>>>>> dahua
	if err != nil {
		return nil, err
	}

<<<<<<< HEAD
	mdat := &gomp4.Mdat{} // <mdat/>
=======
	mdat := &mp4.Mdat{} // <mdat/>
>>>>>>> dahua

	dataSize := 0
	videoDataSize := 0

	if videoTrack != nil {
		for _, e := range videoSamples {
			dataSize += len(e.avcc)
		}
		videoDataSize = dataSize
	}

	if audioTrack != nil {
		for _, e := range audioSamples {
			dataSize += len(e.au)
		}
	}

	mdat.Data = make([]byte, dataSize)
	pos := 0

	if videoTrack != nil {
		for _, e := range videoSamples {
			pos += copy(mdat.Data[pos:], e.avcc)
		}
	}

	if audioTrack != nil {
		for _, e := range audioSamples {
			pos += copy(mdat.Data[pos:], e.au)
		}
	}

<<<<<<< HEAD
	mdatOffset, err := w.WriteBox(mdat)
=======
	mdatOffset, err := w.writeBox(mdat)
>>>>>>> dahua
	if err != nil {
		return nil, err
	}

	if videoTrack != nil {
		videoTrun.DataOffset = int32(mdatOffset - moofOffset + 8)
<<<<<<< HEAD
		err = w.RewriteBox(videoTrunOffset, videoTrun)
=======
		err = w.rewriteBox(videoTrunOffset, videoTrun)
>>>>>>> dahua
		if err != nil {
			return nil, err
		}
	}

	if audioTrack != nil && audioTrun != nil {
		audioTrun.DataOffset = int32(videoDataSize + mdatOffset - moofOffset + 8)
<<<<<<< HEAD
		err = w.RewriteBox(audioTrunOffset, audioTrun)
=======
		err = w.rewriteBox(audioTrunOffset, audioTrun)
>>>>>>> dahua
		if err != nil {
			return nil, err
		}
	}

<<<<<<< HEAD
	return w.Bytes(), nil
=======
	return w.bytes(), nil
>>>>>>> dahua
}

func fmp4PartName(id uint64) string {
	return "part" + strconv.FormatUint(id, 10)
}

type muxerVariantFMP4Part struct {
	videoTrack *gortsplib.TrackH264
	audioTrack *gortsplib.TrackAAC
	id         uint64
<<<<<<< HEAD
=======
	startDTS   time.Duration
>>>>>>> dahua

	isIndependent    bool
	videoSamples     []*fmp4VideoSample
	audioSamples     []*fmp4AudioSample
	renderedContent  []byte
	renderedDuration time.Duration
}

func newMuxerVariantFMP4Part(
	videoTrack *gortsplib.TrackH264,
	audioTrack *gortsplib.TrackAAC,
	id uint64,
<<<<<<< HEAD
=======
	startDTS time.Duration,
>>>>>>> dahua
) *muxerVariantFMP4Part {
	p := &muxerVariantFMP4Part{
		videoTrack: videoTrack,
		audioTrack: audioTrack,
		id:         id,
<<<<<<< HEAD
=======
		startDTS:   startDTS,
>>>>>>> dahua
	}

	if videoTrack == nil {
		p.isIndependent = true
	}

	return p
}

func (p *muxerVariantFMP4Part) name() string {
	return fmp4PartName(p.id)
}

func (p *muxerVariantFMP4Part) reader() io.Reader {
	return bytes.NewReader(p.renderedContent)
}

func (p *muxerVariantFMP4Part) duration() time.Duration {
	if p.videoTrack != nil {
		ret := time.Duration(0)
		for _, e := range p.videoSamples {
			ret += e.duration()
		}
		return ret
	}

	// use the sum of the default duration of all samples,
	// not the real duration,
	// otherwise on iPhone iOS the stream freezes.
	return time.Duration(len(p.audioSamples)) * time.Second *
		time.Duration(aac.SamplesPerAccessUnit) / time.Duration(p.audioTrack.ClockRate())
}

func (p *muxerVariantFMP4Part) finalize() error {
	if len(p.videoSamples) > 0 || len(p.audioSamples) > 0 {
		var err error
		p.renderedContent, err = mp4PartGenerate(
			p.videoTrack,
			p.audioTrack,
			p.videoSamples,
<<<<<<< HEAD
			p.audioSamples)
=======
			p.audioSamples,
			p.startDTS)
>>>>>>> dahua
		if err != nil {
			return err
		}

		p.renderedDuration = p.duration()
	}

	p.videoSamples = nil
	p.audioSamples = nil

	return nil
}

func (p *muxerVariantFMP4Part) writeH264(sample *fmp4VideoSample) {
	if sample.idrPresent {
		p.isIndependent = true
	}
	p.videoSamples = append(p.videoSamples, sample)
}

func (p *muxerVariantFMP4Part) writeAAC(sample *fmp4AudioSample) {
	p.audioSamples = append(p.audioSamples, sample)
}
