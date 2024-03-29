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

// Package juzu provides for interacting with an instance of juzu server using
// GRPC for performing speech recognition.
package juzu

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/cobaltspeech/sdk-juzu/grpc/go-juzu/juzupb"
	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const defaultStreamingBufSize uint32 = 8192
const defaultConnectTimeout = 2 * time.Second

// Client is an object for interacting with the juzu GRPC API.
//
// All methods except Close may be called concurrently.
type Client struct {
	conn             *grpc.ClientConn
	juzu             juzupb.JuzuClient
	insecure         bool
	tlsCfg           tls.Config
	streamingBufSize uint32
	connectTimeout   time.Duration
}

// NewClient creates a new Client that connects to a juzu Server listening on
// the provided address.  Transport security is enabled by default.  Use Options
// to override default settings if necessary.
func NewClient(addr string, opts ...Option) (*Client, error) {
	c := Client{}
	c.streamingBufSize = defaultStreamingBufSize
	c.connectTimeout = defaultConnectTimeout

	for _, opt := range opts {
		err := opt(&c)
		if err != nil {
			return nil, fmt.Errorf("unable to create a client: %v", err)
		}
	}

	var dopt grpc.DialOption
	if c.insecure {
		dopt = grpc.WithInsecure()
	} else {
		dopt = grpc.WithTransportCredentials(credentials.NewTLS(&c.tlsCfg))
	}

	ctx, cancel := context.WithTimeout(context.Background(), c.connectTimeout)
	defer cancel()

	conn, err := grpc.DialContext(ctx, addr, dopt, grpc.WithBlock())
	if err != nil {
		return nil, fmt.Errorf("unable to create a client: %v", err)
	}
	c.conn = conn
	c.juzu = juzupb.NewJuzuClient(c.conn)
	return &c, nil
}

// Option configures how we setup the connection with a server.
type Option func(*Client) error

// WithInsecure returns an Option which disables transport security for this
// Client.  Use this when connecting to a non-TLS enabled juzu server, such as
// during debugging.
func WithInsecure() Option {
	return func(c *Client) error {
		c.insecure = true
		return nil
	}
}

// WithServerCert returns an Option which sets up the given PEM certificate as a
// root certificate that can validate the certificate presented by the server we
// are connecting to.  Use this when connecting to an instance of juzu server
// that is using a self-signed certificate.
func WithServerCert(cert []byte) Option {
	return func(c *Client) error {
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(cert); !ok {
			return fmt.Errorf("unable to use given caCert")
		}
		c.tlsCfg.RootCAs = caCertPool
		return nil
	}
}

// WithClientCert returns an Option which sets up the given PEM certificate and
// key as the credentials presented by this Client when connecting to a server.
// Use this when setting up mutually authenticated TLS.
func WithClientCert(certPem []byte, keyPem []byte) Option {
	return func(c *Client) error {
		clientCert, err := tls.X509KeyPair(certPem, keyPem)
		if err != nil {
			return err
		}

		c.tlsCfg.Certificates = []tls.Certificate{clientCert}
		return nil
	}
}

// WithStreamingBufferSize returns an Option that sets up the buffer size
// (bytes) of each message sent from the Client to the server during streaming
// GRPC calls.  Use this only if Cobalt recommends you to do so.  A value n>0 is
// required.
func WithStreamingBufferSize(n uint32) Option {
	return func(c *Client) error {
		if n == 0 {
			return fmt.Errorf("invalid streaming buffer size of 0")
		}
		c.streamingBufSize = n
		return nil
	}
}

// WithConnectTimeout returns an Option that configures the timeout for
// establishing grpc connection with the server.  Use this only when you are on
// a slow network and when Cobalt recommends you to do so.
func WithConnectTimeout(t time.Duration) Option {
	return func(c *Client) error {
		c.connectTimeout = t
		return nil
	}
}

// Close closes the connection to the API service.  The user should only invoke
// this when the client is no longer needed.  Pending or in-progress calls to
// other methods may fail with an error if Close is called, and any subsequent
// calls with this client will also fail.
func (c *Client) Close() error {
	return c.conn.Close()
}

// Version queries the server for its version
func (c *Client) Version() (*juzupb.VersionResponse, error) {
	return c.juzu.Version(context.Background(), &empty.Empty{})
}

// ListModels queries the server for the list of available diarization models
func (c *Client) ListModels() (*juzupb.ListModelsResponse, error) {
	return c.juzu.ListModels(context.Background(), &empty.Empty{})
}

// DiarizationResponseHandler is a type of callback function that will be called
// when the `StreamingDiarize` method is running.  For each response received
// from juzu server, this method will be called once.  The provided
// DiarizationResponse is guaranteed to be non-nil.  Since this function is
// executed as part of the streaming process, it should preferably return
// quickly and certainly not block.
type DiarizationResponseHandler func(*juzupb.DiarizationResponse)

// StreamingDiarize uses the bidirectional streaming API for performing speech
// diarization. It sets up diarization using the given cfg.
//
// Data is read from the given audio reader into a buffer and streamed to Juzu
// server.  The default buffer size may be overridden using Options when
// creating the Client.
//
// After results are received from the Juzu server, they will be sent to the
// provided handlerFunc.
//
// If any error occurs while reading the audio or sending it to the server, this
// method will immediately exit, returning that error.
//
// This function returns only after all results have been passed to the
// resultHandler.
func (c *Client) StreamingDiarize(
	ctx context.Context,
	cfg *juzupb.DiarizationConfig,
	audio io.Reader,
	handlerFunc DiarizationResponseHandler,
) error {

	stream, err := c.juzu.StreamingDiarize(ctx)
	if err != nil {
		return fmt.Errorf("unable to start streaming diarization: %v", err)
	}

	// There are two concurrent processes going on.  We will create a new
	// goroutine to read audio and stream it to the server.  This goroutine
	// will receive results from the stream.  Errors could occur in both
	// goroutines.  We therefore setup a channel, errCh, to hold these
	// errors. Both goroutines are designed to send up to one error, and
	// return immediately. Therefore we use a bufferred channel with a
	// capacity of two.
	errCh := make(chan error, 2)

	// start streaming audio in a separate goroutine
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		if err := sendaudio(
			stream, cfg, audio, c.streamingBufSize,
		); err != nil && err != io.EOF {
			// if sendaudio encountered io.EOF, it's only a
			// notification that the stream has closed.  The actual
			// status will be obtained in a subsequent Recv call, in
			// the other goroutine below.  We therefore only forward
			// non-EOF errors.
			errCh <- err
		}
		wg.Done()
	}()

	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			errCh <- err
			break
		}

		handlerFunc(in)
	}

	wg.Wait()

	select {
	case err := <-errCh:
		// There may be more than one error in the channel, but it is
		// very likely they are related (e.g. connection reset causing
		// both the send and recv to fail) and we therefore return the
		// first error and discard the other.
		return fmt.Errorf("streaming recognition failed: %v", err)
	default:
		return nil
	}
}

// sendaudio sends audio to a stream.
func sendaudio(stream juzupb.Juzu_StreamingDiarizeClient,
	cfg *juzupb.DiarizationConfig, audio io.Reader, bufSize uint32) error {

	// The first message needs to be a config message, and all subsequent
	// messages must be audio messages.

	// Send the diarization config
	if err := stream.Send(&juzupb.StreamingDiarizeRequest{
		Request: &juzupb.StreamingDiarizeRequest_Config{Config: cfg},
	}); err != nil {
		// if this failed, we don't need to CloseSend
		return err
	}

	// Stream the audio.
	buf := make([]byte, bufSize)
	for {
		n, err := audio.Read(buf)
		if n > 0 {
			if err2 := stream.Send(&juzupb.StreamingDiarizeRequest{
				Request: &juzupb.StreamingDiarizeRequest_Audio{
					Audio: &juzupb.DiarizationAudio{Data: buf[:n]},
				},
			}); err2 != nil {
				// if we couldn't Send, the stream has
				// encountered an error and we don't need to
				// CloseSend.
				return err2
			}
		}

		if err != nil {
			// err could be io.EOF, or some other error reading from
			// audio.  In any case, we need to CloseSend, send the
			// appropriate error to errCh and return from the function
			if err2 := stream.CloseSend(); err2 != nil {
				return err2
			}
			if err != io.EOF {
				return err
			}
			return nil
		}
	}
}
