# vtools
A collection of personal experiment tools.

## vecho
A command tool which can support to set up both echo server/client by TCP/UDP.

Command Usage Example:
```bash
# 1. Set up TCP echo server
vecho server --listen_ip=0.0.0.0 --listen_port=64886
# the same as above, default protocol is TCP
vecho server -l=0.0.0.0 -L=64886

# 2. Set up UDP echo server (support set SO_REUSEADDR/SO_REUSEPORT sock option)
# default listen ip is '127.0.0.1'
vecho server --protocol=udp --listen_port=54996 --reuse_addr --reuse_port

# 3. Set up UDP echo client (console interaction)
vecho client --protocol=udp --remote_ip=192.168.0.1 --remote_port=54666

# 4. Write echo content quickly
vecho client content -r=192.168.0.1 -R=54666

# 5. Send content 20 times in a loop
for i in {1..20}; do vecho client hello -p=tcp -R=54666; done
```
