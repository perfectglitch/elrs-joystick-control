// SPDX-FileCopyrightText: © 2023 OneEyeFPV oneeyefpv@gmail.com
// SPDX-License-Identifier: GPL-3.0-or-later
// SPDX-License-Identifier: FS-0.9-or-later

package crossfire

import (
	"github.com/kaack/elrs-joystick-control/pkg/crc"
	"github.com/kaack/elrs-joystick-control/pkg/util"
	"time"
)

func CreateModelIDFrame(modelId uint8) []uint8 {
	frame := []uint8{
		/* 0: */ uint8(UartSyncFrame),
		/* 1: */ 8,
		/* 2: */ uint8(CommandFrame),
		/* 3: */ uint8(ModuleEndpoint),
		/* 4: */ uint8(HandsetEndpoint),
		/* 5: */ uint8(SubcommandFrame),
		/* 6: */ uint8(CmdModelSelectFrame),
		/* 7: */ modelId, //model id
		/* 8: */ 0, //crc BA
		/* 9: */ 0, //crc D5
	}

	frame[8] = crc.BA(frame[2:8])
	frame[9] = crc.D5(frame[2:9])
	return frame
}

func CreatePingDevicesFrame() []uint8 {
	frame := []uint8{
		/* 0: */ uint8(UartSyncFrame),
		/* 1: */ 5,
		/* 2: */ uint8(PingDevicesFrame),
		/* 3: */ uint8(AllEndpoint),
		/* 4: */ uint8(LuaEndpoint),
		/* 5: */ 0, //crc BA
		/* 6: */ 0, //crc D5
	}

	frame[5] = crc.BA(frame[2:5])
	frame[6] = crc.D5(frame[2:6])
	return frame
}

func CreateParameterSettingsReadFrame(deviceId uint8, fieldId uint8, fieldChunk uint8) []uint8 {
	frame := []uint8{
		/* 0: */ uint8(UartSyncFrame),
		/* 1: */ 7,
		/* 2: */ uint8(ParameterSettingsReadFrame),
		/* 3: */ deviceId,
		/* 4: */ uint8(LuaEndpoint),
		/* 5: */ fieldId,
		/* 6: */ fieldChunk,
		/* 7: */ 0, //crc BA
		/* 8: */ 0, //crc D5
	}

	frame[7] = crc.BA(frame[2:7])
	frame[8] = crc.D5(frame[2:8])
	return frame
}

func CreateParameterSettingWriteFrameUint8(deviceId uint8, fieldId uint8, fieldValue uint8) []uint8 {
	frame := []uint8{
		/* 0: */ uint8(UartSyncFrame),
		/* 1: */ 7,
		/* 2: */ uint8(ParameterSettingsWriteFrame),
		/* 3: */ deviceId,
		/* 4: */ uint8(LuaEndpoint),
		/* 5: */ fieldId,
		/* 6: */ fieldValue,
		/* 7: */ 0, //crc BA
		/* 8: */ 0, //crc D5
	}

	frame[7] = crc.BA(frame[2:7])
	frame[8] = crc.D5(frame[2:8])
	//fmt.Printf("%x\n", frame)
	return frame
}

// CreateParameterSettingWriteFrameUint16 encodes a 16-bit value big-endian.
// Note: ELRS firmware only acts on the most-significant byte for most parameters;
// TBS Crossfire reads both bytes correctly.
func CreateParameterSettingWriteFrameUint16(deviceId uint8, fieldId uint8, fieldValue uint16) []uint8 {
	frame := []uint8{
		/* 0: */ uint8(UartSyncFrame),
		/* 1: */ 8,
		/* 2: */ uint8(ParameterSettingsWriteFrame),
		/* 3: */ deviceId,
		/* 4: */ uint8(LuaEndpoint),
		/* 5: */ fieldId,
		/* 6: */ uint8((fieldValue >> 8) & 0xFF),
		/* 7: */ uint8(fieldValue & 0xFF),
		/* 8: */ 0, //crc BA
		/* 9: */ 0, //crc D5
	}

	frame[8] = crc.BA(frame[2:8])
	frame[9] = crc.D5(frame[2:9])
	return frame
}

func GetRefreshRate(baudRate int32) time.Duration {
	if baudRate <= 115200 {
		return 16 * 1000 * time.Microsecond
	}

	if baudRate <= 400000 {
		return 4 * 1000 * time.Microsecond
	}

	//921600, 1870000, 3750000, 5250000
	return MinRefreshRate
}

// GetRefreshRateForModule returns the appropriate send interval for the given baud
// rate and module type.  TBS Crossfire runs at 150 Hz on the 400 kbps link
// (≈6.67 ms), while ELRS uses 250 Hz (4 ms) at that baud rate.
func GetRefreshRateForModule(baudRate int32, moduleType ModuleType) time.Duration {
	if baudRate <= 115200 {
		return 16 * 1000 * time.Microsecond
	}

	if baudRate <= 400000 {
		if moduleType == ModuleTypeCrossfire {
			return 6667 * time.Microsecond // 150 Hz
		}
		return 4 * 1000 * time.Microsecond // 250 Hz
	}

	return MinRefreshRate
}

// --- Module-type-aware extended frame builders ---
// TBS Crossfire extended frames use a single D5 CRC (no preceding BA CRC).
// ELRS extended frames use both a BA CRC and a D5 CRC.
// The helpers below select the correct format automatically.

// CreatePingDevicesFrameForModule creates a PingDevices frame with CRC encoding
// appropriate for the detected module type.
func CreatePingDevicesFrameForModule(moduleType ModuleType) []uint8 {
	if moduleType == ModuleTypeCrossfire {
		frame := []uint8{
			/* 0: */ uint8(UartSyncFrame),
			/* 1: */ 4,
			/* 2: */ uint8(PingDevicesFrame),
			/* 3: */ uint8(AllEndpoint),
			/* 4: */ uint8(LuaEndpoint),
			/* 5: */ 0, //crc D5
		}
		frame[5] = crc.D5(frame[2:5])
		return frame
	}
	return CreatePingDevicesFrame()
}

// CreateParameterSettingsReadFrameForModule creates a ParameterSettingsRead frame
// with CRC encoding appropriate for the detected module type.
func CreateParameterSettingsReadFrameForModule(moduleType ModuleType, deviceId uint8, fieldId uint8, fieldChunk uint8) []uint8 {
	if moduleType == ModuleTypeCrossfire {
		frame := []uint8{
			/* 0: */ uint8(UartSyncFrame),
			/* 1: */ 6,
			/* 2: */ uint8(ParameterSettingsReadFrame),
			/* 3: */ deviceId,
			/* 4: */ uint8(LuaEndpoint),
			/* 5: */ fieldId,
			/* 6: */ fieldChunk,
			/* 7: */ 0, //crc D5
		}
		frame[7] = crc.D5(frame[2:7])
		return frame
	}
	return CreateParameterSettingsReadFrame(deviceId, fieldId, fieldChunk)
}

// CreateParameterSettingWriteFrameUint8ForModule creates a ParameterSettingsWrite
// (uint8) frame with CRC encoding appropriate for the detected module type.
func CreateParameterSettingWriteFrameUint8ForModule(moduleType ModuleType, deviceId uint8, fieldId uint8, fieldValue uint8) []uint8 {
	if moduleType == ModuleTypeCrossfire {
		frame := []uint8{
			/* 0: */ uint8(UartSyncFrame),
			/* 1: */ 6,
			/* 2: */ uint8(ParameterSettingsWriteFrame),
			/* 3: */ deviceId,
			/* 4: */ uint8(LuaEndpoint),
			/* 5: */ fieldId,
			/* 6: */ fieldValue,
			/* 7: */ 0, //crc D5
		}
		frame[7] = crc.D5(frame[2:7])
		return frame
	}
	return CreateParameterSettingWriteFrameUint8(deviceId, fieldId, fieldValue)
}

// CreateParameterSettingWriteFrameUint16ForModule creates a ParameterSettingsWrite
// (uint16, big-endian) frame with CRC encoding appropriate for the detected module type.
func CreateParameterSettingWriteFrameUint16ForModule(moduleType ModuleType, deviceId uint8, fieldId uint8, fieldValue uint16) []uint8 {
	if moduleType == ModuleTypeCrossfire {
		frame := []uint8{
			/* 0: */ uint8(UartSyncFrame),
			/* 1: */ 7,
			/* 2: */ uint8(ParameterSettingsWriteFrame),
			/* 3: */ deviceId,
			/* 4: */ uint8(LuaEndpoint),
			/* 5: */ fieldId,
			/* 6: */ uint8((fieldValue >> 8) & 0xFF),
			/* 7: */ uint8(fieldValue & 0xFF),
			/* 8: */ 0, //crc D5
		}
		frame[8] = crc.D5(frame[2:8])
		return frame
	}
	return CreateParameterSettingWriteFrameUint16(deviceId, fieldId, fieldValue)
}

func PackChannels(channels *[16]util.CRSFValue) (result []byte) {
	var CrossfireChBits uint8 = 11
	var buf [26]byte
	var offset uint8 = 0
	var bits util.CRSFValue
	var bitsAvailable uint8 = 0

	buf[offset] = 0xEE
	offset += 1

	buf[offset] = 24 // 1(ID) + 22 + 1(CRC)
	offset += 1

	buf[offset] = 0x16
	offset += 1

	for i := 0; i < 16; i++ {
		var val = channels[i]
		shifted := val << bitsAvailable
		bits |= shifted
		bitsAvailable += CrossfireChBits
		for bitsAvailable >= 8 {
			buf[offset] = byte(bits)
			offset += 1
			bits >>= 8
			bitsAvailable -= 8
		}
	}

	buf[25] = crc.D5(buf[2:25])
	//fmt.Printf("%d\n", channels)
	//fmt.Printf("%x\n", buf)
	return buf[:]
}

func AdjustSendRate(rate int32, offset int32) time.Duration {
	duration := time.Duration((rate+offset)/10) * time.Microsecond
	if duration <= 0 {
		return MinRefreshRate
	}

	if duration > MaxRefreshRate {
		return MaxRefreshRate
	}

	return duration
}
