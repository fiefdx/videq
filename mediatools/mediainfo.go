package mediatools

import (
	"errors"
	"fmt"
	alog "github.com/cenkalti/log"
	"github.com/codeskyblue/go-sh"
	"os"
	//	"strconv"
	"strings"
	"time"
)

type MediaInfo struct {
	log alog.Logger
}

func NewMediaInfo(log alog.Logger) *MediaInfo {
	m := new(MediaInfo)
	m.log = log
	return m
}

type MediaFileInfo struct {
	FileName       string
	FileSize_bytes string
	VideoCount     string
	AudioCount     string
	Duration_ms    string
	Duration       time.Duration
	Format         string
	CodecID        string
	Resolution     string
	Width          string
	Height         string
	Standard       string
	Codec          string
	Bitrate_bps    string
	Framerate      string
	AspectRatio    string
	Audio          string
}

/*
FileName: r2w_1080p.mov
FileSize_bytes: 104937987
VideoCount: 1
AudioCount: 1
Duration_ms: 90125
Format: MPEG-4
CodecID: qt
Resolution: 1920x816
Width: 1920
Height: 816
Standard:
Codec: AVC Main@L4.0
Bitrate_bps: 9185470
Framerate: 24.000 fps
AspectRatio: 2.35:1
Audio: English 128 Kbps CBR 2 chnls AAC LC
*/

/*
fetches media info
usage:
mt := mediatools.NewMediaInfo(log)
minfo, err := mt.GetMediaInfo("_test/master_1080.mp4")
if err != nil {
	log.Fatal(err)
}
log.Infof("%#v", minfo)
*/
func (m *MediaInfo) GetMediaInfo(fileName string) (fileInfo MediaFileInfo, err error) {

	// timeout should be a session
	//	out, err := sh.Command("ping", "-t", "127.0.0.1").SetTimeout(time.Second * 60).Output()

	mediaInfoParams := `--Inform=General;FileName:: %FileName%.%FileExtension%\r\nFileSize_bytes:: %FileSize%\r\nVideoCount:: %VideoCount%\r\nAudioCount:: %AudioCount%\r\nDuration_ms:: %Duration%\r\nFormat:: %Format%\r\nCodecID:: %CodecID%\r\n
Video;Resolution:: %Width%x%Height%\r\nWidth:: %Width%\r\nHeight:: %Height%\r\nStandard:: %Standard%\r\nCodec:: %Codec/String% %Format_Profile%\r\nBitrate_bps:: %BitRate%\r\nFramerate:: %FrameRate% fps\r\nAspectRatio:: %DisplayAspectRatio/String%\r\n
Audio;Audio:: %Language/String% %BitRate/String% %BitRate_Mode% %Channel(s)% chnls %Codec/String%\r\n
Text;%Language/String%
Text_Begin;Subs:
Text_Middle;, 
Text_End;.\r\n
`
	out, err := sh.Command("mediainfo", mediaInfoParams, fileName).SetTimeout(time.Second * 60).Output()
	// fmt.Printf("output:(%s), err(%v)\n", string(out), err)
	if err == sh.ErrExecTimeout {
		m.log.Errorf("shell exec timeouteded.", err)
	}
	if err != nil {
		m.log.Errorf("sh.Command error:", err)
		return fileInfo, err
	}

	//output = strings.Replace(string(out), "\n", "<br>", -1)

	lines := strings.Split(string(out), "\n")
	//fileInfo := new(MediaFileInfo)

	for _, line := range lines {
		paramArr := strings.Split(line, "::")
		if len(paramArr) != 2 {
			continue
		}

		paramName := strings.Trim(paramArr[0], " ")
		paramValue := strings.Trim(paramArr[1], " ")

		switch paramName {
		case `FileName`:
			fileInfo.FileName = paramValue
		case `FileSize_bytes`:
			fileInfo.FileSize_bytes = paramValue
		case `VideoCount`:
			fileInfo.VideoCount = paramValue
		case `AudioCount`:
			fileInfo.AudioCount = paramValue
		case `Duration_ms`:
			fileInfo.Duration_ms = paramValue
			// durationDurationInt64, err := strconv.ParseInt(fileInfo.Duration_ms, 10, 64)
			// fileInfo.Duration = time.Millisecond * durationDurationInt64 / 100
			dur, err := time.ParseDuration(paramValue + "ms")
			if err == nil {
				fileInfo.Duration = dur
				m.log.Debug(dur)
			}
		case `Format`:
			fileInfo.Format = paramValue
		case `CodecID`:
			fileInfo.CodecID = paramValue
		case `Resolution`:
			fileInfo.Resolution = paramValue
		case `Width`:
			fileInfo.Width = paramValue
		case `Height`:
			fileInfo.Height = paramValue
		case `Standard`:
			fileInfo.Standard = paramValue
		case `Codec`:
			fileInfo.Codec = paramValue
		case `Bitrate_bps`:
			fileInfo.Bitrate_bps = paramValue
		case `Framerate`:
			fileInfo.Framerate = paramValue
		case `AspectRatio`:
			fileInfo.AspectRatio = paramValue
		case `Audio`:
			fileInfo.Audio = paramValue
		}
	}

	return fileInfo, nil
}

func (m *MediaInfo) CheckMedia(fileName string) (ok bool, fileInfo MediaFileInfo, err error) {

	exists, err := m.checkIfFileExists(fileName)
	if err != nil {
		return false, fileInfo, err
	}
	if exists == false {
		return false, fileInfo, errors.New(fmt.Sprintf("File '%s' does not exists.", fileName))
	}

	fileInfo, err = m.GetMediaInfo(fileName)
	if err != nil {
		m.log.Error(err)
		return false, fileInfo, err
	}
	//	m.log.Infof("%#v", fileInfo)

	return true, fileInfo, nil
}

func (m *MediaInfo) checkIfFileExists(fileName string) (bool, error) {
	_, err := os.Stat(fileName)
	if os.IsNotExist(err) {
		return false, nil
	}
	if err != nil {
		return false, err
	}
	return true, nil
}
