#!/usr/bin/python

from mininet.clean import cleanup
from mininet.cli import CLI
from mininet.node import Host, OVSController
from mininet.topo import Topo
from mininet.net import Mininet
from mininet.util import waitListening
from mininet.log import setLogLevel

class Server( Host ):
    "A host running sshd and nginx"

    def __init__( self, name, **kwargs ):
        Host.__init__( self, name, **kwargs )
        self.services = []

    def config( self, **params ):
        super( Server, self).config( **params )
        self.service( '/usr/sbin/sshd', 22 )
        self.service( '/usr/sbin/nginx', 80 )

    def service(self, cmd, port):
        self.cmd(cmd + ' &')
        self.services.append( int( self.cmd('echo $!' ) ) )
        waitListening( client=self, port=port, timeout=5 )

    def terminate( self ):
        for pid in self.services:
            self.cmd( 'kill -15 %d' % pid )
        super( Server, self ).terminate()

class NetworkTopo( Topo ):
    "Setup networks"

    def build( self ):
        private = self.addSwitch( 's1' )
        public  = self.addSwitch( 's2' )
        
        proxy        = self.addHost( 'p1' )
        workstation1 = self.addHost( 'ws1' )
        workstation2 = self.addHost( 'ws2' )
        localServer1 = self.addHost( 'srv1', cls=Server )

        publicServer1 = self.addHost( 'srv2', cls=Server )

        self.addLink(proxy       , private)
        self.addLink(workstation1, private)
        self.addLink(workstation2, private)
        self.addLink(localServer1, private)

        self.addLink(proxy        , public)
        self.addLink(publicServer1, public)

def start():
    cleanup()
    topo = NetworkTopo()
    net = Mininet(topo, controller=OVSController)
    net.start()
    return net

if __name__ == '__main__':
    "Run network and drop user to cli"
    setLogLevel( 'info' )
    net = start()
    CLI( net )
    net.stop()
    cleanup()
