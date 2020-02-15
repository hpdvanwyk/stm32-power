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

#include <libopencm3/usb/usbd.h>

typedef struct {
    float dc;
    float rms;
    float real_power;
    float apparent_power;
    float powerfactor;
} current_sensor;

typedef struct {
    float dc;
    float rms;
} voltage_sensor;

void calculate_power(
    usbd_device*        usbd_dev,
    volatile uint16_t** i,
    volatile uint16_t*  v,
    current_sensor*     cs,
    voltage_sensor*     vs,
    int                 len);