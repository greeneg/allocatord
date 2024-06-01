# allocatord - Allocation Control Service

This daemon manages requests for reimaging a system in your environment.

## What does Allocator Daemon Do?

The Allocator Daemon records and allows client systems to know if they need to reimage themselves, and with what image, settings, etc.

## Does Allocator Have a Client?

Eventually, yes, there will be a client available that will use a miniature Linux distribution based on Busybox and assorted other components that does the heavy lifting of imaging the system, managing the bootloader, etc.
