// Copyright (2021) Cobalt Speech and Language Inc.

// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package juzu_test

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"testing"
	"time"

	juzu "github.com/cobaltspeech/sdk-juzu/grpc/go-juzu"
	"github.com/cobaltspeech/sdk-juzu/grpc/go-juzu/juzupb"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
)

// MockJuzuServer implements juzupb.JuzuServer so we can use it to test our
// client.  The implementation is interleaved with the appropriate test
// functions below.
type MockJuzuServer struct{}

// Test Version

var ExpectedVersionResponse = &juzupb.VersionResponse{Juzu: "none", Server: "test"}

func (s *MockJuzuServer) Version(ctx context.Context, e *empty.Empty) (*juzupb.VersionResponse, error) {
	return ExpectedVersionResponse, nil
}

func TestVersion(t *testing.T) {
	svr, port, err := setupGRPCServer()
	defer svr.Stop()

	if err != nil {
		t.Fatalf("could not set up testing server: %v", err)
	}

	c, err := juzu.NewClient(fmt.Sprintf("localhost:%d", port), juzu.WithInsecure())
	if err != nil {
		t.Errorf("could not create client: %v", err)
	}
	defer c.Close()

	v, err := c.Version()
	if err != nil {
		t.Errorf("did not expect error in version; got %v", err)
	}

	if !proto.Equal(v, ExpectedVersionResponse) {
		t.Errorf("version failed; got %v, want %v", v, ExpectedVersionResponse)
	}
}

// Test ListModels

var ExpectedListModelsResponse = &juzupb.ListModelsResponse{
	Models: []*juzupb.Model{
		&juzupb.Model{Id: "1"},
	},
}

func (s *MockJuzuServer) ListModels(ctx context.Context, e *empty.Empty) (*juzupb.ListModelsResponse, error) {
	return ExpectedListModelsResponse, nil
}

func TestListModels(t *testing.T) {
	svr, port, err := setupGRPCServer()
	defer svr.Stop()

	if err != nil {
		t.Fatalf("could not set up testing server: %v", err)
	}

	c, err := juzu.NewClient(fmt.Sprintf("localhost:%d", port), juzu.WithInsecure())
	if err != nil {
		t.Errorf("could not create client: %v", err)
	}
	defer c.Close()

	m, err := c.ListModels()
	if err != nil {
		t.Errorf("did not expect error in listmodels; got %v", err)
	}

	if !proto.Equal(m, ExpectedListModelsResponse) {
		t.Errorf("listmodels failed; got %v, want %v", m, ExpectedListModelsResponse)
	}
}

// Test Streaming Diarize

var ExpectedDiarizationResult = &juzupb.DiarizationResult{
	SpeakerLabels: []string{"0", "1"},
	IsPartial:     false,
	Segments: []*juzupb.Segment{
		&juzupb.Segment{
			SpeakerLabel: "0",
			Transcript:   "Hello",
		},
		&juzupb.Segment{
			SpeakerLabel: "1",
			Transcript:   "Goodbye",
		},
	},
}

var ExpectedStreamingDiarizeResponse = &juzupb.DiarizationResponse{
	Results: []*juzupb.DiarizationResult{ExpectedDiarizationResult},
}

func (s *MockJuzuServer) StreamingDiarize(stream juzupb.Juzu_StreamingDiarizeServer) error {
	// verify that first message is config
	msg, err := stream.Recv()
	if err != nil {
		return fmt.Errorf("streaming diarization failed: missing first message")
	}

	if msg.GetConfig() == nil {
		return fmt.Errorf("streaming diarization failed: first message should be a config message")
	}

	if msg.GetConfig().GetModelId() == "test-nil" {
		// send a response with a nil result.  grpc does not allow this, and
		// this should fail.  but we test this to ensure that our client does
		// not need to check for nil before using the object.
		err = stream.Send(&juzupb.DiarizationResponse{Results: []*juzupb.DiarizationResult{nil}})
		return err
	}

	// verify that remaining messages are audio messages, and there are at least three of those.
	count := 0
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("streaming diarization failed: %v", err)
		}

		if req.GetAudio() == nil {
			return fmt.Errorf("streaming diarization failed: all messages after the first should be audio messages")
		}

		count++
	}

	if count < 3 {
		return fmt.Errorf("streaming diarization failed: expecting at least 3 test audio messages, got %d", count)
	}

	if err := stream.Send(ExpectedStreamingDiarizeResponse); err != nil {
		return fmt.Errorf("streaming diarization failed: %v", err)
	}

	return nil
}

func TestStreamingDiarize(t *testing.T) {
	svr, port, err := setupGRPCServer()
	defer svr.Stop()

	if err != nil {
		t.Fatalf("could not set up testing server: %v", err)
	}

	c, err := juzu.NewClient(fmt.Sprintf("localhost:%d", port), juzu.WithInsecure())
	if err != nil {
		t.Errorf("could not create client: %v", err)
	}
	defer c.Close()

	var got *juzupb.DiarizationResponse

	handleResult := func(resp *juzupb.DiarizationResponse) {
		got = resp
	}

	audio := make([]byte, 10*4096) // all zeros

	err = c.StreamingDiarize(context.Background(), &juzupb.DiarizationConfig{}, bytes.NewReader(audio), handleResult)
	if err != nil {
		t.Errorf("did not expect error in streaming diarization; got %v", err)
	}

	if !proto.Equal(got, ExpectedStreamingDiarizeResponse) {
		t.Errorf("streaming diarization failed: got %v; want %v", got, ExpectedStreamingDiarizeResponse)
	}

	// Now check that even if the server sends a nil message, it won't be a
	// nil pointer when we receive the output.
	err = c.StreamingDiarize(context.Background(), &juzupb.DiarizationConfig{ModelId: "test-nil"}, bytes.NewReader(audio), handleResult)
	if err != nil {
		t.Errorf("streaming diarization should not have failed when server sent a nil message, but got error %v", err)
	}

	if got.Results[0] == nil {
		t.Errorf("streaming diarization should not have provided a nil pointer in the result")
	}

}

// Test Streaming Buffer Size Option
func TestStreamingBufSize(t *testing.T) {
	svr, port, err := setupGRPCServer()
	defer svr.Stop()

	if err != nil {
		t.Fatalf("could not set up testing server: %v", err)
	}

	_, err = juzu.NewClient(fmt.Sprintf("localhost:%d", port), juzu.WithInsecure(), juzu.WithStreamingBufferSize(0))
	if err == nil {
		t.Errorf("client creation with streaming buffer size 0, want failure, got success")
	}

	c, err := juzu.NewClient(fmt.Sprintf("localhost:%d", port), juzu.WithInsecure(), juzu.WithStreamingBufferSize(1))
	if err != nil {
		t.Errorf("client creation with streaming buffer size 1, want success, got %v", err)
	}
	defer c.Close()

}

func TestClient_InvalidURL(t *testing.T) {
	if _, err := juzu.NewClient(fmt.Sprintf("wrong_localhost:2727"), juzu.WithInsecure(),
		juzu.WithConnectTimeout(200*time.Millisecond)); err == nil {
		t.Errorf("connecting to invalid server: want error, got nil")
	}
}

func TestClient_Insecure(t *testing.T) {
	svr, port, err := setupGRPCServer() // server without TLS
	defer svr.Stop()

	if err != nil {
		t.Fatalf("could not set up testing server: %v", err)
	}

	// connection with TLS
	if _, err := juzu.NewClient(fmt.Sprintf("localhost:%d", port),
		juzu.WithConnectTimeout(200*time.Millisecond)); err == nil {
		t.Errorf("tls connection to non-tls server: want error, got nil")
	}
}

func setupGRPCServer() (*grpc.Server, int, error) {
	lis, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, 0, err
	}

	s := grpc.NewServer()
	juzupb.RegisterJuzuServer(s, &MockJuzuServer{})
	go func() { _ = s.Serve(lis) }()
	return s, lis.Addr().(*net.TCPAddr).Port, nil
}
