#include <stdio.h>
#include <stdlib.h>

#include <simavr/sim_avr.h>
#include <simavr/sim_elf.h>
#include <simavr/avr_uart.h>
#include <simavr/sim_irq.h>
#include <simavr/avr_ioport.h>

static void usage() {
    printf("usage: runner <firmware.elf>\n");

    exit(1);
}

static int load_firmware(const char *filename, elf_firmware_t *firmware) {
    firmware->flashsize = 0;
    return elf_read_firmware(filename, firmware) != 0 || firmware->flashsize == 0 ? -1 : 0;
}

static void uart_char_out_hook(avr_irq_t *irq, uint32_t value, void *param) {
    putchar(value);
}

static void gpio_write_hook(avr_irq_t *irq, uint32_t value, void *param) {
    printf("gpio write: %d\n", value);
}

static void setup_uart(avr_t *avr) {
    uint32_t flags;

    avr_ioctl(avr, AVR_IOCTL_UART_GET_FLAGS('0'), &flags);
    flags &= ~AVR_UART_FLAG_STDIO;
    avr_ioctl(avr, AVR_IOCTL_UART_SET_FLAGS('0'), &flags);

    const char* irq_names[1] = {"UART_OUT"};

    avr_irq_t *irq = avr_alloc_irq(&avr->irq_pool, 0, 1, irq_names);
    avr_irq_register_notify(irq, &uart_char_out_hook, NULL);
    avr_connect_irq(avr_io_getirq(avr, AVR_IOCTL_UART_GETIRQ('0'), UART_IRQ_OUTPUT), irq);
}

static void setup_gpio(avr_t *avr) {
    const char* irq_names[1] = {"GPIO_WRITE"};

    avr_irq_t *irq = avr_alloc_irq(&avr->irq_pool, 0, 1, irq_names);
    avr_irq_register_notify(irq, &gpio_write_hook, NULL);
    avr_connect_irq(avr_io_getirq(avr, AVR_IOCTL_IOPORT_GETIRQ('D'), IOPORT_IRQ_PIN_ALL), irq);
}

int main(int argc, char** argv) {
    if (argc != 2) {
        usage();
    }

    elf_firmware_t firmware;
    const char* firmware_file = argv[1];

    if (load_firmware(firmware_file, &firmware) != 0) {
        fprintf(stderr, "failed to read elf file '%s' \n", firmware_file);
        exit(1);
    }

    avr_t *avr = avr_make_mcu_by_name("atmega328");
    if (!avr) {
        fprintf(stderr, "failed to init avr\n");
        exit(1);
    }

    avr_init(avr);
    avr->frequency = 16000000UL;

    setup_uart(avr);
    setup_gpio(avr);

    avr_load_firmware(avr, &firmware);

    printf("starting simulator...\n\n");

    while (1) {
        switch (avr_run(avr)) {
            case cpu_Done:
                printf("\ncpu stopped, simulation done\n");
                goto done;

            case cpu_Crashed:
                fprintf(stderr, "\ncpu crashed\n");
                goto done;
        }
    }

    done: return 0;
}
