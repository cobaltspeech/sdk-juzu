---
title: "Juzu API Reference"
weight: 3
---

The Juzu API is specified as a [proto
file](https://github.com/cobaltspeech/sdk-juzu/blob/master/grpc/juzu.proto).
This section of the documentation is auto-generated from the spec.  It describes
the data types and functions defined in the spec. The "messages" below
correspond to the data structures to be used, and the "service" contains the
methods that can be called.





## juzu.proto





### Service: Juzu
Service that implements the Cobalt Juzu Diarization API.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Version | .google.protobuf.Empty | VersionResponse | Queries the Version of the Server. |
| ListModels | .google.protobuf.Empty | ListModelsResponse | Retrieves a list of available diarization models. |
| StreamingDiarize | StreamingDiarizeRequest | DiarizationResponse | Performs bidirectional streaming to enable on-the-go processing of audio files, as well as the option to receive partial transcripts of audio along with speaker IDs. This method is not truly streaming for diarization yet, as results are received after specific chunks of audio have been sent. This method is only available via GRPC and not via HTTP+JSON. However, a web browser may use websockets to use this service. |

 <!-- end services -->



### Message: DiarizationAudio
Audio to be sent to the diarizer.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| data | bytes |  | <p></p> |







### Message: DiarizationConfig
Configuration for setting up a Diarizer.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| model_id | string |  | <p>ID of the diarization model to use on the server. Can be obtained by first getting list of models on the server via ListModels().</p> |
| num_speakers | uint32 |  | <p>The number of speakers expected in the audio; If the number of speakers is unknown, set to 0.</p> |
| sample_rate | uint32 |  | <p>Sampling rate of the audio to process.</p> |
| audio_encoding | DiarizationConfig.Encoding |  | <p>Encoding of audio data sent/streamed through the `DiarizationAudio` messages. For encodings like WAV/MP3 that have headers, the headers are expected to be sent at the beginning of the stream, not in every `DiarizationAudio` message.</p><p>If not specified, the default encoding is RAW_LINEAR16.</p><p>Depending on how they are configured, server instances of this service may not support all the encodings enumerated above. They are always required to accept RAW_LINEAR16. If any other `Encoding` is specified, and it is not available on the server being used, the recognition request will result in an appropriate error message.</p> |
| cubic_model_id | string |  | <p>Unique identifier of the cubic model to be used for speech recognition. If this value is specified, transcription results from the cubic model with the given ID will also be returned alongside speaker labels. If it omitted or blank, the results will not include transcripts, even if Cubic server was included in the deployed image.</p> |
| enable_raw_transcript | bool |  | <p>Returns unformatted transcript.</p> |







### Message: DiarizationResponse
Collection of sequence of diarization results in a portion of audio.
Juzu currently requires the full audio to determine which audio segments
belong to which speaker.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| results | DiarizationResult | repeated | <p></p> |







### Message: DiarizationResult
A diarization result corresponding to a portion of audio.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| segments | Segment | repeated | <p>Diarized segments containing speaker labels, timestamps and transcripts.</p> |
| speaker_labels | string | repeated | <p>Set of labels used to identify speakers in each segment.</p> |
| is_partial | bool |  | <p>If this is set to true, it denotes that the result is an interim partial result, and could change after more audio is processed. If unset, or set to false, it denotes that this is a final result and will not change.</p><p>Servers are not required to implement support for returning partial results, and clients should generally not depend on their availability.</p> |







### Message: ListModelsResponse
The message sent by the server for the `ListModels` method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| models | Model | repeated | <p>List of models available for use that match the request.</p> |







### Message: Model
Description of a Juzu Diarization Model.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| id | string |  | <p>Unique identifier of the model. This identifier is used to choose the model that should be used for diarization, and is specified in the `DiarizationConfig` message.</p> |
| name | string |  | <p>Model name. This is a concise name describing the model, and maybe presented to the end-user, for example, to help choose which model to use.</p> |
| attributes | ModelAttributes |  | <p>Model attributes.</p> |







### Message: ModelAttributes
Attributes of a Juzu Diarization Model.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| sample_rate | uint32 |  | <p>Audio sample rate supported by the model.</p> |
| segmentation_type | string |  | <p>The type of segmentation (fixed / variable) supported by the model.</p> |







### Message: Segment
A diarized segment of audio.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| speaker_label | string |  | <p>The identity of the speaker for this segment.</p> |
| start_time | google.protobuf.Duration |  | <p>Time offset relative to the beginning of audio received by the diarizer and corresponding to the start of this segment.</p> |
| end_time | google.protobuf.Duration |  | <p>Time offset relative to the beginning of audio received by the diarizer and corresponding to the end of this segment.</p> |
| transcript | string |  | <p>Text representing the transcription of the words that the speaker spoke. Formatting options are set in cubicsvr.</p> |
| words | WordInfo | repeated | <p>Words in the transcript, their timestamps and confidence scores.</p> |







### Message: StreamingDiarizeRequest
The top-level message sent by the client for the `StreamingDiarize`
request.  Multiple `StreamingDiarizeRequest` messages are sent. The first
message must contain a `DiarizationConfig` message only, and all subsequent
messages must contain `DiarizationAudio ` only.  All `DiarizationAudio `
messages must contain non-empty audio.  If audio content is empty, the server
may interpret it as end of stream and stop accepting any further messages.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| config | DiarizationConfig |  | <p></p> |
| audio | DiarizationAudio |  | <p></p> |







### Message: VersionResponse
The message sent by the server for the `Version` method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| juzu | string |  | <p>version of the juzu library handling the recognition.</p> |
| server | string |  | <p>version of the server handling these requests.</p> |







### Message: WordInfo
Word-specific information for recognized words.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| word | string |  | <p>The actual word in the text.</p> |
| confidence | double |  | <p>Confidence estimate between 0 and 1. A higher number represents a higher likelihood that the word was correctly recognized.</p> |
| start_time | google.protobuf.Duration |  | <p>Time offset relative to the beginning of audio received by the recognizer and corresponding to the start of this spoken word.</p> |
| duration | google.protobuf.Duration |  | <p>Duration of the current word in the spoken audio.</p> |





 <!-- end messages -->



### Enum: DiarizationConfig.Encoding
The encoding of the audio data to be sent for recognition.

For best results, the audio source should be captured and transmitted using
the RAW_LINEAR16 encoding.

| Name | Number | Description |
| ---- | ------ | ----------- |
| RAW_LINEAR16 | 0 | Raw (headerless) Uncompressed 16-bit signed little endian samples (linear PCM), single channel, sampled at the rate expected by the chosen `Model`. |
| WAV | 1 | WAV (data with RIFF headers), with data sampled at a rate equal to or higher than the sample rate expected by the chosen Model. |
| FLAC | 2 | FLAC data, sampled at a rate equal to or higher than the sample rate expected by the chosen Model. |


 <!-- end enums -->

 <!-- end HasExtensions -->




## Scalar Value Types

| .proto Type | Notes | Go Type | Python Type |
| ----------- | ----- | ------- | ----------- |
| double |  | float64 | float |
| float |  | float32 | float |
| int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int |
| int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | int/long |
| uint32 | Uses variable-length encoding. | uint32 | int/long |
| uint64 | Uses variable-length encoding. | uint64 | int/long |
| sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int |
| sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | int/long |
| fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int |
| fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | int/long |
| sfixed32 | Always four bytes. | int32 | int |
| sfixed64 | Always eight bytes. | int64 | int/long |
| bool |  | bool | boolean |
| string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | str/unicode |
| bytes | May contain any arbitrary sequence of bytes. | []byte | str |

