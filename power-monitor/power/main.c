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

#include <libopencm3/stm32/rcc.h>
#include <libopencm3/stm32/gpio.h>
#include <libopencm3/stm32/usart.h>
#include <libopencm3/stm32/adc.h>
#include <libopencm3/stm32/dma.h>
#include <libopencm3/cm3/nvic.h>
#include <libopencm3/usb/usbd.h>
#include <libopencm3/usb/cdc.h>
#include <libopencm3/stm32/memorymap.h>

#include <wctype.h>
#include <ctype.h>
#include <time.h>
#include <string.h>

#include "debug.h"
#include "power.h"
#include "usb.h"

#include "adc_constants.h"

volatile uint16_t adc_raw[ADC_OVERSAMPLED_LEN * 2];

volatile uint16_t adc1_out[ADC_OUT_LEN * 2];
volatile uint16_t adc2_out[ADC_OUT_LEN * 2];
volatile uint16_t adc3_out[ADC_OUT_LEN * 2];
volatile uint16_t adc4_out[ADC_OUT_LEN * 2];
volatile uint16_t vrefdata_print;

int           out_index      = 0;
volatile int  missed_samples = 0;
volatile bool half           = false;
volatile bool full           = false;

uint32_t Vrefcal;

current_sensor currents[CURRENT_SENSOR_COUNT];
voltage_sensor voltage;

static inline void delay(int32_t n) {
    for (int32_t i = 0; i < n; i++) {
        __asm__("nop");
    }
}

static void clock_setup(void) {
    rcc_clock_setup_pll(&rcc_hse8mhz_configs[RCC_CLOCK_HSE8_72MHZ]);
    rcc_periph_clock_enable(RCC_GPIOA);
    rcc_periph_clock_enable(RCC_GPIOC);
}

static void adc_setup(void) {
    rcc_periph_clock_enable(RCC_ADC12);

    gpio_mode_setup(GPIOA, GPIO_MODE_ANALOG, GPIO_PUPD_NONE, GPIO0);
    gpio_mode_setup(GPIOA, GPIO_MODE_ANALOG, GPIO_PUPD_NONE, GPIO1);
    gpio_mode_setup(GPIOA, GPIO_MODE_ANALOG, GPIO_PUPD_NONE, GPIO2);
    gpio_mode_setup(GPIOA, GPIO_MODE_ANALOG, GPIO_PUPD_NONE, GPIO3);

    adc_power_off(ADC1);
    adc_enable_vrefint();

    rcc_adc_prescale(RCC_CFGR2_ADCxPRES_PLL_CLK_DIV_10, RCC_CFGR2_ADCxPRES_PLL_CLK_DIV_8);
    adc_set_continuous_conversion_mode(ADC1);

    adc_disable_external_trigger_regular(ADC1);
    adc_set_right_aligned(ADC1);
    adc_set_sample_time_on_all_channels(ADC1, ADC_SMPR_SMP_19DOT5CYC);

    uint8_t channels[] = {ADC_CHANNEL_CT0,
                          ADC_CHANNEL_CT1,
                          ADC_CHANNEL_CT2,
                          ADC_CHANNEL_V,
                          ADC_CHANNEL_VREF};
    adc_set_regular_sequence(ADC1, 5, channels);
    adc_set_resolution(ADC1, ADC_CFGR1_RES_12_BIT);
    adc_enable_dma_circular_mode(ADC1);

    adc_calibrate(ADC1);
    delay(10); // Does not work without this.
    adc_power_on(ADC1);

    adc_enable_dma(ADC1);
    adc_start_conversion_regular(ADC1);
}

static void dma_setup(void) {
    rcc_periph_clock_enable(RCC_DMA1);
    nvic_enable_irq(NVIC_DMA1_CHANNEL1_IRQ);

    dma_channel_reset(DMA1, DMA_CHANNEL1);

    dma_enable_circular_mode(DMA1, DMA_CHANNEL1);
    dma_enable_memory_increment_mode(DMA1, DMA_CHANNEL1);
    dma_set_peripheral_size(DMA1, DMA_CHANNEL1, DMA_CCR_PSIZE_16BIT);
    dma_set_memory_size(DMA1, DMA_CHANNEL1, DMA_CCR_MSIZE_16BIT);
    dma_set_read_from_peripheral(DMA1, DMA_CHANNEL1);
    dma_enable_transfer_complete_interrupt(DMA1, DMA_CHANNEL1);
    dma_enable_half_transfer_interrupt(DMA1, DMA_CHANNEL1);

    dma_set_peripheral_address(DMA1, DMA_CHANNEL1, (uint32_t)&ADC_DR(ADC1));
    dma_set_memory_address(DMA1, DMA_CHANNEL1, (uint32_t)&adc_raw);
    dma_set_number_of_data(DMA1, DMA_CHANNEL1, ADC_OVERSAMPLED_LEN * 2);

    dma_enable_channel(DMA1, DMA_CHANNEL1);
}

static void usart_setup(void) {
    rcc_periph_clock_enable(RCC_USART1);

    gpio_mode_setup(GPIOA, GPIO_MODE_AF, GPIO_PUPD_NONE, GPIO9 | GPIO10);
    gpio_set_af(GPIOA, GPIO_AF7, GPIO9 | GPIO10);

    usart_set_baudrate(USART_CONSOLE, 115200);
    usart_set_databits(USART_CONSOLE, 8);
    usart_set_stopbits(USART_CONSOLE, USART_STOPBITS_1);
    usart_set_mode(USART_CONSOLE, USART_MODE_TX_RX);
    usart_set_parity(USART_CONSOLE, USART_PARITY_NONE);
    usart_set_flow_control(USART_CONSOLE, USART_FLOWCONTROL_NONE);

    usart_enable(USART1);
}

inline static void calculate_oversampling(volatile uint16_t* oversampled, int* idx) {
    int      j;
    uint32_t sum0         = 0;
    uint32_t sum1         = 0;
    uint32_t sum2         = 0;
    uint32_t sum3         = 0;
    uint32_t sum_vrefdata = 0;
    for (j = 0; j < ADC_OVERSAMPLING * ADC_CHANS; j += ADC_CHANS) {
        sum0 += oversampled[j];
        sum1 += oversampled[j + 1];
        sum2 += oversampled[j + 2];
        sum3 += oversampled[j + 3];
        sum_vrefdata += oversampled[j + 4];
    }
    uint16_t vrefdata = sum_vrefdata >> 2;
    vrefdata_print    = vrefdata;
    adc1_out[*idx]    = ((sum0 >> 2) * vrefdata) / Vrefcal;
    adc2_out[*idx]    = ((sum1 >> 2) * vrefdata) / Vrefcal;
    adc3_out[*idx]    = ((sum2 >> 2) * vrefdata) / Vrefcal;
    adc4_out[*idx]    = ((sum3 >> 2) * vrefdata) / Vrefcal;
    (*idx)++;
}

void dma1_channel1_isr(void) {
    volatile uint16_t* current_res;
    if (dma_get_interrupt_flag(DMA1, DMA_CHANNEL1, DMA_HTIF)) {
        dma_clear_interrupt_flags(DMA1, DMA_CHANNEL1, DMA_HTIF);
        gpio_set(GPIOC, GPIO13);
        current_res = &adc_raw[0];
    } else {
        dma_clear_interrupt_flags(DMA1, DMA_CHANNEL1, DMA_TCIF);
        gpio_clear(GPIOC, GPIO13);
        current_res = &adc_raw[ADC_OVERSAMPLED_LEN];
    }
    int i;
    for (i = 0; i < ADC_OVERSAMPLED_LEN; i += ADC_OVERSAMPLING * ADC_CHANS) {
        calculate_oversampling(&current_res[i], &out_index);
        if (out_index == ADC_OUT_LEN) {
            if (half) {
                missed_samples++;
            }
            half = true;
        }
        if (out_index == ADC_OUT_LEN * 2) {
            if (full) {
                missed_samples++;
            }
            full = true;
        }
        out_index %= (ADC_OUT_LEN * 2);
    }
}

int main(void) {
    int          i = 0;
    usbd_device* usbd_dev;
    clock_setup();
    usart_setup();
    gpio_mode_setup(GPIOC, GPIO_MODE_OUTPUT, GPIO_PUPD_NONE, GPIO13);
    dma_setup();
    adc_setup();
    usbd_dev = usblib_init();
    Vrefcal  = ST_VREFINT_CAL << 2;
    debuglogf("vrefcal %lu\n", Vrefcal);

    for (i = 0; i < CURRENT_SENSOR_COUNT; i++) {
        currents[i].dc = 8170;
    }
    voltage.dc = 8134;

    debuglogf("sample rate %lu\n", ADC_SAMPLE_RATE / ADC_OVERSAMPLING / ADC_CHANS);
    while (1) {
        if (half || full) {
            int idx;
            if (half) {
                half = false;
                idx  = 0;
            } else {
                full = false;
                idx  = ADC_OUT_LEN;
            }
            volatile uint16_t* current_measurements[] = {
                &adc1_out[idx],
                &adc2_out[idx],
                &adc3_out[idx],
            };

            calculate_power(usbd_dev, current_measurements, &adc4_out[idx], currents, &voltage, CURRENT_SENSOR_COUNT);
            debuglogf("missed sample interrupts %d \n", missed_samples);
            debuglogf("vrefdata %u\n", vrefdata_print);
        }
        usbd_poll(usbd_dev);
    }
    return 0;
}