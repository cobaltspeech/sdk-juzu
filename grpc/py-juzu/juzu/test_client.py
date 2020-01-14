# -*- coding: utf-8 -*-
#
# Copyright(2020) Cobalt Speech and Language Inc.
#
# Licensed under the Apache License, Version 2.0 (the "License")
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
# http: // www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

import unittest
import io
import grpc
from concurrent import futures

from juzu_types import DurationMilliseconds
from client import Client, DiarizationAudio, DiarizationConfig
import juzu_pb2
import juzu_pb2_grpc

expectedResponses = {}
expectedResponses['Version'] = juzu_pb2.VersionResponse(juzu='2727', server='v.27.0')
expectedResponses['ListModels'] = juzu_pb2.ListModelsResponse(models=[juzu_pb2.Model(id="1"), juzu_pb2.Model(id="2")])

diarizationResult = [
    juzu_pb2.DiarizationResult(
        segments=[
            juzu_pb2.Segment(
                speaker_label="1",
                start_time=DurationMilliseconds(1000),
                end_time=DurationMilliseconds(3000),
                transcript="Hello Goodbye",
                words=[
                    juzu_pb2.WordInfo(word="Hello", start_time=DurationMilliseconds(1000), duration=DurationMilliseconds(1000), confidence=1.0),
                    juzu_pb2.WordInfo(word="Goodbye", start_time=DurationMilliseconds(2000), duration=DurationMilliseconds(1000), confidence=1.0),
                ]
            ),
            juzu_pb2.Segment(
                speaker_label="2",
                start_time=DurationMilliseconds(3300),
                end_time=DurationMilliseconds(5300),
                transcript="Hello Goodbye",
                words=[
                    juzu_pb2.WordInfo(word="Hello", start_time=DurationMilliseconds(3300), duration=DurationMilliseconds(1000), confidence=1.0),
                    juzu_pb2.WordInfo(word="Goodbye", start_time=DurationMilliseconds(4300), duration=DurationMilliseconds(1000), confidence=1.0),
                ]
            ),
        ],

        speaker_labels=[
            "1", "2"
        ],

        is_partial=False,
    )
]

expectedResponses['StreamingDiarize'] = juzu_pb2.DiarizationResponse(results=diarizationResult)

def setupGRPCServer():
    # create a gRPC server and adding the defined juzu servicer class to it
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
    juzu_pb2_grpc.add_JuzuServicer_to_server(JuzuServicer(), server)
    serverAddress = 'localhost:2727'
    server.add_insecure_port(serverAddress)
    server.start()
    return server, serverAddress

# Mock server service implementations
class JuzuServicer(juzu_pb2_grpc.JuzuServicer):
    
    def Version(self, request, context):
        return expectedResponses['Version']

    def ListModels(self, request, context):
        return expectedResponses['ListModels']

    def StreamingDiarize(self, request_iterator, context):
        
        # first message must be config
        request = next(request_iterator)
        if request.config == juzu_pb2.DiarizationConfig(): # empty config message
            context.set_code(grpc.StatusCode.FAILED_PRECONDITION)
            context.set_details('streamingdiarize failed: first message should be a config message')
            return juzu_pb2.DiarizationResponse()
        # rest should be audio messages
        try:
            while True:
                request = next(request_iterator)
                if request.audio == juzu_pb2.DiarizationAudio():   # empty audio message
                    context.set_code(grpc.StatusCode.INVALID_ARGUMENT)
                    context.set_details('streamingdiarize failed: all messages after the first should be audio messages')
                    return juzu_pb2.DiarizationResponse()
        except StopIteration:
            yield expectedResponses['StreamingDiarize']
        except Exception:
            context.set_code(grpc.StatusCode.INTERNAL)
            context.set_details('streamingdiarize failed: unexpected exception arose')
            return juzu_pb2.DiarizationResponse()

class TestClient(unittest.TestCase):
    
    @classmethod
    def setUpClass(cls):
        cls.server, cls.serverAddress = setupGRPCServer()

    @classmethod
    def tearDownClass(cls):
        cls.server.stop(0)

    def test_Version(self):
        client = Client(self.serverAddress, insecure=True)
        response = client.Version()
        self.assertEqual(response, expectedResponses['Version'])

    def test_ListModels(self):
        client = Client(self.serverAddress, insecure=True)
        response = client.ListModels()
        self.assertEqual(response, expectedResponses['ListModels'])

    def test_StreamingDiarize(self):
        client = Client(self.serverAddress, insecure=True)
        audio = io.BytesIO(b"0"*8192*5)

        cfg = DiarizationConfig(
            model_id="1",
            audio_encoding="WAV"
        )
        for response in client.StreamingDiarize(cfg, audio):
            self.assertEqual(response, expectedResponses['StreamingDiarize'])

if __name__ == "__main__":
    unittest.main()
