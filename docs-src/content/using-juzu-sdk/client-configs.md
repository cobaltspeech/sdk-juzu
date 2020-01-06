---
title: "Diarization Configurations"
description: "Describes the configurable options for diarization requests."
weight: 25
---

An in-depth explanation of the various fields of the complete SDK can be found [here](../../protobuf). The sub-section [DiarizationConfig](../../protobuf/#message-diarizationconfig) details the options that can be set when sending a diarization request.

This page here discusses the more common combinations sent to the server.

## Fields


| Field | Required | Default | Description |
| :-----: | :----: | :-----: | ----------- |
| model_id | Yes |  | <p>ID of the diarization model to use on the server. Can be obtained by first getting list of models on the server via ListModels().</p> |
| num_speakers | Yes |  | <p>The number of speakers expected in the audio; specifying the correct number of speakers improves the accuracy of the speaker labels. If the number of speakers is unknown, set to 0. </p> |
| audio_encoding | Yes |  | <p>Encoding of audio data sent/streamed through the `DiarizationAudio` messages. For encodings like WAV/FLAC that have headers, the headers are expected to be sent at the beginning of the stream, not in every `DiarizationAudio` message.</p> |
| sample_rate | Yes |  | <p>Sampling rate of the audio to process.</p> |
| cubic_model_id | Yes if transcription required | "" | <p>Unique identifier of the cubic model to be used for speech recognition. If this value is specified, transcription results from the cubic model with the given ID will also be returned alongside speaker labels. If it omitted or blank, the results will not include transcripts, even if Cubic server was included in the deployed image.</p> |
| enable_raw_transcript | No | False | <p>If true, the raw transcript (unformatted) will be included in the results (only has an effect if Cubicsvr also set up with Juzusvr).</p> |
