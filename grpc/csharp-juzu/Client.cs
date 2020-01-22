// Copyright (2019) Cobalt Speech and Language Inc.

using System.Collections.Generic;
using System.IO;
using System.Threading.Tasks;
using Grpc.Core;

namespace JuzusvrClient {

    public delegate void ResponseHandler(CobaltSpeech.Juzu.DiarizationResponse resp);

    public class Client {

        private string serverURL;
        private Grpc.Core.ChannelCredentials creds;
        private Grpc.Core.Channel channel;
        private CobaltSpeech.Juzu.Juzu.JuzuClient client;

        // Returns the default channel options specified when creating a
        // channel. Corresponds to grpc_channel_args from grpc/grpc.h in the
        // grpc repo. Commonly Any of the GRPC_ARG_* channel options names
        // defined in grpc_types.h can be used.
        private List<ChannelOption> defaultChannelOptions() {
            return new List<ChannelOption> {
                // After a duration of this time the client pings its peer to see if the
                // transport is still alive. Int valued, milliseconds. Set to 15 seconds.
                new ChannelOption("grpc.keepalive_time_ms", 15 * 1000),

                // After waiting for a duration of this time, if the keepalive ping sender does
                // not receive the ping ack, it will close the transport. Int valued,
                // milliseconds. Set to 20 seconds.
                new ChannelOption("grpc.keepalive_timeout_ms", 20 * 1000),

                // Is it permissible to send keepalive pings without any outstanding streams.
                // Int valued, 0(false) / 1(true).                
                new ChannelOption("grpc.keepalive_permit_without_calls", 1),

                // How many pings can we send before needing to send a data frame or header
                // frame? (0 indicates that an infinite number of pings can be sent without
                // sending a data frame or header frame)
                new ChannelOption("grpc.http2.max_pings_without_data", 0),

                // Minimum time between sending successive ping frames without receiving any
                // data frame, Int valued, milliseconds. Set to 10 seconds.
                new ChannelOption("grpc.http2.min_time_between_pings_ms", 10 * 1000),
            };
        }

        // Creates a client to Juzusvr. If insecure is set 
        // to True, TLS will be disabled.
        public Client(string url, bool insecure) {

            this.serverURL = url;

            if (insecure) {
                // no TLS
                this.creds = Grpc.Core.ChannelCredentials.Insecure;
            } else {
                // SSL credentials loaded from disk file pointed to by the 
                // GRPC_DEFAULT_SSL_ROOTS_FILE_PATH environment variable.
                // If that fails, gets the roots certificates from a well 
                // known place on disk.
                this.creds = new Grpc.Core.SslCredentials();
            }
            this.channel = new Grpc.Core.Channel(url, this.creds, this.defaultChannelOptions());
            this.client = client = new CobaltSpeech.Juzu.Juzu.JuzuClient(channel);
        }

        // Creates a client to Juzusvr with SSL credentials
        // from a string containing PEM encoded root certificates,
        // that can validate the certificate presented by the server.
        public Client(string url, string rootCert) {
            this.serverURL = url;
            this.creds = new Grpc.Core.SslCredentials(rootCert);
            this.channel = new Grpc.Core.Channel(url, this.creds, this.defaultChannelOptions());
            this.client = client = new CobaltSpeech.Juzu.Juzu.JuzuClient(channel);
        }

        // Creates a client to Juzusvr with mutually authenticated TLS.
        // The PEM encoded root certificates, PEM encoded client certificate
        // and the client's PEM private key must be provided as strings.
        public Client(string url, string rootCert, string clientCert, string clientKey) {
            this.serverURL = url;
            var keyCertPair = new Grpc.Core.KeyCertificatePair(clientCert, clientKey);
            this.creds = new Grpc.Core.SslCredentials(rootCert, keyCertPair);
            this.channel = new Grpc.Core.Channel(url, this.creds, this.defaultChannelOptions());
            this.client = client = new CobaltSpeech.Juzu.Juzu.JuzuClient(channel);
        }

        // Queries version of the server
        public CobaltSpeech.Juzu.VersionResponse Version() {
            return this.client.Version(
                new Google.Protobuf.WellKnownTypes.Empty());
        }
        // Gets list of diarization models on the server
        public CobaltSpeech.Juzu.ListModelsResponse ListModels() {
            return this.client.ListModels(
                new Google.Protobuf.WellKnownTypes.Empty());
        }

        // Sets up the bi-directional gRPC stream to Juzusvr for diarization
        // (and transcription if cubic model specified in the diarization
        // config). The stream can be any readable stream (i.e needs to only
        // implement the Stream.Read method). The results are returned by the
        // server and passed into the given ResponseHandler function after all
        // the data in the stream has been read.
        public async Task StreamingDiarizeAsync(
            System.IO.Stream stream, DiarizationConfig diarCfg, ResponseHandler handleFunc) {

            // Creating config for the diarization request
            // Mapping audio encoding
            CobaltSpeech.Juzu.DiarizationConfig.Types.Encoding encoding;
            switch (diarCfg.Encoding) {
                case (AudioEncoding.WAV):
                    encoding = CobaltSpeech.Juzu.DiarizationConfig.Types.Encoding.Wav;
                    break;
                case (AudioEncoding.FLAC):
                    encoding = CobaltSpeech.Juzu.DiarizationConfig.Types.Encoding.Flac;
                    break;
                case (AudioEncoding.RAW):
                    encoding = CobaltSpeech.Juzu.DiarizationConfig.Types.Encoding.RawLinear16;
                    break;
                default:
                    throw new InvalidDataException("unknown audio encoding");
            }
            // Add the config to the request
            var request = new CobaltSpeech.Juzu.StreamingDiarizeRequest();
            request.Config = new CobaltSpeech.Juzu.DiarizationConfig {
                ModelId = diarCfg.JuzuModelID,
                CubicModelId = diarCfg.CubicModelID,
                NumSpeakers = diarCfg.NumSpeakers,
                SampleRate = diarCfg.SampleRate,
                AudioEncoding = encoding,
            };

            // Setting up bidirectional stream
            var call = this.client.StreamingDiarize();

            // Setting up receive task
            using(call) {
                var responseReaderTask = Task.Run(async() => {
                    // Wait for response
                    while (await call.ResponseStream.MoveNext()) {
                        var response = call.ResponseStream.Current;
                        // Do stuff with the response
                        handleFunc(response);
                    }
                });

                // Send config first, followed by the audio
                {
                    await call.RequestStream.WriteAsync(request);

                    // Setup object for streaming audio
                    request.Config = null;
                    request.Audio = new CobaltSpeech.Juzu.DiarizationAudio { };

                    // Send the audio, in 8kb chunks; audio can be streamed 
                    // from a mic instead of a file in the same way, in which
                    // case the encoding and sample rate configs needs to be
                    // set appropriately.
                    const int chunkSize = 8192;

                    int bytesRead;
                    var buffer = new byte[chunkSize];
                    while ((bytesRead = stream.Read(buffer, 0, buffer.Length)) > 0) {
                        var bytes = Google.Protobuf.ByteString.CopyFrom(buffer);
                        request.Audio.Data = bytes;
                        await call.RequestStream.WriteAsync(request);
                    }

                    // Close the sending stream
                    await call.RequestStream.CompleteAsync();
                }

                // Wait for the response to come back through the receiving stream
                await responseReaderTask;
            }
        }
    }

    public struct DiarizationConfig {
        public string JuzuModelID;
        public string CubicModelID;
        public uint NumSpeakers;
        public uint SampleRate;
        public AudioEncoding Encoding;
    }

    public enum AudioEncoding {
        WAV,
        FLAC,
        RAW,
    }
}
