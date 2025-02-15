# MDK-ARM_LL_UART_RTT0_ESC_STM32F030R8-NUCLEO-64

Delivers TRICE logs over UART and RTT channel 0 in "esc" encoding:

## Usage
- Connect board over USB.
- Enter command `trice s` to discover the new COM port.
- Enter command `trice l -e esc -p COM3 -s` when the new COM port is COM3. The `-s` switch shows the raw arriving bytes for issue tracking.
- Enter command `trice l -e esc -p JLINK` when a (SEGGER) JLINK debug probe is active.
- Enter command `trice l -e esc -p STLINK` when a STLINK debug probe is active.
## Infos
- The trices are stored in `esc` format and read out in 8 byte packages.
- The transfer goes over RTT channel 0 and UART.
- A sync package is inserted from time to time in UART transfer.
