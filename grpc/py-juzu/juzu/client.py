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

import os
import sys
import grpc

from google.protobuf.empty_pb2 import Empty

sys.path.append(os.path.join(os.path.dirname(os.path.realpath(__file__))))
from juzu_pb2_grpc import JuzuStub
from juzu_pb2 import DiarizationConfig, DiarizationAudio
from juzu_pb2 import StreamingDiarizeRequest

class Client(object):
    """ A client for interacting with Cobalt's Juzu GRPC API."""

    def __init__(self, serverAddress, insecure=False,
                serverCertificate=None,
                clientCertificate=None,
                clientKey=None,
                bufferSize=8192):
        """  Creates a new juzu Client object.

        Args:
            serviceAddress: host:port of where juzu server is running

            insecure: If set to true, an insecure grpc channel is used.
                      Otherwise, a channel with transport security is used.

            serverCertificate:  PEM certificate as byte string which is used as a  
                                root certificate that can validate the certificate 
                                presented by the server we are connecting to. Use this 
                                when connecting to an instance of juzu server that is 
                                using a self-signed certificate.

            clientCertificate:  PEM certificate as bytes string presented by this Client when
                                connecting to a server. Use this when setting up mutually 
                                authenticated TLS. The clientKey must also be provided.

            clientKey:  PEM key as byte string presented by this Client when
                        connecting to a server. Use this when setting up mutually 
                        authenticated TLS. The clientCertificate must also be provided.

            bufferSize: Bytes of audio sent to juzu server in each chunk if using streaming
                        diarization. Default is 8192 bytes.
        """

        self.serverAddress = serverAddress
        self.insecure = insecure
        self.bufferSize = bufferSize

        if self.bufferSize <= 0:
            raise ValueError('buffer size must be greater than 0')
        
        # channel options enabling keepalive pings from the client
        self._channelOpts = []

        # After a duration of this time the client pings its peer to see if the
        # transport is still alive. Int valued, milliseconds. Set to 15 seconds.
        self._channelOpts.append(('grpc.keepalive_time_ms', 15 * 1000))

        # After waiting for a duration of this time, if the keepalive ping
        # sender does not receive the ping ack, it will close the transport. Int
        # valued, milliseconds. Set to 20 seconds.
        self._channelOpts.append(('grpc.keepalive_timeout_ms', 20 * 1000))

        # Is it permissible to send keepalive pings without any outstanding streams.
        # Int valued, 0(false) / 1(true).  
        self._channelOpts.append(('grpc.keepalive_permit_without_calls', 1))

        # How many pings can we send before needing to send a data frame or header
        # frame? (0 indicates that an infinite number of pings can be sent without
        # sending a data frame or header frame)
        self._channelOpts.append(('grpc.http2.max_pings_without_data', 0))

        # Minimum time between sending successive ping frames without receiving any
        # data frame, Int valued, milliseconds. Set to 6 minutes.
        self._channelOpts.append(('grpc.http2.min_time_between_pings_ms', 6 * 60 * 1000))

        if insecure:
            # no transport layer security (TLS)
            self._channel = grpc.insecure_channel(serverAddress, options=self._channelOpts)        
        else:
            # using a TLS endpoint with optional certificates for mutual authentication
            if clientCertificate is not None and clientKey is None:
                raise ValueError("client key must also be provided")
            if clientKey is not None and clientCertificate is None:
                raise ValueError("client certificate must also be provided")
            self._creds = grpc.ssl_channel_credentials(
                root_certificates=serverCertificate, 
                private_key=clientKey,
                certificate_chain=clientCertificate)
            self._channel = grpc.secure_channel(serverAddress, self._creds, options=self._channelOpts)

        self._client = JuzuStub(self._channel)

    def __del__(self):
        """ Closes and cleans up after Client. """
        try:
            self._channel.close()
        except AttributeError:
            # client wasn't fully instantiated, no channel to close
            pass

    def Version(self):
        """ Queries the server for its version. """
        return self._client.Version(Empty())

    def ListModels(self):
        """ Retrieves a list of available speech diarization models. """
        return self._client.ListModels(Empty())

    def StreamingDiarize(self, cfg, audio):
        """ Performs bidirectional streaming speech diarization. Results are received after 
        the all the data has been streamed.

        Args:
            cfg: DiarizationConfig object containing the model ID, audio encoding etc.

            audio: a binary stream of data to send to juzu. The object passed in should have a 
                    read(nBytes) method that returns nBytes from the binary stream. An example is
                    the object created using open('filepath', 'rb'); For streaming from a microphone,
                    the object could be a PyAudio() stream. Byte chunks are read from this stream
                    and sent to juzu server sequentially. The size of each chunk is equal to the 
                    buffer size set for the client.
        """

        stream = _audioStreamer(cfg, audio, self.bufferSize)
        for resp in self._client.StreamingDiarize(stream):
            yield resp
        
def _audioStreamer(cfg, audio, bufferSize):
    """ A generator that streams audio data packaged for diarization. The first 
    yield is a config message, and the following are all audio messages. """
    yield StreamingDiarizeRequest(config=cfg)
    while True:
        data = audio.read(bufferSize)
        if len(data) == 0:
            break
        rcgAudio = DiarizationAudio(data=data)
        yield StreamingDiarizeRequest(audio=rcgAudio)
