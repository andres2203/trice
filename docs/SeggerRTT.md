# Segger Real Time Transfer (RTT)

- Prerequisite is a processor with memory background access support like ARM Cortex-M cores.
- If you can use a Segger JLINK or an STM STLINK debug probe (ST Microelectronics eval boards have it) this is an easy and fast way to use trice without any UART or other port.
- Detailed description can be found in document [UM08001_JLink.pdf](https://github.com/rokath/trice/blob/master/third_party/segger.com/UM08001_JLink.pdf) in chapter 16 which is part of [https://www.segger.com/downloads/jlink/#J-LinkSoftwareAndDocumentationPack](https://www.segger.com/downloads/jlink/#J-LinkSoftwareAndDocumentationPack).
- Following examples are for Windows, but should work similar also on Linux and Darwin (MacOS).
- The trice tool can use the SEGGER RTT protocol in different ways.
  - Hardware paths:
    - Use [J-Link](https://www.segger.com/products/debug-probes/j-link/) or [J-Link OB (on-board)](https://www.segger.com/products/debug-probes/j-link/models/j-link-ob/).
      J-Link OB one is flashable to many ST Microelectronics evaluation boards (v2.0 link hardware) and for example also usable with NXP and Atmel
    - Use ST-Link with [gostlink](https://github.com/rokath/trice/blob/master/third_party/goST/ReadMe.md).
      It uses only one USB endpoint so debugging and trice output in parallel is not possible.
    - Use some other Debug-Probe with target memory access (support welcome)
  - RTT channel selection (on target and on host)
    - RECOMMENDED:
      - `trice l -p JLINK` or shorter `trice l` (default port is JLINK) starts in background a `JLinkRTTLogger.exe` which \
      connects to J-Link and writes to a logfile which in turn is read by the trice tool. On exit the `JLinkRTTLogger.exe` is killed automatically.\
      It expects a target sending messages over RTT channel 0 (other channels supported too).\
      It is possible to start several instances on different channels as well as on different targets.
      - `trice l -p STLINK` starts in background a `trice/third_party/goST/stRttLogger.exe` which connects to J-Link and\
      writes to a logfile which in turn is read by the trice tool. On exit the `stRttLogger.exe` is killed automatically.\
      It expects a target sending messages over RTT channel 0 (other channels supported too).\
      It is possible to start several instances on different channels as well as on different targets.
      - If you have the choice prefer J-Link. It allowes parallel debugging and trice output.

## J-Link option

- Prerequisite is a SEGGER J-Link debug probe or a development board with an on-board J-Link option.

### Convert a STM NUCLEO or DISCOVERY onboard ST-Link (valid for ST-Link v2 & v2.1, not for v3)

Following steps describe the needed action for a ST Microelectronics evaluation board and windows - adapt them to your environment.

#### First step (to do if some issues occur - otherwise you can skip it)

- Get & install [STM32 ST-LINK utility](https://www.st.com/en/development-tools/stsw-link004.html)
- Run from default install location `"C:\Program Files (x86)\STMicroelectronics\STM32 ST-LINKUtility\ST-LINK Utility\ST-LinkUpgrade.exe"`)
- Enable checkbox `Change Type` and select radio button `STM32 Debug+Mass storage + VCP`. *The `STM32Debug+ VCP` won´t be detected by Segger reflash utility.*
  ![ST-LINK-Upgrade.PNG](./README.media/ST-LINK-Upgrade.PNG)

#### Second step

- Check [Converting ST-LINK On-Board Into a J-Link](https://www.segger.com/products/debug-probes/j-link/models/other-j-links/st-link-on-board/)
- Use `STLinkReflash.exe` to convert NUCLEO from ST-Link on-board to J-Link on-board. *`STM32 Debug+ VCP` won´t be detected by Segger reflash utility.*

### Some SEGGER tools in short

- `JLink.exe` is the SEGGER J-Link commander. It starts the **J-Link driver/server** and the `trice` tool can connect to it:
  - Example:
    - Compile and flash `../test/MDK-ARM_LL_UART_RTT0_PACK_STM32F030R8-NUCLEO-64` project
    - Open a commandline and run:

      ```b
      trice log -p JLINK
      ```

    - trice outout is visible
    - With `h`alt and `g`o inside the Jlink window the MCU can get haltes and released
    - It is possible in parallel to debug-step with a debugger (tested with ARM-MDK)

- **PLUS:**
- works reliable
- trice can connect over TCP localhost:19021 and display logs over RTT channel 0

- **MINUS:**
- Uses RTT up-channel 0 and therefore RTT up-channel 0 is not usable differently
- No down-channel usable
- `JLinkRTTLogger.exe` is usable for writing RTT channel N data from target into a file. The trice tool can watch this file and display them.
  - **PLUS:**
    - Can use RTT up-channel 0,1,2,...
    - Can run in several instances for the same target each on a different channel
  - **MINUS:**
    - Needs a running **J-Link driver/server**
    - Logs in a file, so `trice` needs to read from that file
    - ...
  - **UNKNOWN:**
    - If several instances can watch on different RTT chanels and all goes in parallel to debugging.
  - **TESTED:**
    - Create file

    /c/Program\ Files\ \(x86\)/SEGGER/JLink/JLinkRTTLogger.exe -Device STM32F030R8 -if SWD -Speed 4000 -RTTChannel 0 triceRaw.log

    - Play file with *-device RTTF*

- `JLinkRTTClient.exe` can be used for simple text transmitting to the target, it also displays strings  from target coming over channel 0. It is not used by the trice tool.
  - **PLUS:**
    - target stimulation with proprietary protocol over RTT down-channel 0 possible
  - **MINUS:**
    - Unfortunately it cannot run separately parallel to stimulate the target with any proprietary protocol.

## Software Setup example (JLinkRTTViewer.exe is depreciated, use JLink.exe instead!)

- Build and flash `../test/MDK-ARM_LL_UART_RTT0_PACK_STM32F030R8-NUCLEO-64`
  - Download [J-Link Software and Documentation Pack](https://www.segger.com/downloads/jlink/#J-LinkSoftwareAndDocumentationPack) and install.
  - Run `trice log -port JLINK` or in short `trice l`.
  - Now the trice output is visible.
- In the SeggerRTT example projects you see how to setup.

<!---

- Start `"C:\Program Files (x86)\SEGGER\JLink\JLinkRTTViewer.exe"` and connect to the J-Link. You only need this as a running server to connect to.
  - Unfortunately the JLinkRTTViewer "steals" from time to time some trice data packages and displays them as data garbage.
  - Better use JLink.exe or the *Segger J-Link SDK* instead.

-->

## Segger RTT

- The main advantage here is, that no `triceServe()` nor any interrupt is needed in the background, because this job is automatically done by SeggerRTT. This way one can debug code as comfortable as with `printf()` but with all the TRICE advantages. Have a look here: ![SeggerRTTD.gif](./README.media/JLINK-DebugSession.gif)

## Segger RTT unbuffered (PC side receiver -device=JLINK)

- Avoid trice buffering inside target and write with TRICE macro directly into the RTT buffer
- Write the bytes per trace directly (little time & some space overhead on target, but no changes on host side) - this is implemented as test example.
  
  ![triceBlockDiagramWithSeggerRTT.svg](./README.media/triceBlockDiagramWithSeggerRTTD.svg)

## Segger J-Link SDK (~800 EUR) Option

- Segger offers a SeggerRTT SDK which allows to use more than just channel 0 and you can develop your own tooling with it.
- The `trice -port JLINK` is ok for usage **as is** right now. However if you wish more comfort check here:
- Question: [How-to-access-multiple-RTT-channels](https://forum.segger.com/index.php/Thread/6688-SOLVED-How-to-access-multiple-RTT-channels-from-Telnet/)
  - "Developer pack used to write your own program for the J-Link. Please be sure you agree to the terms of the associated license found on the Licensing Information tab before purchasing this SDK. You will benefit from six months of free email support from the time that this product is ordered."
  - [https://www.segger.com/products/debug-probes/j-link/technology/j-link-sdk/](https://www.segger.com/products/debug-probes/j-link/technology/j-link-sdk/)

## Additional Notes (leftovers)

- `Downloading RTT target package` from [https://www.segger.com/products/debug-probes/j-link/technology/about-real-time-transfer/](https://www.segger.com/products/debug-probes/j-link/technology/about-real-time-transfer/).
- Read the manual `"C:\Program Files (x86)\SEGGER\JLink\Doc\Manuals\UM08001_JLink.pdf"`.
- Extract `"C:\Program Files (x86)\SEGGER\JLink\Samples\RTT\SEGGER_RTT_V672b.zip"` to target project. Check for an update @ SEGGER.
- Add `SEGGER_RTTI.c` to target project
- Put into main(): `char* s = "Hi5!\n"; SEGGER_RTT_Write(0, s, strlen(s));`, compile & flash image
- Start `"C:\Program Files (x86)\SEGGER\JLink\JLinkRTTViewer.exe"` & connect.
- Reset target & `Hi5!` should be visible inside viewer window.
- Now `Hi5!` should be visible inside terminal window when resetting target board.
- Instead of `"C:\Program Files (x86)\SEGGER\JLink\JLinkRTTViewer.exe"` also `"C:\Program Files (x86)\SEGGER\JLink\JLink.exe"` could be started & then `https://github.com/stalehd/jlinklogviewer`.
- What also works: `"C:\Program Files (x86)\SEGGER\JLink\JLink.exe"` followed by a terminal program like TeraTerm connected to `localhost:19021`.
- [https://github.com/stalehd/jlinklogviewer](https://github.com/stalehd/jlinklogviewer) is integrated into the trice tool (`-port JLINK` option)

## Further development

- The GoST project offers a way around JLINK. Used -port STLINK instead.
- Maybe `libusb` together with `libjaylink` offer some options too.
- Checkout [https://github.com/deadsy/jaylink](https://github.com/deadsy/jaylink).
- `"C:\Program Files (x86)\SEGGER\JLink\JMem.exe"` shows a memory dump.

- Go to [https://libusb.info/](https://libusb.info/)
  - -> Downloads -> Latest Windows Binaries
  - extract `libusb-1.0.23` (or later version)

```b
libusb-1.0.23\examples\bin64> .\listdevs.exe
2109:2811 (bus 2, device 8) path: 6
1022:145f (bus 1, device 0)
1022:43d5 (bus 2, device 0)
0a12:0001 (bus 2, device 1) path: 13
1366:0105 (bus 2, device 10) path: 5
```

- Repeat without connected Segger JLink

```b
libusb-1.0.23\examples\bin64> .\listdevs.exe
2109:2811 (bus 2, device 8) path: 6
1022:145f (bus 1, device 0)
1022:43d5 (bus 2, device 0)
0a12:0001 (bus 2, device 1) path: 13
```

- In this case `1366:0105 (bus 2, device 10) path: 5` is missing, so `vid=1366`, `did=0105` as example
- On Windows install WSL2. The real Linux kernel is needed for full USB access.
