/*
 * Copyright (C) 2020 Hendrik van Wyk
 *
 * This library is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Lesser General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This library is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Lesser General Public License for more details.
 *
 * You should have received a copy of the GNU Lesser General Public License
 * along with this library.  If not, see <http://www.gnu.org/licenses/>.
 */

#include <math.h>

#include "debug.h"
#include "power.h"
#include "adc_constants.h"

#include "usb.h"
#include "pb_decode.h"
#include "pb_encode.h"
#include "pb.h"
#include "power.pb.h"

#define A 0.0001

PowerMessage msg = PowerMessage_init_default;

static void message_to_usb(usbd_device* usbd_dev, PowerMessage* m) {

    int          ret;
    bool         status;
    pb_ostream_t ostream;

    uint8_t* buf = cdc_get_buf();
    if (buf == NULL) {
        debuglogf("no usb cdc buffers left\n");
        return;
    }

    ostream = pb_ostream_from_buffer(cdc_get_buf(), USB_TX_BUF_LEN);

    status = pb_encode_delimited(&ostream, PowerMessage_fields, m);
    if (!status) {
        debuglogf("failed encoding power message\n");
        return;
    }

    debuglogf("msg len %d\n", ostream.bytes_written);
    ret = cdc_write(usbd_dev, ostream.bytes_written);

    if (ret != 0) {
        debuglogf("usb write failure\n");
    }
}

void calculate_power(
    usbd_device*        usbd_dev,
    volatile uint16_t** i,
    volatile uint16_t*  v,
    current_sensor*     cs,
    voltage_sensor*     vs,
    int                 len) {
    int j, k;
    for (k = 0; k < len; k++) {
        for (j = 0; j < ADC_OUT_LEN; j++) {
            cs[k].dc += A * ((float)(i[k][j]) - cs[k].dc);
        }
    }
    for (j = 0; j < ADC_OUT_LEN; j++) {
        vs->dc += A * ((float)(v[j]) - vs->dc);
    }

    for (k = 0; k < len; k++) {
        float sum_inst_power = 0;
        for (j = 0; j < ADC_OUT_LEN; j++) {
            float v_instant = ((float)v[j] - vs->dc);
            sum_inst_power += (((float)i[k][j] - cs[k].dc)) * v_instant;
        }
        cs[k].real_power = (sum_inst_power ) / (float)ADC_OUT_LEN;
    }

    float v_rms_tmp = 0;
    for (j = 0; j < ADC_OUT_LEN; j++) {
        v_rms_tmp += pow(((float)v[j] - vs->dc), 2);
    }
    v_rms_tmp /= (float)ADC_OUT_LEN;
    vs->rms = sqrt(v_rms_tmp);

    for (k = 0; k < len; k++) {
        float i_rms = 0;
        for (j = 0; j < ADC_OUT_LEN; j++) {
            i_rms += pow(((float)i[k][j] - cs[k].dc), 2);
        }
        i_rms /= (float)ADC_OUT_LEN;
        cs[k].rms            = sqrt(i_rms);
        cs[k].apparent_power = cs[k].rms * vs->rms;
        cs[k].powerfactor    = cs[k].real_power / cs[k].apparent_power;
    }

    debuglogf("dc %lu.%02lu %lu.%02lu %lu.%02lu %lu.%02lu \n",
              (uint32_t)(cs[0].dc), (uint32_t)(cs[0].dc * 100.0) % 100,
              (uint32_t)(cs[1].dc), (uint32_t)(cs[1].dc * 100.0) % 100,
              (uint32_t)(cs[2].dc), (uint32_t)(cs[2].dc * 100.0) % 100,
              (uint32_t)(vs->dc), (uint32_t)(vs->dc * 100.0) % 100);

    debuglogf("real power     %ld %ld %ld \n",
              (int32_t)(cs[0].real_power),
              (int32_t)(cs[1].real_power),
              (int32_t)(cs[2].real_power));

    debuglogf("apparent power %lu %lu %lu \n",
              (uint32_t)(cs[0].apparent_power),
              (uint32_t)(cs[1].apparent_power),
              (uint32_t)(cs[2].apparent_power));

    debuglogf("power factor   %ld %ld %ld \n",
              (int32_t)(cs[0].powerfactor * 100),
              (int32_t)(cs[1].powerfactor * 100),
              (int32_t)(cs[2].powerfactor * 100));

    debuglogf("I_rms %lu.%03lu %lu.%03lu %lu.%03lu\r\n",
              (uint32_t)(cs[0].rms), (uint32_t)(cs[0].rms * 1000.0) % 1000,
              (uint32_t)(cs[1].rms), (uint32_t)(cs[1].rms * 1000.0) % 1000,
              (uint32_t)(cs[2].rms), (uint32_t)(cs[2].rms * 1000.0) % 1000);

    debuglogf("V_rms %lu.%03lu \n",
              (uint32_t)(vs->rms), (uint32_t)(cs->rms * 1000.0) % 1000);

    msg.VoltageRms = vs->rms;
    msg.DC         = vs->dc;
    for (j = 0; j < (int)(sizeof(msg.Voltage) / sizeof(float)); j++) {
        msg.Voltage[j] = (uint32_t)v[j];
    }
    for (k = 0; k < len; k++) {
        for (j = 0; j < (int)(sizeof(msg.Powers[k].Current) / sizeof(float)); j++) {
            msg.Powers[k].Current[j] = (uint32_t)i[k][j];
        }
        msg.Powers[k].DC            = cs[k].dc;
        msg.Powers[k].RealPower     = cs[k].real_power;
        msg.Powers[k].ApparentPower = cs[k].apparent_power;
        msg.Powers[k].PowerFactor   = cs[k].powerfactor;
        msg.Powers[k].CurrentRms    = cs[k].rms;
    }
    msg.Powers_count = len;
    message_to_usb(usbd_dev, &msg);
}
