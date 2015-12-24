# goTcpProxy
a go tcp proxy with multi listener ports via xml config file

keep the settings.xml file in the same directory as the executable
<?xml version="1.0" encoding="UTF-8"?>
<proxy>
    <proxyserver name="Test proxy" buffersize="16384">
        <source port="21000" quedconnections="20" receivebuffersize="65535"
                sendbuffersize="65535" bindaddress="0.0.0.0"/>
        <destination port="9123" ipaddress="192.168.1.70"
                     receivebuffersize="65535" sendbuffersize="65535"/>
    </proxyserver>
</proxy>

you can supply as many nodes as you want listening tcp sockets

To run it 

./tcpproxy -debug true 
or 
./tcpproxy -debug false  [you can also omit the debug flag] 

The application should look like this when it starts

2015/12/24 15:19:35 Proxy Server Starting

2015/12/24 15:19:35 XML data loaded ok

2015/12/24 15:19:35 Name : Test proxy

2015/12/24 15:19:35 Buffersize : 16384

2015/12/24 15:19:35 	Source

2015/12/24 15:19:35 		Port : 21000

2015/12/24 15:19:35 		Bindaddress : 0.0.0.0

2015/12/24 15:19:35 		Quedconnections : 20

2015/12/24 15:19:35 		Receivebuffersize : 65535

2015/12/24 15:19:35 		Sendbuffersize : 

2015/12/24 15:19:35 	Destination

2015/12/24 15:19:35 		Port : 9123


2015/12/24 15:19:35 		Receivebuffersize : 65535
2015/12/24 15:19:35 		Sendbuffersize : 

2015/12/24 15:19:35 		Ipaddress : 192.168.1.70

2015/12/24 15:19:35 Starting Listening Server Test proxy

2015/12/24 15:19:35 Opening Listener for proxy : Test proxy

2015/12/24 15:19:35 Service : 0.0.0.0:21000

2015/12/24 15:19:35 Binding listener IP[0.0.0.0] and Port[21000]

Accepted connection : 127.0.0.1:48270

2015/12/24 15:19:41 Bytes read write 1:1 direction source->destination

2015/12/24 15:19:41 Bytes read write 1:1 direction destination->source

2015/12/24 15:19:42 Bytes read write 1:1 direction source->destination

2015/12/24 15:19:42 Bytes read write 1:1 direction destination->source

2015/12/24 15:19:50 EOF

2015/12/24 15:19:50 EOF

mail : Michael@soffex.co.za
Web : www.soffex.co.za
