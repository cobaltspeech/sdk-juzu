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
        self.assertEqual(context.exception.details(),
                         "failed to connect to all addresses")

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
        self.assertEqual(context.exception.details(), "Socket closed")

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
        self.assertEqual(context.exception.details(), "Socket closed")

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
# openssl req -x509 -nodes -days 36500 -newkey rsa:2048 -keyout key.pem -out cert.pem
#
# and are used for testing only. Do not use in production.
certPem = b"""-----BEGIN CERTIFICATE-----
MIIDlTCCAn2gAwIBAgIUISO9AdBzEIxv366ruyRniPVrA8AwDQYJKoZIhvcNAQEL
BQAwWTELMAkGA1UEBhMCQVUxEzARBgNVBAgMClNvbWUtU3RhdGUxITAfBgNVBAoM
GEludGVybmV0IFdpZGdpdHMgUHR5IEx0ZDESMBAGA1UEAwwJbG9jYWxob3N0MCAX
DTIwMDExNDE4NTY0OVoYDzIxMTkxMjIxMTg1NjQ5WjBZMQswCQYDVQQGEwJBVTET
MBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50ZXJuZXQgV2lkZ2l0cyBQ
dHkgTHRkMRIwEAYDVQQDDAlsb2NhbGhvc3QwggEiMA0GCSqGSIb3DQEBAQUAA4IB
DwAwggEKAoIBAQDdF183r5B4ZeM677SYQrf7JWc26eBN60mqOSC18aeRxkeQsgO4
X3Y7wpQRcVdsm/i291otNIOWm4mJLLJFVE87Yq65gH4O4MHxQNlhZ0Bf1J8WsbsF
RHk3LF2rhUBll6cG+Z1OX7mCtmM33znXDFxTf3/DZ5XgeleNG98umeUMg8rHgj2y
UB4nwoMbeJIjk7e5tBQKCCNOYM1Mda1wzrvxo3blXsIFzpxLqQ+tVnYuql9CYjX1
69Nwq+Dsgv6zNWzWMlPTPAKdbOVVvXV2hfQ3LmnuzCv9t/TUdwkdyUMDUkbF+T8v
eD5bMP3k8lYuaNu0YQmbgbKvklK7voaEFte9AgMBAAGjUzBRMB0GA1UdDgQWBBTO
3ItCZrqn4cPtvGiVT1BQ+gzZhDAfBgNVHSMEGDAWgBTO3ItCZrqn4cPtvGiVT1BQ
+gzZhDAPBgNVHRMBAf8EBTADAQH/MA0GCSqGSIb3DQEBCwUAA4IBAQCGSvX5vN2F
siW23KnlurUtNiSoPzwSOwRGzaPcQYc1rdTXfN0F3Hj4qRnVJ+9jl3z6/xQnrgzg
iQ+4bZcnJeebmPI0jMZZXDXdDnp/Ze4klELpG53DzVzGZ7FvENmfNFIEIx6hrT2K
TDrWfZomjYu4Tn5rGAA4TflA9u8AWHYcZDLtjiwuHNFJLY6PdZSJU+OaPJXztRDV
jg/KVCsPH5LLxZy1U175YbWN7nIDvG2/H7o2vQdBs9A8lJ8CEGr2/jXTl8GDyOYN
6uHPipr8Y2rvh/jY8HalYo3x5gXM7AJk4OBh6x/Fcw6L0BmUHrUYvaU8UbwjraKF
PISX1spqbI8a
-----END CERTIFICATE-----"""

keyPem = b"""-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDdF183r5B4ZeM6
77SYQrf7JWc26eBN60mqOSC18aeRxkeQsgO4X3Y7wpQRcVdsm/i291otNIOWm4mJ
LLJFVE87Yq65gH4O4MHxQNlhZ0Bf1J8WsbsFRHk3LF2rhUBll6cG+Z1OX7mCtmM3
3znXDFxTf3/DZ5XgeleNG98umeUMg8rHgj2yUB4nwoMbeJIjk7e5tBQKCCNOYM1M
da1wzrvxo3blXsIFzpxLqQ+tVnYuql9CYjX169Nwq+Dsgv6zNWzWMlPTPAKdbOVV
vXV2hfQ3LmnuzCv9t/TUdwkdyUMDUkbF+T8veD5bMP3k8lYuaNu0YQmbgbKvklK7
voaEFte9AgMBAAECggEBAK18w4jMyQ7Q1KfQpOO9puT6Cq36g7pg4OMkBNkAkT9A
WbPfHDA3KG3oV4wAZluhYF8iZa6HQKKT1i6/1fu1Fp9A5l5Fx6UhFM6c1ncqMEeC
bnu+Z0TQ4FU9CRuoaknN4JEGmjt/vfAl8mFLVvW6i1AyAi1xQRhup/jgYBcPR76y
zcmgdyaDNX/Z4rPAIzPTJbEjPhV8L3J8SjM6CML9k03QA6GALyVIJwUOUy/qJtYL
9UrUFwE+jjbifWjIpNjdKPFi6ltgSk9cOQZXcKWw+0CBJ5q3cZ6m6faqnWyDMHqP
garQd8bHzxG5yrVqrYiP8Mv2vD+cham8Nl2frJsK4cECgYEA/yR5rTpaFkZ9fV12
/r/KLapClg9R6COwBWk5nPvGde7rLMPVg2sQgyzfj7fRvZiQDx9zjmn96LJX+uwr
do88rWFUCmq5S/G0p1UD3DEW1b/YcN5D/nnG0QjeP0PSmG4Aj/SxDUD8AX1zwCpn
/G7dJame4A45tW3oxlGDMUPIwHkCgYEA3dWZV2VIHk5Aa4PJ92iTCtIPQnVb7140
WL2cCmsvIWe/UkHvHJZYcgB78OOuNET3ijI83rx/Ry2iBapKnuHxbl8npOw+8btZ
GakR0Htbns/sOqfb1jruen6kw+HTZaYYYEckHw8DiARLc4PlC6WSsYhL8W2Jl+SD
mwBb4NJQKGUCgYBGcR2e9BNXPxL6f8mQwAbj4LQNliE5BFFezRR5ARJkERig/Vh/
thmS/dqjZU7lF6/+XOKcmSrfCg48WuQNEbLg85QuZBTQoOUNpe0w5+S0EwmA7/y5
z4lSwS4LLYCBUS2akSYo0J5DEw3YKl0XVsx7z37rwUGxk6zGxE6CVYKhkQKBgGGh
PyJqjcngsJNg5gM//+70MgkSs4pukGU51bH0KELwcRBXuk9/j59kvSdwXNveOn+U
yptQpEeEOtl5b+vrDqF/uWfpHW6wAG+9q/xwPgtwAMxz0dnAB/LbR9J50drbtcCx
rqEIr4ouMbK+KpDsptoBXUL87WBvDsip6MXSabrNAoGBAKKjpwZXPBUfjRMEWSrv
AQLrBifzWHi6NbC+yFxdbaidzOzOLfxMXUKW8MRWK9Wsqza1r/bubbVTakzzFuZa
sxjcfOmE0/ArjIxeMLVvaXbY1hCgwdplJf0JyksidT7njIkXL+vFwq6VD5R+JCC8
MqcOXU7V6srirJOObsIHyC/i
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
