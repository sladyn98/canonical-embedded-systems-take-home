#!/bin/bash

set -e

# Check for superuser privileges
if [ "$(id -u)" != "0" ]; then
  echo "This script requires superuser privileges. Please run with sudo."
  exit 1
fi

# Install required packages
apt-get update
apt-get install -y qemu-system-x86 qemu-utils busybox-static

# Create working directory
WORK_DIR="$(pwd)/qemu-image"
mkdir -p "$WORK_DIR"

# Create the initramfs
cat > "$WORK_DIR/init" << 'EOF'
#!/bin/busybox sh
/bin/busybox --install -s

echo "Mounting proc and sys"
mount -t proc none /proc
mount -t sysfs none /sys

echo "Creating device nodes"
mknod /dev/null c 1 3
mknod /dev/tty c 5 0
mknod /dev/console c 5 1
mknod /dev/ttyS0 c 4 64

echo "hello world"

echo "Starting a shell..."
setsid cttyhack /bin/sh

umount /proc
umount /sys
poweroff -f
EOF

chmod +x "$WORK_DIR/init"

# Create initramfs
pushd "$WORK_DIR"

find . | cpio -H newc -o > "$WORK_DIR/initramfs.cpio"
popd

# Download pre-built Linux kernel
KERNEL_URL="https://cdn.kernel.org/pub/linux/kernel/v5.x/linux-5.14.9.tar.xz"
KERNEL_ARCHIVE="$WORK_DIR/linux-5.14.9.tar.xz"
KERNEL_DIR="$WORK_DIR/linux-5.14.9"
wget -O "$KERNEL_ARCHIVE" "$KERNEL_URL"
tar -C "$WORK_DIR" -xf "$KERNEL_ARCHIVE"
cp "$KERNEL_DIR/arch/x86/boot/bzImage" "$WORK_DIR/vmlinuz"

# Run QEMU with the kernel and initramfs
qemu-system-x86_64 \
  -kernel "$WORK_DIR/vmlinuz" \
  -initrd "$WORK_DIR/initramfs.cpio" \
  -append "console=ttyS0" \
  -nographic \
  -m 256M

