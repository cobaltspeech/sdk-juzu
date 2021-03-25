---
title: "Setup Connection"
description: "Describes how to connect to a running Juzu server instance."
weight: 22
---

Once you have your Juzu server up and running, let's see how we can use the SDK
to connect to it.

<!--more-->

First, you need to know the address (`host`:`port`) where the server is running.
This document will assume the values `127.0.0.1:2727`, but be sure to change
those to point to your server instance.

## Default Connection

The following code snippet connects to the server and queries its version.  It
uses our recommended default setup, expecting the server to be listening on a
TLS encrypted connection. Examples showing how to connect to a server not using
TLS is also shown in the [Insecure Connection](#insecure-connection) section.

{{< tabs >}}

{{< tab "Python" "py" >}}

import juzu

serverAddress = "127.0.0.1:2727"

client = juzu.Client(serverAddress)

resp = client.Version()
print(resp)

{{< /tab >}}

{{< tab "C#" "csharp" >}}

var serverAddress = "127.0.0.1:2727";

// change this to true if using server without TLS
var insecure = false;
var client = new Client (serverAddress, insecure);

// Get Version of the server
var ver = client.Version ();
Console.WriteLine ("Juzu: {0} Server: {1}", ver.Juzu, ver.Server);

{{< /tab >}}

{{< /tabs >}}

## Insecure Connection

It is sometimes required to connect to Juzu server without TLS enabled, such as
during debugging.

Please note that if the server has TLS enabled, attempting to connect with an
insecure client will fail. To connect to an instance of Juzu server without TLS enabled, you
can use:

{{< tabs >}}

{{< tab "Python" "py" >}}

client = juzu.Client(serverAddress, insecure=True)

{{< /tab >}}

{{< tab "C#" "csharp" >}}

var insecure = true;
var client = new Client (serverAddress, insecure);

{{< /tab >}}

{{< /tabs >}}

## Client Authentication

In our recommended default setup, TLS is enabled in the gRPC setup, and when
connecting to the server, clients validate the server's SSL certificate to make
sure they are talking to the right party.  This is similar to how "https"
connections work in web browsers.

In some setups, it may be desired that the server should also validate clients
connecting to it and only respond to the ones it can verify. If your Juzu
server is configured to do client authentication, you will need to present the
appropriate certificate and key when connecting to it.

Please note that in the client-authentication mode, the client will still also
verify the server's certificate, and therefore this setup uses mutually
authenticated TLS. This can be done with:

{{< tabs >}}

{{< tab "Python" "py" >}}

client = juzu.Client(serverAddress, clientCertificate=certPem, clientKey=keyPem)

{{< /tab >}}

{{< tab "C#" "csharp" >}}

// Authenticating Server Certificate
var rootPem = File.ReadAllText("root.pem");
var client = new Client (serverAddress, rootPem);

// OR

// Mutual Authentication
var rootPem = File.ReadAllText("root.pem");
var certPem = File.ReadAllText("cert.pem");
var keyPem = File.ReadAllText("key.pem");
var client = new Client (serverAddress, rootPem, certPem, keyPem);

{{< /tab >}}

{{< /tabs >}}

where rootPem is the bytes of the certificate used to validate the server
certificate and certPem & keyPem are the bytes of the client certificate and key
provided to you respectively.
