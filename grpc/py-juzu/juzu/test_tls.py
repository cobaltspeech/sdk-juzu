# -*- coding: utf-8 -*-
#
# Copyright(2019) Cobalt Speech and Language Inc.
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
import grpc
from concurrent import futures

from client import Client
import juzu_pb2
import juzu_pb2_grpc

expectedResponses = {}
expectedResponses['Version'] = juzu_pb2.VersionResponse(juzu='2727', server='v.27.0')


def setupGRPCServerWithTLS(certPem, keyPem, mutual=False):
    # create a gRPC server and adding the defined juzu servicer class to it
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=1))
    juzu_pb2_grpc.add_JuzuServicer_to_server(JuzuServicer(), server)
    # adding SSL credentials
    if mutual:
        creds = grpc.ssl_server_credentials(
            [(keyPem, certPem)], root_certificates=certPem, require_client_auth=True)
    else:
        creds = grpc.ssl_server_credentials([(keyPem, certPem)])
    serverAddress = 'localhost:2727'
    server.add_secure_port(serverAddress, creds)
    server.start()
    return server, serverAddress

# Mock server service implementations
class JuzuServicer(juzu_pb2_grpc.JuzuServicer):

    def Version(self, request, context):
        return expectedResponses['Version']


class TestTLS(unittest.TestCase):

    def test_ServerTLS(self):
        # TestServerTLS starts a server with TLS(not mutual) and makes sure only the
        # appropriate clients can connect to it. Since the server uses a self-signed
        # certificate, a default client should not be able to call methods, as the
        # certificate can not be validated.
        server, serverAddress = setupGRPCServerWithTLS(
            certPem, keyPem, mutual=False)

        # default client with self-signed server cert; should fail
        with self.assertRaises(grpc.RpcError) as context:
            client = Client(serverAddress, insecure=False)
            response = client.Version()
        self.assertEqual(context.exception.details(), "failed to connect to all addresses")

        # client provides incorrect server certificate; should fail
        with self.assertRaises(grpc.RpcError) as context:
            client = Client(serverAddress, insecure=False,
                            serverCertificate=fakeCertPem)
            response = client.Version()
        self.assertEqual(context.exception.details(), "failed to connect to all addresses")

        # client tries to connect with insecure channel; should fail
        with self.assertRaises(grpc.RpcError) as context:
            client = Client(serverAddress, insecure=True,
                            serverCertificate=certPem)
            response = client.Version()
        self.assertEqual(context.exception.details(), "failed to connect to all addresses")

        # client provides correct server certificate; should succeed
        client = Client(serverAddress, insecure=False,
                        serverCertificate=certPem)
        response = client.Version()
        self.assertEqual(response, expectedResponses['Version'])

        server.stop(0)

    def test_MutualTLS(self):
        # TestMutualTLS starts a server with mutual TLS enabled and makes
        # sure only the appropriate clients can connect to it.
        server, serverAddress = setupGRPCServerWithTLS(
            certPem, keyPem, mutual=True)

        # mutual tls with correct cert but wrong CA; should fail (client can not validate the server)
        with self.assertRaises(grpc.RpcError) as context:
            client = Client(serverAddress, insecure=False, serverCertificate=fakeCertPem,
                            clientCertificate=certPem, clientKey=keyPem)
            response = client.Version()
        self.assertEqual(context.exception.details(), "failed to connect to all addresses")

        # mutual tls with wrong cert but correct CA; should fail (server can not validate the client)
        with self.assertRaises(grpc.RpcError) as context:
            client = Client(serverAddress, insecure=False, serverCertificate=certPem,
                            clientCertificate=fakeCertPem, clientKey=fakeKeyPem)
            response = client.Version()
        self.assertEqual(context.exception.details(), "failed to connect to all addresses")

        # mutual tls with correct cert and CA but no key provided; should fail (client raises exception)
        with self.assertRaises(ValueError):
            client = Client(serverAddress, insecure=False,
                            serverCertificate=certPem, clientCertificate=certPem)
            response = client.Version()

        # client tries to connect with insecure channel; should fail
        with self.assertRaises(grpc.RpcError) as context:
            client = Client(serverAddress, insecure=True, serverCertificate=certPem,
                            clientCertificate=certPem, clientKey=keyPem)
            response = client.Version()
        self.assertEqual(context.exception.details(), "failed to connect to all addresses")

        # client presented a bad client cert/key; should fail
        with self.assertRaises(grpc.RpcError) as context:
            client = Client(serverAddress, insecure=False, serverCertificate=certPem,
                            clientCertificate=certPem[:3], clientKey=keyPem)
            response = client.Version()
        self.assertEqual(context.exception.details(), "Empty update")

        # client provides appropriate certificates; should succeed
        client = Client(serverAddress, insecure=False, serverCertificate=certPem,
                        clientCertificate=certPem, clientKey=keyPem)
        response = client.Version()
        self.assertEqual(response, expectedResponses['Version'])

        server.stop(0)


# The TLS certificates below were generated using:
#
# openssl req -new -x509 -newkey rsa:2048 -days 36500 -keyout key.pem -out cert.pem -nodes -subj "/CN=localhost" -addext "subjectAltName=DNS:localhost"
#
# and are used for testing only. Do not use in production.
certPem = b"""-----BEGIN CERTIFICATE-----
MIIDITCCAgmgAwIBAgIURF9ueWPW7wLSODqDuGMxFrVQIdYwDQYJKoZIhvcNAQEL
BQAwFDESMBAGA1UEAwwJbG9jYWxob3N0MCAXDTIxMDMyNTE1MjkwOFoYDzIxMjEw
MzAxMTUyOTA4WjAUMRIwEAYDVQQDDAlsb2NhbGhvc3QwggEiMA0GCSqGSIb3DQEB
AQUAA4IBDwAwggEKAoIBAQCb8iDrB27cKi4ObbPAth3gybv4d8baWa03RwgM9k7m
mG6cHD3ZKZykiMZhtn0QrmIzx+pUU7J3lCkYcpBuPREYtNjlwo17Btr0/B3duVeT
jA7hb2HZbL3h2mUeTArS0ljoJ+6qhdjBTeDE5isc/favceK/4p5Q6mLU7x1udeyP
RGY8p+zkzZ0u62YMVuuaz9NWAJ3wAhC29ci+8Ub7oJBP01Lk15yZlkipSRpKGixq
gcLZQtCQpCmZP8cIwvYPc8xvxMq0cl78tyr+AD+6jZo1wEgyS5l/gPgTjlTwVDn/
1mRZBXiwP2Mza/SkTk5JwlapFAuZ3J7mTCbmsHawtb8JAgMBAAGjaTBnMB0GA1Ud
DgQWBBQIoEPhVqkC4kjoMJnIz57YZHHbxTAfBgNVHSMEGDAWgBQIoEPhVqkC4kjo
MJnIz57YZHHbxTAPBgNVHRMBAf8EBTADAQH/MBQGA1UdEQQNMAuCCWxvY2FsaG9z
dDANBgkqhkiG9w0BAQsFAAOCAQEAeP43p+fc1IrDAdVb0l8Wt5NexPae+Xqigxm0
vfR74lokSWvGeFtXH1nX1bsfvC5ydwS7x1+j5iZsC5x1w0pXSbp9n/go6VOPhpLH
BfzlJrNxHxwd3R12B18tyE0Bn6NcXTr4xgAZOgRSXigQnkt+6nO6sa3EpH8C/Bm0
xFOSIvJzQr8ksqEvZnFS4TFNRt+Un/SfbYfIxn1U9YPeXswZg+c3RKCrtqAUULOb
6+zR6Qml1TTuQFWRmpY++j0HC/0YCaeACLgQF7b3Hy21j6kFe1uFhnZgRYO+tpS4
DWWJnuh3aXY43/Fz32Srdl7RWgw7Y6j4LGKgF4n07h7qcsdp+g==
-----END CERTIFICATE-----"""

keyPem = b"""-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQCb8iDrB27cKi4O
bbPAth3gybv4d8baWa03RwgM9k7mmG6cHD3ZKZykiMZhtn0QrmIzx+pUU7J3lCkY
cpBuPREYtNjlwo17Btr0/B3duVeTjA7hb2HZbL3h2mUeTArS0ljoJ+6qhdjBTeDE
5isc/favceK/4p5Q6mLU7x1udeyPRGY8p+zkzZ0u62YMVuuaz9NWAJ3wAhC29ci+
8Ub7oJBP01Lk15yZlkipSRpKGixqgcLZQtCQpCmZP8cIwvYPc8xvxMq0cl78tyr+
AD+6jZo1wEgyS5l/gPgTjlTwVDn/1mRZBXiwP2Mza/SkTk5JwlapFAuZ3J7mTCbm
sHawtb8JAgMBAAECggEANGCy5b3tYk1Ygsgd49BHYMThOachRx1vrG1TKudfwNvZ
6t0heNsdHbZCy3b7sqqfmTzAW+pwuvEsjGWOVTh1kZLMdv9NfdDjof08GsixvYEB
Tn6Woux/KmisJezshbrY2cDvzJ6AGw0JPZkBSg+S14Ks4j4/ZnYTE+nsPLmm5sdQ
RxLyVFhENzb7Rc5ILoLNQse1XaP+tBM7zh2NKEDuZgf8Upylom8g9XhHQDfLGxUv
JY0l4ytOCAbtxaVLv4uFqDUJ8SToQoVYEZVJ+bET5frIdKhiTKSDMbHyATFjiKug
Ebd7jHSE7z5XL7naQdcO50Fwna63+c8vcaHZeHyswQKBgQDIXoTMkL1Gc/dSxN7q
UpIpXOMQdpxdqPND4wMzfnBoXbCCLwhCKfzv2I5FpwceuWVVcRMrE566eA8fl+6k
pdoFWesMa9IRKfzIVovx8ok+wXSCUXhMCmFWcC/eQI94UNcwDxJU9QVcmW+oeLc+
H3s/Umr/Aa/d5hpqvsJOdCYYdQKBgQDHPinm4mKUbASnb7yxKlcc880GUVPg1GU8
aWG7xxWP7idwnbOICUrf5xdxlwlug2DhBio1yjsOy1luZoNDiixZNURLudKrPfwW
j0tq9hU5vPnGuNrZ19pdaf8Z17Ca3KSTFMhTjgIjlhYWgjLKEB485ltctlPbpo9R
nZCFzsuZxQKBgQCSqr1CkGA8Bf7wEzOpCi/7HwcGsgvbcyQ+eRmKw+68NnhRkuPU
qq9/UWEzEwRfNQ478L07Nu1lSSSo16nT06M3b45iVTo7XiSu8MsXm5Vi3nDOQ5xu
7No7T3t6kH90bphVkXK1xXTbR2U4SMS5MIoi281+RFnbb+e+Y0wP6W1JUQKBgQCF
yQIN3JlhHB8L5wWG2O1hrXtP8LAn2baYmBJeKlIaMUoeXmY8xDCYo+2kAkz8/g7B
SohxUffR+U11yjP91/vrcs8HIiGJhreM1m9HxTF/cZsuitW5TS19aD05NxGJCvsf
H+db45EgIBl4x7ge4i/LqMpUOMGaTf8SSIDQuIkSGQKBgDdF7tfXbJ6E8EdD3EMW
norwTQAFKoAyJlF7BUEVgjyZt1GKkbOZirO4QjtTgZk55TpesXpTcbN5Fl4I0Zll
QdhDNHP01pJ/E6AEfTWFnZGR0kzywhS8sPKjXGYgFXM05i9Fu4jaSKkkb3eCiW+V
w5EChlKjaj2TFOdi/2jHhMxF
-----END PRIVATE KEY-----"""
#
# define a new "fake" certificate and key that can be
# used in testing credential validation
#
fakeCertPem = b"""-----BEGIN CERTIFICATE-----
MIIBZjCCARCgAwIBAgIJAMn7swGSO8/UMA0GCSqGSIb3DQEBCwUAMA0xCzAJBgNV
BAYTAlVTMB4XDTE5MDMwNzEyMjAxM1oXDTE5MDQwNjEyMjAxM1owDTELMAkGA1UE
BhMCVVMwXDANBgkqhkiG9w0BAQEFAANLADBIAkEA2fY7pvsLKU0SKjYEat8lbU8n
OFVicZSuqKVuE/n/0PW69MLAiA/8nNgX/RGads1udSNe1LptlYH8scE76qi4VwID
AQABo1MwUTAdBgNVHQ4EFgQUEN5OTssGhgt04gmfYEPGhBPRVLwwHwYDVR0jBBgw
FoAUEN5OTssGhgt04gmfYEPGhBPRVLwwDwYDVR0TAQH/BAUwAwEB/zANBgkqhkiG
9w0BAQsFAANBALF4l2Szu9agWyu5LtOxd4I7eyMfyPZhqa0Lff5j/VdyTjIVDrgo
XKJaldxh9WV0fSXLIqczEqmAQbmbj/CTvz8=
-----END CERTIFICATE-----"""

fakeKeyPem = b"""-----BEGIN PRIVATE KEY-----
MIIBVQIBADANBgkqhkiG9w0BAQEFAASCAT8wggE7AgEAAkEA2fY7pvsLKU0SKjYE
at8lbU8nOFVicZSuqKVuE/n/0PW69MLAiA/8nNgX/RGads1udSNe1LptlYH8scE7
6qi4VwIDAQABAkBduVEbU3YQM3DtL78kiYHZiCDQS38CYjHcmQ5Fjsne+xBcDJxZ
8dtnCu4/OVNUevTGAObATIQBW5eBieUf8tZRAiEA7F2Bx/p+N8fzjcyCifMTUyOk
qpRXPgwc/wQ4KfdvwTkCIQDsEVztoPr1wnnmBRFhZsXzOvzGpSMczm2x3EmvIA2W
DwIgIvgym0OUKOyMPA5lwcMUuNgtJI+N2MAyCgi1xn+1KQECIQC/YXgwIgEy+n4u
r88eYt56SUkilkB4GxatSgTmmBrLmwIhAJ66U/pu4nL2HIhGcUdxawRC0DJRwRKc
B6KD9XmVFWXX
-----END PRIVATE KEY-----"""


if __name__ == "__main__":
    unittest.main()
