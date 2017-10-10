#!/usr/bin/python

from mininet.node import OVSController
from mininet.topo import Topo
from mininet.net import Mininet
from mininet.util import dumpNodeConnections
from mininet.log import setLogLevel

class NetworkTopo(Topo):
    "Setup networks"
    def build(self):
        private = self.addSwitch('s1')
        public  = self.addSwitch('s2')
        
        proxy        = self.addHost('proxy')
        workstation1 = self.addHost('ws1')
        workstation2 = self.addHost('ws2')
        localSsh1    = self.addHost('ssh1')
        localHttp1   = self.addHost('http1')

        publicSsh  = self.addHost('pub-ssh')
        publicHttp = self.addHost('pub-http')

        self.addLink(proxy       , private)
        self.addLink(workstation1, private)
        self.addLink(workstation2, private)
        self.addLink(localSsh1   , private)
        self.addLink(localHttp1  , private)

        self.addLink(proxy     , public)
        self.addLink(publicSsh , public)
        self.addLink(publicHttp, public)

def run():
    "Create and test a simple network"
    topo = NetworkTopo()
    net = Mininet(topo, controller=OVSController)
    net.start()
    print "Dumping host connections"
    dumpNodeConnections(net.hosts)
    print "Testing network connectivity"
    net.pingAll()
    net.stop()

if __name__ == '__main__':
    # Tell mininet to print useful information
    setLogLevel('info')
    run()
