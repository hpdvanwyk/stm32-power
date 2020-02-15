PROJECT = powerstm
BUILD_DIR = bin

LIBNAME		= opencm3_stm32f3
DEFS		+= -DSTM32F3

FP_FLAGS	?= -mfloat-abi=hard -mfpu=fpv4-sp-d16
ARCH_FLAGS	= -mthumb -mcpu=cortex-m4 $(FP_FLAGS)
SHARED_DIR = pb_usb
SHARED_DIR += pb
SHARED_DIR += nanopb

CFILES = power/main.c
CFILES += power/power.c
CFILES += power/debug.c
CFILES += power/usb.c
CFILES += pb/power.pb.c
CFILES += nanopb/pb_encode.c
CFILES += nanopb/pb_common.c
LDLIBS = -lm

CFLAGS += -DPB_FIELD_16BIT=1
ifdef DEBUG
CFLAGS += -DDEBUG
endif


# TODO - you will need to edit these two lines!
DEVICE=stm32f303cct6
OOCD_FILE = board/stm32f303cct6.cfg

# You shouldn't have to edit anything below here.
VPATH += $(SHARED_DIR)
INCLUDES += $(patsubst %,-I%, . $(SHARED_DIR))
OPENCM3_DIR=libopencm3

include $(OPENCM3_DIR)/mk/genlink-config.mk
include rules.mk
include $(OPENCM3_DIR)/mk/genlink-rules.mk

BLACKMAGIC_TARGET := 192.168.0.12:2002

blackmagicflash: powerstm.elf
	arm-none-eabi-gdb -nx --batch \
	-ex 'target extended-remote $(BLACKMAGIC_TARGET)' \
	-ex 'monitor swdp_scan' \
	-ex 'attach 1' \
	-ex 'monitor erase_mass' \
	-ex 'load' \
	-ex 'compare-sections' \
	-ex 'kill' \
	powerstm.elf

proto:
	cd pb&&make