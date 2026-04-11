// SPDX-FileCopyrightText: © 2023 OneEyeFPV oneeyefpv@gmail.com
// SPDX-License-Identifier: GPL-3.0-or-later
// SPDX-License-Identifier: FS-0.9-or-later

package settings

import (
	"bytes"
	"fmt"
	"github.com/kaack/elrs-joystick-control/pkg/crossfire/telemetry"
	"github.com/kaack/elrs-joystick-control/pkg/proto/generated/pb"
	"strings"
)

// VTXField parses a CRSF VTX-type parameter (field type 14).
// Binary layout: name\0 | band(1) | channel(1) | power(1) | pitmode(1) | freq_lsb(1) | freq_msb(1) | units\0
//
// It is surfaced to the webapp as a TEXT_SELECT so the existing read/write
// machinery works without any protocol, proto, or JS changes.
type VTXField struct {
	id       uint32
	parentId uint32
	data     []uint8
	nameEnd  int
}

func NewVTXField(id uint32, parentId uint32, data []uint8) FieldType {
	nameEnd := bytes.IndexByte(data, 0)
	if nameEnd < 0 {
		nameEnd = len(data)
	}
	return &VTXField{id: id, parentId: parentId, data: data, nameEnd: nameEnd}
}

func (f *VTXField) Name() string {
	return string(f.data[0:f.nameEnd])
}

func (f *VTXField) Type() telemetry.CRSFFieldType {
	return telemetry.CrsfVtx
}

func (f *VTXField) Id() uint32 {
	return f.id
}

func (f *VTXField) ParentId() uint32 {
	return f.parentId
}

// channel returns the 0-based channel index sent by the device.
func (f *VTXField) channel() uint32 {
	if f.nameEnd+2 >= len(f.data) {
		return 0
	}
	return uint32(f.data[f.nameEnd+2])
}

func (f *VTXField) Proto() *pb.CRSFDeviceFieldData {
	ch := f.channel()

	return &pb.CRSFDeviceFieldData{
		Data: &pb.CRSFDeviceFieldData_TextSelect{
			TextSelect: &pb.CRSFDeviceFieldTextSelect{
				Name:     strings.ToValidUTF8(f.Name(), ""),
				Type:     pb.CRSFDeviceFieldType_TEXT_SELECT,
				Id:       f.Id(),
				ParentId: f.ParentId(),
				Options:  []string{"1", "2", "3", "4", "5", "6", "7", "8"},
				Value:    ch,
				Min:      0,
				Max:      7,
				Default:  0,
				Units:    "",
			},
		},
	}
}

func (f *VTXField) String() string {
	return fmt.Sprintf("(vtx) name: %s, type: %s, id: %v, pid: %v, channel: %v",
		f.Name(),
		f.Type(),
		f.Id(),
		f.ParentId(),
		f.channel()+1,
	)
}
