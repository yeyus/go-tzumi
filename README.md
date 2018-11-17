Tzumi MagicTV
=============

Disable auto-shutdown
---------------------

modify `/app/shutdown.sh` to do nothing and backup your original file.

Configure ethernet access
-------------------------

Interface `eth0` in Tzumi side is automatically set into `br-lan` bridge which has ip `192.168.1.1`. 

On the computer side:

```console
# ip addr add 192.168.1.100/16 dev enp0s25
# ip link set enp0s25 up
```