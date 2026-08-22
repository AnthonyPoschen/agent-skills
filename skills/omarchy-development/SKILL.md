---
name: omarchy-development
description: Develop, test, debug, or contribute to Omarchy source code in an isolated environment. Use whenever the task touches an Omarchy checkout, Omarchy shell/QML code, installation source, migrations, upstream Omarchy PRs, `omarchy dev link`, or testing Omarchy changes in a VM. Keep the host Omarchy installation untouched; create a disposable KVM/QEMU overlay, register it in Virtual Machine Manager, and clean it up through libvirt when finished. Do not use for ordinary end-user Omarchy configuration changes.
---

# Omarchy Development

Develop Omarchy in a checkout and test it in a disposable guest. The host's installed Omarchy must stay on its packaged version: do not run `omarchy dev link`, restart the host shell, or reboot the host to test source changes.

## Read source guidance first

In the Omarchy checkout, read `AGENTS.md`. Then read every linked task guide relevant to the change, especially:

- `agents/skills/shell-dev.md` for QML/shell work.
- `agents/skills/acceptance-tests.md` for graphical acceptance tests.
- `agents/skills/visual-verification.md` for visible UI changes.

Run focused source tests before using the VM. For shell changes, this normally includes the affected `test/shell.d/*-test.sh` and `qmllint` on changed QML files.

## Prepare the host once

Install the virtualization stack only with the user's authorization. Use `sudo` from an interactive terminal; when an agent cannot enter a password, use `pkexec`.

```bash
omarchy pkg add qemu-desktop qemu-full libvirt virt-manager dnsmasq edk2-ovmf swtpm
sudo systemctl enable --now libvirtd.service
sudo usermod -aG libvirt "$USER"
sudo virsh net-start default
sudo virsh net-autostart default
```

The user must log out and back in before Virtual Machine Manager can use the new `libvirt` group membership. Confirm KVM is available with `test -r /dev/kvm`.

Use the official `omacom-io/omarchy-iso` harness to make an installed base once. Keep the ISO outside `/tmp` and verify its published checksum. The harness's `--install-only` mode creates a reusable `base.qcow2`; never boot or write that base directly.

## Create a GUI-managed disposable overlay

Use a separate source worktree per change and an overlay directory per VM. Store overlays under `~/VMs/omarchy-<feature>/`, not `/tmp`.

```bash
VM_NAME=omarchy-<feature>
VM_DIR="$HOME/VMs/$VM_NAME"
BASE_DISK=/absolute/path/to/base.qcow2
BASE_VARS=/absolute/path/to/OVMF_VARS.4m.fd

mkdir -p "$VM_DIR"
qemu-img create -f qcow2 -b "$BASE_DISK" -F qcow2 "$VM_DIR/disk.qcow2"
cp "$BASE_VARS" "$VM_DIR/OVMF_VARS.4m.fd"
```

Generate a libvirt domain XML that uses:

- `type='kvm'`, UEFI OVMF code plus the copied per-overlay NVRAM file.
- The overlay as a `virtio` qcow2 disk.
- 8 GiB memory and 8 vCPUs unless the task needs more or less.
- SPICE graphics, a virtio video device, USB tablet input, and host-passthrough CPU.
- A persistent domain named exactly `$VM_NAME`.

Use QEMU user-mode networking for these Omarchy development VMs:

```xml
<interface type='user'>
  <model type='virtio'/>
</interface>
```

This is deliberate. It matches the ISO harness's proven direct-QEMU networking and avoids a guest DHCP failure observed with libvirt's `default` NAT network. Do not substitute `<interface type='network'><source network='default'/></interface>` unless the guest's network path is explicitly verified. User-mode networking gives the guest outbound Internet for cloning and package operations; it does not automatically provide host-port forwarding.

Define the VM through libvirt so it appears in Virtual Machine Manager:

```bash
pkexec virsh -c qemu:///system define /absolute/path/to/$VM_NAME.xml
pkexec virsh -c qemu:///system dominfo "$VM_NAME"
```

If redefining an existing domain, preserve its UUID from `virsh dominfo`/`virsh dumpxml`; omitting it makes libvirt treat the definition as a conflicting new VM. Confirm the interface is `user` with:

```bash
pkexec virsh -c qemu:///system domiflist "$VM_NAME"
```

Start it in Virtual Machine Manager, or with `pkexec virsh -c qemu:///system start "$VM_NAME"` and then open its console in Virtual Machine Manager. Never run the same overlay through direct QEMU and libvirt at the same time.

## Link source inside the guest only

In the guest, clone the exact feature branch and enable development mode there:

```bash
git clone --branch <branch> <fork-url> ~/omarchy
omarchy dev link ~/omarchy --no-reboot
reboot
```

Later source updates are:

```bash
cd ~/omarchy
git pull --ff-only
omarchy restart shell
```

If the guest cannot clone, first check `nmcli device status`. With the required user-mode NIC it should obtain networking automatically. Do not alter the host network before proving the guest's NIC and domain interface type.

## Manage and remove overlays safely

Use libvirt status as the source of truth:

```bash
pkexec virsh -c qemu:///system list --all
pkexec virsh -c qemu:///system shutdown "$VM_NAME"
```

Wait until `domstate` is `shut off` before changing the definition or deleting files. To delete an overlay after the user is finished with it:

1. Confirm `VM_NAME`, `VM_DIR`, and the domain are the intended disposable test VM.
2. Confirm it is shut off.
3. Remove the libvirt definition: `pkexec virsh -c qemu:///system undefine "$VM_NAME"`.
4. Verify it no longer appears in `virsh list --all`.
5. Delete only that resolved overlay directory. Never delete the reusable base disk or an unresolved/broad path.

The normal cleanup removes both the Virtual Machine Manager entry and the overlay. Preserve the overlay when the user wants to resume testing; it contains all guest-side source links and test changes.

## Report handoff

State clearly:

- the source worktree and branch tested;
- VM name and overlay path;
- whether the host stayed unmodified;
- guest test commands and results;
- whether the overlay remains available in Virtual Machine Manager or was removed.
