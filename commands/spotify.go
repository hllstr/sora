package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"sora/lib"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"
)

type SpotifyData struct {
	Songs []struct {
		Title     string `json:"title"`
		Artist    string `json:"artist"`
		Thumbnail string `json:"thumbnail"`
		URL       string `json:"url"`
		Duration  string `json:"duration"`
	} `json:"songs"`
	ContentType string `json:"contentType"`
}

func init() {
	Plugin(Cmd{
		Name:     "spotify",
		Category: "tools",
		Alias:    []string{"sp"},
		Desc:     "Download song",
		Exec:     spotify,
	})
}

func spotify(ctx *CommandContext) {
	API := "https://spotdown.org/api"
	searchAPI := API + "/song-details?url="
	donlodAPI := API + "/direct-download?url="
	search := searchAPI + url.QueryEscape(ctx.RawArgs)
	req, err := http.NewRequest("GET", search, nil)
	if err != nil {
		ctx.Client.Log.Errorf("Error Creating Request: %v", err)
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var result SpotifyData
	err = json.Unmarshal(body, &result)
	if err != nil {
		ctx.Client.Log.Errorf("Error unmarshaling data: %v", err)
		return
	}
	topResult := result.Songs[0]
	resp, err = http.Get(donlodAPI + topResult.URL)
	if err != nil {
		ctx.Client.Log.Errorf("Error downloading file: %v", err)
		return
	}
	defer resp.Body.Close()
	data, _ := io.ReadAll(resp.Body)
	song, err := ctx.Client.Upload(ctx.Ctx, data, whatsmeow.MediaAudio)
	if err != nil {
		ctx.Client.Log.Errorf("Error uploading file: %v", err)
		return
	}
	audioMsg := &waE2E.AudioMessage{
		URL:           proto.String(song.URL),
		Mimetype:      proto.String("audio/mpeg"),
		FileSHA256:    song.FileSHA256,
		FileLength:    proto.Uint64(song.FileLength),
		MediaKey:      song.MediaKey,
		FileEncSHA256: song.FileEncSHA256,
		DirectPath:    proto.String(song.DirectPath),
		ContextInfo: &waE2E.ContextInfo{
			StanzaID:    &ctx.Message.Info.ID,
			Participant: proto.String(ctx.Message.Info.Sender.String()),
			ExternalAdReply: &waE2E.ContextInfo_ExternalAdReplyInfo{
				Title:                 proto.String(topResult.Title),
				Body:                  proto.String(topResult.Artist),
				SourceURL:             proto.String(topResult.URL),
				ThumbnailURL:          proto.String(topResult.Thumbnail),
				RenderLargerThumbnail: proto.Bool(true),
				MediaType:             waE2E.ContextInfo_ExternalAdReplyInfo_IMAGE.Enum(),
			},
		},
	}
	if exp, ok := lib.GetEphemeralDuration(ctx.Message); ok {
		audioMsg.ContextInfo.Expiration = proto.Uint32(exp)
	}
	_, err = ctx.Client.SendMessage(ctx.Ctx, ctx.Message.Info.Chat, &waE2E.Message{AudioMessage: audioMsg}, lib.Bypass(ctx.Client, ctx.Message.Info.Chat))
	if err != nil {
		ctx.Client.Log.Errorf("Error sending message: %v", err)
	}
}
