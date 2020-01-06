---
title: "Juzu SDK Documentation"
---

# Juzu API Overview

Juzu is Cobalt's speaker diarization engine. It can be deployed on-prem and accessed over the network or on your local machine via an API. We currently support C# and are adding support for more languages.

Once running, Juzu's API provides a method to which you can stream audio. This audio can either be from a microphone or a file. We recommend uncompressed WAV or lossless compression such as FLAC as the encoding, but we can support other formats as well upon request.

Juzu's API returns the diarization results using Google's protobuf library, allowing them to be handled natively by your application. Juzu returns timestamps for each speaker segment along with the speaker labels. Juzu can also return speaker transcripts (via Cubic) and word-level timestamps. Juzu's output format is described further below.

## Diarization Results

Juzu uses protocol buffers to return its results. The fields within these results are listed in the [Juzu Protobuf API Docs](protobuf/#message-diarizationresult). Depending on the programming language, the field names may vary in casing.

<details>
<summary><font color="236ecc"><b>Click here to see an example json representation of Juzu's output along with transcription and word level timestamps</b></font></summary>

``` json
{
    "speaker_labels": [
        "S-0",
        "S-1"
    ],
    "segments": [
        {
            "speaker_label": "S-0",
            "start_time": {
                "seconds": 1,
                "nanos": 740000000
            },
            "end_time": {
                "seconds": 5,
                "nanos": 120000000
            },
            "transcript": "My name is Michael. What can I help you with today",
            "words": [
                {
                    "word": "My",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 2,
                        "nanos": 880000000
                    },
                    "duration": {
                        "nanos": 150000000
                    }
                },
                {
                    "word": "name",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 3,
                        "nanos": 30000000
                    },
                    "duration": {
                        "nanos": 210000000
                    }
                },
                {
                    "word": "is",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 3,
                        "nanos": 240000000
                    },
                    "duration": {
                        "nanos": 90000000
                    }
                },
                {
                    "word": "Michael.",
                    "confidence": 0.99,
                    "start_time": {
                        "seconds": 3,
                        "nanos": 330000000
                    },
                    "duration": {
                        "nanos": 300000000
                    }
                },
                {
                    "word": "What",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 3,
                        "nanos": 690000000
                    },
                    "duration": {
                        "nanos": 150000000
                    }
                },
                {
                    "word": "can",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 3,
                        "nanos": 840000000
                    },
                    "duration": {
                        "nanos": 150000000
                    }
                },
                {
                    "word": "I",
                    "confidence": 0.997,
                    "start_time": {
                        "seconds": 3,
                        "nanos": 990000000
                    },
                    "duration": {
                        "nanos": 60000000
                    }
                },
                {
                    "word": "help",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 4,
                        "nanos": 50000000
                    },
                    "duration": {
                        "nanos": 150000000
                    }
                },
                {
                    "word": "you",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 4,
                        "nanos": 200000000
                    },
                    "duration": {
                        "nanos": 90000000
                    }
                },
                {
                    "word": "with",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 4,
                        "nanos": 290000000
                    },
                    "duration": {
                        "nanos": 120000000
                    }
                },
                {
                    "word": "today",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 4,
                        "nanos": 410000000
                    },
                    "duration": {
                        "nanos": 270000000
                    }
                }
            ]
        },
        {
            "speaker_label": "S-1",
            "start_time": {
                "seconds": 5,
                "nanos": 120000000
            },
            "end_time": {
                "seconds": 7,
                "nanos": 410000000
            },
            "transcript": "Hi I need to upgrade my service.",
            "words": [
                {
                    "word": "Hi",
                    "confidence": 0.689,
                    "start_time": {
                        "seconds": 5,
                        "nanos": 316000000
                    },
                    "duration": {
                        "nanos": 159000000
                    }
                },
                {
                    "word": "I",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 5,
                        "nanos": 679000000
                    },
                    "duration": {
                        "nanos": 141000000
                    }
                },
                {
                    "word": "need",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 5,
                        "nanos": 820000000
                    },
                    "duration": {
                        "nanos": 180000000
                    }
                },
                {
                    "word": "to",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 6
                    },
                    "duration": {
                        "nanos": 240000000
                    }
                },
                {
                    "word": "upgrade",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 6,
                        "nanos": 240000000
                    },
                    "duration": {
                        "nanos": 448000000
                    }
                },
                {
                    "word": "my",
                    "confidence": 0.887,
                    "start_time": {
                        "seconds": 6,
                        "nanos": 768000000
                    },
                    "duration": {
                        "nanos": 137000000
                    }
                },
                {
                    "word": "service.",
                    "confidence": 1,
                    "start_time": {
                        "seconds": 6,
                        "nanos": 933000000
                    },
                    "duration": {
                        "nanos": 477000000
                    }
                }
            ]
        }
    ]
}
```

</details>

Juzu can handle both short and long audio streams. The diarization process is initiated with the start of the stream and the results are returned within a few seconds of the stream ending.

## Obtaining Juzu

Cobalt will provide you with a package of Juzu that contains the engine,
appropriate diarization models and a server application. The package may include
Cubic as well for transcription and aiding the diarization process.  This server
exports Juzu's functionality over the gRPC protocol.  The
https://github.com/cobaltspeech/sdk-juzu repository contains the SDK that you
can use in your application to communicate with the Juzu server. This SDK is
currently available for C# and we would be happy to talk to you if you need
support for other languages. Most of the core SDK is generated automatically
using the gRPC tools, and Cobalt provides a top level package for more
convenient API calls.
