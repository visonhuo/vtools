# vtools
A tool collection for some personal experiment and interesting.

## vecho
vecho is a command tool which can used to set up a simple **Echo Server** or **Echo Client**.
1. For echo server, support **tcp4/tcp6/udp4/udp6** protocol;
2. For echo client, support **tcp4/tcp6/udp4/udp6/ip:udp4** protocol;
3. Support console interactive mode and fast send mode;

The original purpose of this command is used for personal experiment, details: [《TCP/UDP端口复用实验》](https://www.yuque.com/docs/share/fca5c475-de48-4ea7-9b17-59428618ca49)

### Echo server setup
```bash
# Create TCP echo server on local port 64886
vecho server --protocol=tcp --listen_ip=0.0.0.0 --listen_port=64886
# Use shorthand
vecho server -p=tcp -l=0.0.0.0 -L=64886
# More concise command 
# default value of protocol(p): tcp
# default value of listen_ip(l): 0.0.0.0
vecho server -L=64886
```

If you want to use ipv6 address, you should also set **--zone(-z)** flag.
```bash
# Create TCP6 echo server on local port 64886
vecho server -p=tcp6 -l=fe80::8c:574d:b224:3f2f -L=64886 -z=en0

# By default, it will use ipv6zero [::] address, if so, then no need to set zone flag
vecho server -p=tcp6 -L=64886
```

For UDP echo server, we can change **--protocol(-p)** flag value.  
```bash
# Create UDP echo server on local port 64886
vecho server -p=udp -L=64886

# Use ipv6 address
vecho server -p=udp6 -l=fe80::8c:574d:b224:3f2f -L=64886 -z=en0
```

### Echo client setup
By default, if vecho doesn't receive any args, it will enter console interactive mode, scan input from **system.in**, then send data to remote side. 
```bash
# Create TCP echo client, connect to <192.168.0.2:64886> from <192.168.0.1:54666>
vecho client -l=192.168.0.1 -L=54666 -r=192.168.0.2 -R=64886

# Create UDP echo client, connect to local port 9999.
# Default port value is 0, then OS will assign random port for echo client. 
vecho client -p=udp -R=9999

# Use ipv6 address
vecho client -p=udp6 -r=fe80::8c:574d:b224:3f2f -R=9999 -z=en0
```

If vecho receives any args, vecho will send these args to server side and exit quickly. (**Fast send mode**)
```bash
# Connect to local port 54666, and send 'hello' to server side
vecho client hello -p=tcp -R=54666

# Combine with bash for-loop statement, can use for testing 
for i in {1..20}; do vecho client hello -p=udp -R=54666; done
```

A special protocol value: **ip4:udp**, it will create UDP packet payload by itself and use **Raw Socket** to send data. (Need root permission)
```bash
# Connect to 9999 port from 8888 port on local machine 
sudo vecho client -p=ip4:udp -L=8888 -R=9999
```
