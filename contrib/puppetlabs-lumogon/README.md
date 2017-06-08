# Lumogon type/provider for Puppet

## Overview

This module contains a custom type/provider for [Lumogon](https://lumogon.com), a container scanning and reporting tool. This custom type/provider exposes containers on a system as resources within Puppet to be consumed for discovery purposes.

## Requirements

* Puppet 4 or Later

## Usage

### Using `puppet resource`

```shell
$ puppet resource lumogon

lumogon { '/nginx':
  dpkg            => {'acl' => '2.2.52-2', 'adduser' => '3.113+nmu3', 'apt' => '1.0.9.8.4', 'base-files' => '8+deb8u7', 'base-passwd' => '3.5.37', 'bash' => '4.3-11+deb8u1', 'bsdutils' => '1:2.25.2-6', 'ca-certificates' => '20141019+deb8u2', 'coreutils' => '8.23-4', 'dash' => '0.5.7-4+b1', 'debconf' => '1.5.56', 'debconf-i18n' => '1.5.56', 'debian-archive-keyring' => '2014.3', 'debianutils' => '4.4+b1', 'diffutils' => '1:3.3-1+b1', 'dmsetup' => '2:1.02.90-2.2+deb8u1', 'dpkg' => '1.17.27', 'e2fslibs' => '1.42.12-2+b1', 'e2fsprogs' => '1.42.12-2+b1', 'findutils' => '4.4.2-9+b1', 'fontconfig-config' => '2.11.0-6.3+deb8u1', 'fonts-dejavu-core' => '2.34-1', 'gcc-4.8-base' => '4.8.4-1', 'gcc-4.9-base' => '4.9.2-10', 'gettext-base' => '0.19.3-2', 'gnupg' => '1.4.18-7+deb8u3', 'gpgv' => '1.4.18-7+deb8u3', 'grep' => '2.20-4.1', 'gzip' => '1.6-4', 'hostname' => '3.15', 'inetutils-ping' => '2:1.9.2.39.3a460-3', 'init' => '1.22', 'initscripts' => '2.88dsf-59', 'insserv' => '1.14.0-5', 'iproute2' => '3.16.0-2', 'libacl1' => '2.2.52-2', 'libapt-pkg4.12' => '1.0.9.8.4', 'libasprintf0c2' => '0.19.3-2', 'libattr1' => '1:2.4.47-2', 'libaudit-common' => '1:2.4-1', 'libaudit1' => '1:2.4-1+b1', 'libblkid1' => '2.25.2-6', 'libbz2-1.0' => '1.0.6-7+b3', 'libc-bin' => '2.19-18+deb8u7', 'libc6' => '2.19-18+deb8u7', 'libcap2' => '1:2.24-8', 'libcap2-bin' => '1:2.24-8', 'libcomerr2' => '1.42.12-2+b1', 'libcryptsetup4' => '2:1.6.6-5', 'libdb5.3' => '5.3.28-9', 'libdebconfclient0' => '0.192', 'libdevmapper1.02.1' => '2:1.02.90-2.2+deb8u1', 'libexpat1' => '2.1.0-6+deb8u3', 'libfontconfig1' => '2.11.0-6.3+deb8u1', 'libfreetype6' => '2.5.2-3+deb8u1', 'libgcc1' => '1:4.9.2-10', 'libgcrypt20' => '1.6.3-2+deb8u2', 'libgd3' => '2.1.0-5+deb8u9', 'libgdbm3' => '1.8.3-13.1', 'libgeoip1' => '1.6.2-4', 'libgpg-error0' => '1.17-3', 'libjbig0' => '2.1-3.1', 'libjpeg62-turbo' => '1:1.3.1-12', 'libkmod2' => '18-3', 'liblocale-gettext-perl' => '1.05-8+b1', 'liblzma5' => '5.1.1alpha+20120614-2+b3', 'libmount1' => '2.25.2-6', 'libncurses5' => '5.9+20140913-1+b1', 'libncursesw5' => '5.9+20140913-1+b1', 'libpam-modules' => '1.1.8-3.1+deb8u2', 'libpam-modules-bin' => '1.1.8-3.1+deb8u2', 'libpam-runtime' => '1.1.8-3.1+deb8u2', 'libpam0g' => '1.1.8-3.1+deb8u2', 'libpcre3' => '2:8.35-3.3+deb8u4', 'libperl5.20' => '5.20.2-3+deb8u6', 'libpng12-0' => '1.2.50-2+deb8u3', 'libprocps3' => '2:3.3.9-9', 'libreadline6' => '6.3-8+b3', 'libselinux1' => '2.3-2', 'libsemanage-common' => '2.3-1', 'libsemanage1' => '2.3-1+b1', 'libsepol1' => '2.3-2', 'libslang2' => '2.3.0-2', 'libsmartcols1' => '2.25.2-6', 'libss2' => '1.42.12-2+b1', 'libssl1.0.0' => '1.0.1t-1+deb8u6', 'libstdc++6' => '4.9.2-10', 'libsystemd0' => '215-17+deb8u6', 'libtext-charwidth-perl' => '0.04-7+b3', 'libtext-iconv-perl' => '1.7-5+b2', 'libtext-wrapi18n-perl' => '0.06-7', 'libtiff5' => '4.0.3-12.3+deb8u2', 'libtinfo5' => '5.9+20140913-1+b1', 'libudev1' => '215-17+deb8u6', 'libusb-0.1-4' => '2:0.1.12-25', 'libustr-1.0-1' => '1.0.4-3+b2', 'libuuid1' => '2.25.2-6', 'libvpx1' => '1.3.0-3', 'libx11-6' => '2:1.6.2-3', 'libx11-data' => '2:1.6.2-3', 'libxau6' => '1:1.0.8-1', 'libxcb1' => '1.10-3+b1', 'libxdmcp6' => '1:1.1.1-1+b1', 'libxml2' => '2.9.1+dfsg1-5+deb8u4', 'libxpm4' => '1:3.5.12-0+deb8u1', 'libxslt1.1' => '1.1.28-2+deb8u2', 'login' => '1:4.2-3+deb8u3', 'lsb-base' => '4.1+Debian13+nmu1', 'mawk' => '1.3.3-17', 'mount' => '2.25.2-6', 'multiarch-support' => '2.19-18+deb8u7', 'ncurses-base' => '5.9+20140913-1', 'ncurses-bin' => '5.9+20140913-1+b1', 'netbase' => '5.3', 'nginx' => '1.11.12-1~jessie', 'nginx-module-geoip' => '1.11.12-1~jessie', 'nginx-module-image-filter' => '1.11.12-1~jessie', 'nginx-module-njs' => '1.11.12.0.1.9-1~jessie', 'nginx-module-perl' => '1.11.12-1~jessie', 'nginx-module-xslt' => '1.11.12-1~jessie', 'openssl' => '1.0.1t-1+deb8u6', 'passwd' => '1:4.2-3+deb8u3', 'perl' => '5.20.2-3+deb8u6', 'perl-base' => '5.20.2-3+deb8u6', 'perl-modules' => '5.20.2-3+deb8u6', 'procps' => '2:3.3.9-9', 'readline-common' => '6.3-8', 'sed' => '4.2.2-4+deb8u1', 'sensible-utils' => '0.0.9', 'startpar' => '0.59-3', 'systemd' => '215-17+deb8u6', 'systemd-sysv' => '215-17+deb8u6', 'sysv-rc' => '2.88dsf-59', 'sysvinit-utils' => '2.88dsf-59', 'tar' => '1.27.1-2+deb8u1', 'tzdata' => '2017a-0+deb8u1', 'ucf' => '3.0030', 'udev' => '215-17+deb8u6', 'util-linux' => '2.25.2-6', 'zlib1g' => '1:1.2.8.dfsg-2+b1'},
  hostname        => 'ea100da360fe',
  id              => '0ea7bbc3d1815f3b7e159dce1e68755a99f8aee55f85911831d80ade99c4a605',
  os              => 'linux',
  platform        => 'debian',
  platformfamily  => 'debian',
  platformversion => '8.7',
}
```
