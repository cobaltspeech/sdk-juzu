---
title: "Installing the SDK"
description: "Gives language-specific instructions about how to add the SDK to your project."
weight: 21
---

Instructions for installing the SDK are language specific.

<!--more-->

### Python

The Python SDK depends on Python >= 3.5. You may use pip to perform a system-wide install, or use virtualenv for a local install.

```bash
pip install --upgrade pip
pip install "git+https://github.com/cobaltspeech/sdk-juzu#egg=cobalt-juzu&subdirectory=grpc/py-juzu"
```

### C\#

The C# SDK utilizes the [NuGet package manager](https://www.nuget.org).  The package is called `Juzu-SDK`, under the owners name of `CobaltSpeech`.

NuGet allows 4 different ways to install.  Further instructions can be found on the [nuget webpage](https://www.nuget.org/packages/Juzu-SDK/).  Installing via the `dotnet` cli through the command:

``` bash
dotnet add package Juzu-SDK
```

You can include the SDK in your *.csproj file as well:

```csharp
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
