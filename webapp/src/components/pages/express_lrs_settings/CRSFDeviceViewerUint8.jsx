// SPDX-FileCopyrightText: © 2023 OneEyeFPV oneeyefpv@gmail.com
// SPDX-License-Identifier: GPL-3.0-or-later
// SPDX-License-Identifier: FS-0.9-or-later

import React, {useCallback, useState} from "react";
import {CRSFDeviceFieldData, CRSFDeviceFieldUint8} from "../../../pbwrap";
import Box from "@mui/material/Box";
import {FormControl, InputLabel, MenuItem, Select} from "@mui/material";
import {setCRSFDeviceField} from "../../misc/server";
import {showError} from "../../misc/notifications";

export function CRSFDeviceViewerUint8({device, field, setReload}) {
    const [fieldValue, setFieldValue] = useState(field.getValue());

    const handleFieldChange = useCallback(async function (field, value) {
        try {
            let wrapper = new CRSFDeviceFieldData();
            let uint8Field = new CRSFDeviceFieldUint8();
            uint8Field.setId(field.getId());
            uint8Field.setValue(value);
            uint8Field.setMin(field.getMin());
            uint8Field.setMax(field.getMax());
            wrapper.setUint8(uint8Field);
            await setCRSFDeviceField(device, wrapper);
            setReload(r => ++r);
        } catch (ex) {
            showError(ex.message);
        }
    }, []);

    const onValueChange = useCallback(function (field, value) {
        setFieldValue(value);
        setTimeout(() => handleFieldChange(field, value), 0);
    }, [fieldValue]);

    const options = [];
    for (let i = field.getMin(); i <= field.getMax(); i++) {
        options.push(i);
    }

    return <Box
        key={`${field.getId()}-box`}
        style={{
            display: "flex", alignItems: "center", alignContent: "flex-end", justifyContent: "flex-start",
            paddingTop: 5,
        }}>
        <FormControl style={{width: "100%"}}>
            <Box style={{
                width: '100%',
                marginBottom: '20px',
                marginTop: '-10px'
            }}>
                <InputLabel style={{left: -15}}
                            id={`field-${field.getId()}-label`}>{field.getName()}</InputLabel>
                <Select
                    variant="standard"
                    style={{width: "100%"}}
                    labelId={`field-${field.getId()}-label`}
                    value={fieldValue}
                    onChange={(event) => onValueChange(field, event.target.value)}
                    id={`field-${field.getId()}`}
                >
                    {options.map(v => <MenuItem key={v} value={v}>{v}</MenuItem>)}
                </Select>
            </Box>
        </FormControl>
    </Box>;
}
