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

{{%tabs %}}

{{% tab "C#" %}}

#### Program.cs

``` csharp
using System;
using System.IO;
using System.Net;
using System.Text;
using System.Threading.Tasks;

namespace JuzusvrClient {
    class Program {

        static async Task Main (string[] args) {

            var url = "127.0.0.1:2727";
            string audioFile = "test.wav";

            var insecure = true;
            var client = new Client (url, insecure);

            // Getting list of diarization models on the server
            var modelResp = client.ListModels ();
            Console.WriteLine ("\nAvailable models:\n");
            foreach (var model in modelResp.Models) {
                Console.WriteLine ("{0}\t{1}\t{2}\n", model.Id, model.Name, model.Attributes.SampleRate);
            }

            // Creating config for Diarizing + Transcribing file with the first
            // Juzu Model available and the Cubic model with ID "1" (assigned by
            // cubicsvr config).
            var diarCfg = new DiarizationConfig {
                JuzuModelID = modelList.Models[0].Id,
                CubicModelID = "1",
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
```

#### JuzusvrClient.csproj

``` csharp
<Project Sdk="Microsoft.NET.Sdk">

  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>netcoreapp3.0</TargetFramework>
  </PropertyGroup>

  <ItemGroup>
    <PackageReference Include="Juzu-SDK" Version="0.9.3" />
  </ItemGroup>

</Project>
```

{{% /tab %}}

{{%/tabs %}}

### Streaming from microphone

Streaming audio from microphone input typically needs us to interact with system
libraries. There are several options available, and although the examples here use
one, you may choose to use an alternative as long as the recording audio format is
chosen correctly.

{{%tabs %}}

{{% tab "C#" %}}

We do not currently have example C# code for streaming from a microphone. Simply
pass the bytes from the microphone the same as is done from the file in the
[`Streaming from an audio file`](#streaming-from-an-audio-file) example above via
a class derived from [`Stream.IO`](https://docs.microsoft.com/en-us/dotnet/api/system.io.stream).
with the [`int Read(buffer byte[], offset int, count int)`](https://docs.microsoft.com/en-us/dotnet/api/system.io.stream.read#System_IO_Stream_Read_System_Byte___System_Int32_System_Int32_) method implemented.
{{% /tab %}}

{{%/tabs %}}
