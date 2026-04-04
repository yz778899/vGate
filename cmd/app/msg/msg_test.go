package msg

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/yz778899/vGate/net/data"
)

func TestDecoder_LoginRequest(t *testing.T) {
	// WsMsg declares its own Content field, which shadows BaseMsg.Content for ws.Content.
	ws := &data.WebsocketMsg{
		BaseMsg: data.BaseMsg{
			Cmd:   data.Request,
			Topic: "login",
		},
		Content: json.RawMessage(`{"User":"alice","Pass":"secret"}`),
	}
	var req LoginRequest
	if err := Decoder(ws, &req); err != nil {
		t.Fatal(err)
	}
	if req.User != "alice" || req.Pass != "secret" {
		t.Fatalf("%+v", req)
	}
}

func TestDecoder_LoginRequest_InvalidJSON(t *testing.T) {
	ws := &data.WebsocketMsg{
		BaseMsg: data.BaseMsg{
			Content: json.RawMessage(`{`),
		},
	}
	var req LoginRequest
	if err := Decoder(ws, &req); err == nil {
		t.Fatal("expected error")
	}
}

func TestGameListResponse_JSONRoundTrip(t *testing.T) {
	ts := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)
	resp := GameListResponse{
		Games: []Game{
			{
				Id:         1,
				Name:       "G1",
				Desc:       "d",
				Icon:       "i",
				Url:        "u",
				Status:     1,
				CreateTime: ts,
				UpdateTime: ts,
			},
		},
	}
	b, err := json.Marshal(resp)
	if err != nil {
		t.Fatal(err)
	}
	var out GameListResponse
	if err := json.Unmarshal(b, &out); err != nil {
		t.Fatal(err)
	}
	if len(out.Games) != 1 || out.Games[0].Name != "G1" || !out.Games[0].CreateTime.Equal(ts) {
		t.Fatalf("%+v", out)
	}
}
