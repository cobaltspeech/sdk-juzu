---
title: "Streaming Diarization"
description: "Describes how to stream audio for diarization."
weight: 24
---

The following example shows how to diarize and transcribe an audio file using
Juzu’s Streaming Diarize Request. The stream can come from a file on disk or be
directly from a microphone in real time. The diarization (and transcription)
results are returned after the stream is ended and all the audio has been sent
to the server.

For real-time streaming transcription without diarization, call cubicsvr directly.
(See [StreamingRecognize](https://cobaltspeech.github.io/sdk-cubic/using-cubic-sdk/streaming/) in the Cubic SDK documentation.)

<!--more-->

### Streaming from an audio file

We support several file formats including RAW, WAV and FLAC. For more details, please
see the protocol buffer specification file in the SDK repository ([grpc/juzu.proto](https://github.com/cobaltspeech/sdk-juzu/blob/master/grpc/juzu.proto)).
The examples below use a WAV file as input to the streaming diarization (and
transcription).

{{< tabs >}}

{{< tab "Python" "py" >}}

import juzu

# This is Cobalt's demo server. Change it to your own local Juzusvr if you have one running.
serverAddress = 'demo.cobaltspeech.com:2727'

# set insecure=True for connecting to server not using TLS
client = juzu.Client(serverAddress, insecure=False)

# get list of available models
modelResp = client.ListModels()
for model in modelResp.models:
    print("ID = {}\t Name = {}\t [SampleRate = {} Hz]".format(model.id, model.name, model.attributes.sample_rate))

# use the first available model
juzuModel = modelResp.models[0]
juzuModelID = juzuModel.id

# Using the first cubic model that is compatible with the chosen
# Juzu model to transcribe.
#
# Note that Cubicsvr must also be running and the address:port
# found in cubicsvr.cfg.toml or be obtained via sdk-cubic.
cubicModelID = juzuModel.attributes.compatible_cubic_models[0]

cfg = juzu.DiarizationConfig(
    model_id = juzuModelID,
    cubic_model_id = cubicModelID,
    num_speakers = 2,               # number of speakers expected in the audio file
    audio_encoding = "WAV",         # supported : "RAW_LINEAR16", "FLAC", "WAV", "MP3"
    sample_rate = 16000,            # must match juzu model's expected sample rate
)

# client.StreamingDiarize takes any binary
# stream object that has a read(nBytes) method.
# The method should return nBytes from the stream.

# open audio file stream
audio = open('test.wav', 'rb')

# helper function convert protobuf duration objects
# (which stores the time split into in integer seconds
# and integer nano seconds) into single floating value
# in seconds
def protoDurToSec(dur):
    return float(dur.seconds) + float(dur.nanos) * 1e-9

# defining function to print speaker segments and transcripts to screen
def handleResults(diarizationResp):
    for result in diarizationResp.results:
        for segment in result.segments:
            print("{start:.3f} - {end:.3f}\t{speaker}:\t{transcript}\n".format(
                start = protoDurToSec(segment.start_time),
                end = protoDurToSec(segment.end_time),
                speaker = segment.speaker_label,
                transcript = segment.transcript,
                ))

# sending streaming request to Juzu and
# waiting for results to return
for resp in client.StreamingDiarize(cfg, audio):
    handleResults(resp)

{{< /tab >}}

{{< tab "C#" "csharp">}}

using System;
using System.IO;
using System.Net;
using System.Text;
using System.Threading.Tasks;

namespace JuzusvrClient {
    class Program {

        static async Task Main (string[] args) {
            
            // This is Cobalt's demo server. Change it to your own local Juzusvr if you have one running.
            var url = "demo.cobaltspeech.com:2727";
            string audioFile = "test.wav";

            // Set insecure = true for connecting to server not using TLS.
            var insecure = false;
            var client = new Client (url, insecure);

            // Getting list of diarization models on the server
            var modelResp = client.ListModels ();
            Console.WriteLine ("\nAvailable models:\n");
            foreach (var model in modelResp.Models) {
                Console.WriteLine ("{0}\t{1}\t{2}\n", model.Id, model.Name, model.Attributes.SampleRate);
            }

            // Creating config for Diarizing + Transcribing file with the first
            // Juzu Model available and the first compatible Cubic model.
            var diarCfg = new DiarizationConfig {
                JuzuModelID = modelList.Models[0].Id,
                CubicModelID = modelList.Models[0].CompatibleCubicModels[0],
                NumSpeakers = 2,        // use 0 if unknown
                SampleRate = 16000,
                Encoding = AudioEncoding.WAV,
            };

            // Define callback function to print results on screen; could be
            // modified to do other things with the results as well.
            ResponseHandler handleFunc = delegate (CobaltSpeech.Juzu.DiarizationResponse resp) {
                foreach (var result in resp.Results) {
                    foreach (var seg in result.Segments) {
                        Console.WriteLine ("{0} : {1}\t{2}\t{3}",
                            seg.StartTime, seg.EndTime, seg.SpeakerLabel, seg.Transcript);
                    }
                }
            };

            // StreamingDiarizeAsync takes any readable Stream.IO object, that is
            // only the Stream.IO.Read method needs to be implemented.
            using (FileStream file = File.OpenRead (audioFile)) {
                await client.StreamingDiarizeAsync (file, diarCfg, handleFunc);
            }
        }
    }
}

{{< /tab >}}

{{< /tabs >}}

### Streaming from microphone

Streaming audio from microphone input typically needs us to interact with system
libraries. There are several options available, and although the examples here use
one, you may choose to use an alternative as long as the recording audio format is
chosen correctly.

{{< tabs >}}

{{< tab "Python" "py">}}

# This example requires the pyaudio (http://people.csail.mit.edu/hubert/pyaudio/)
# module to stream audio from a microphone. Instructions for installing pyaudio
# for different systems are available at the link. On most platforms, this is
# simply `pip install pyaudio`

import juzu
import pyaudio
import threading

# This is Cobalt's demo server. Change it to your own local Juzusvr if you have one running.
serverAddress = 'demo.cobaltspeech.com:2727'

# set insecure=True for connecting to server not using TLS
client = juzu.Client(serverAddress, insecure=False)

# get list of available models
modelResp = client.ListModels()

# use the first available model
juzuModel = modelResp.models[0]
juzuModelID = juzuModel.id
cubicModelID = juzuModel.attributes.compatible_cubic_models[0]

# creating diarization config to transcribe + diarize
# audio stream from microphone
cfg = juzu.DiarizationConfig(
    model_id = juzuModelID,
    cubic_model_id = cubicModelID,
    num_speakers = 2,
    audio_encoding = "RAW_LINEAR16",
    sample_rate = juzuModel.attributes.sample_rate,
)

# client.StreamingDiarize takes any binary stream object that has a read(nBytes)
# method. The method should return nBytes from the stream. So pyaudio is a suitable
# library to use here for streaming audio from the microphone. Other libraries or
# modules may also be used as long as they have the read method or have been wrapped
# to do so.

# defining class to wrap around microphone stream from py audio
class MicStream(object):

    def __init__(self, sampleRate):

        self._p = pyaudio.PyAudio()
        # opening mic stream, recording 16 bit little endian integer samples, mono channel
        self._stream = self._p.open(format=pyaudio.paInt16, channels=1, rate=sampleRate, input=True)
        self._stopped = False

    def __del__(self):
        self._stream.close()
        self._p.terminate()

    # streamingDiarize requires a read(nBytes) method
    # that return a list of nBytes from the stream. An
    # empty list signals the end of stream.
    def read(self, nBytes):
        # if stream is stopped, return empty list to
        # signal end of stream to Juzu
        if self._stopped:
            return []
        return self._stream.read(nBytes)

    def pause(self):
        self._stream.stop_stream()

    def resume(self):
        self._stream.start_stream()

    def stop(self):
        self._stopped = True


audio = MicStream(juzuModel.attributes.sample_rate)

# helper function convert protobuf duration objects
# (which stores the time split into in integer seconds
# and integer nano seconds) into single floating value
# in seconds
def protoDurToSec(dur):
    return float(dur.seconds) + float(dur.nanos) * 1e-9

# starting thread to send streaming request to juzu
# and process results once they come back after the
# stream ends.
def streamToJuzu(cfg, audio):
    try:
        for resp in client.StreamingDiarize(cfg, audio):
            for result in resp.results:
                for segment in result.segments:
                    print("{start:.3f} - {end:.3f}\t{speaker}:\t{transcript}\n".format(
                        start = protoDurToSec(segment.start_time),
                        end = protoDurToSec(segment.end_time),
                        speaker = segment.speaker_label,
                        transcript = segment.transcript,
                        ))
    except Exception as ex:
        print("[error]: streaming diarization failed: {}".format(ex))

streamThread = threading.Thread(target=streamToJuzu, args=(cfg,audio))
streamThread.setDaemon(True)
streamThread.start()

# waiting for user to end mic stream
print("\nStreaming audio to Juzu server ...\n")
k = input("-- Press Enter key to stop stream --")

print("\nStopping Stream ...")
audio.stop()

print("Waiting for results ...")
streamThread.join()

{{< /tab >}}

{{< tab "C#" "md" >}}

// We do not currently have example C# code for streaming from a microphone. Simply
// pass the bytes from the microphone the same as is done from the file in the
// `Streaming from an audio file` example above via a class derived from `Stream.IO`
// (https://docs.microsoft.com/en-us/dotnet/api/system.io.stream) with the
// `int Read(buffer byte[], offset int, count int)` method implemented.
//
// See more at:
// https://docs.microsoft.com/en-us/dotnet/api/system.io.stream.read?view=net-5.0#System_IO_Stream_Read_System_Byte___System_Int32_System_Int32_

{{< /tab >}}

{{< /tabs >}}
