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

#define SYSCLK 72000000l
#define ADC_PRESCALAR 10l
#define ADC_SAMPLE_RATE (SYSCLK / ADC_PRESCALAR) / (32l)

#define ADC_OUT_SAMPLES_PER_INT 16
#define ADC_OVERSAMPLING 16
#define ADC_CHANS 5
#define ADC_OVERSAMPLED_LEN (ADC_OVERSAMPLING * ADC_CHANS) * ADC_OUT_SAMPLES_PER_INT
//#define ADC_OUT_LEN 2197
#define ADC_OUT_LEN 675

#define CURRENT_SENSOR_COUNT 3

#define ADC_CHANNEL_CT0 2
#define ADC_CHANNEL_CT1 3
#define ADC_CHANNEL_CT2 4
#define ADC_CHANNEL_V 1
