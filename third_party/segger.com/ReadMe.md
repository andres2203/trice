# SEGGER downloaded Software

- Check in the Internet for newer versions.

## JLink

- Download and install [J-LinkSoftwareAndDocumentationPack](https://www.segger.com/downloads/jlink/#J-LinkSoftwareAndDocumentationPack) or simply use `JLinkRTTLogger.exe` and acompanying `JLinkARM.dll` copied from default install location `C:\Program Files (x86)\SEGGER\JLink`. Both files are inside `JLinkRTTLogger.zip` You need to put to a location in $PATH or extend $PATH.

## SEGGER_RTT

- Target code is expected inside SEGGER_RTT. This is the extracted SEGGER_RTT_V....zip.
- Optionally check for a newer version.

## STLinkReflash_190812.zip

- Tool for exchanging STLINK and JLINK software on STM32 evaluation boards.
  - Works not for v3 Hardware but well for v2 Hardware.
  - In case of not accepting the ST-Link firmware use [../st.com/en.stsw-link007_V2-37-26.zip](../st.com/en.stsw-link007_V2-37-26.zip) for updating the ST-Link firmware first. It could be you need to exchange the ST-Link firmware variant into the variant with mass storage.
