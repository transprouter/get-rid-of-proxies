#!/usr/bin/python

from mininet.clean import cleanup
from mininet.cli import CLI
from mininet.node import Host, OVSController
from mininet.topo import Topo
from mininet.net import Mininet
from mininet.util import waitListening
from mininet.log import setLogLevel
from shutil import copytree
from tempfile import mkdtemp

class Server( Host ):
    """
    A node running some services
    """
    def __init__( self, name, **kwargs ):
        Host.__init__( self, name, **kwargs )
        self.services = []

    def service(self, cmd, port):
        """
        Starts a service with the given command and wait until something is
        listening on given port number
        """
        self.cmd(cmd + ' &')
        self.services.append( int( self.cmd('echo $!' ) ) )
        waitListening( client=self, port=port, timeout=5 )

    def terminate( self ):
        """
        Stops all started services
        """
        for pid in self.services:
            self.cmd( 'kill -15 %d' % pid )
        super( Server, self ).terminate()


class DnsServer( Server ):
    """
    A host running a DNS service

    It takes a 'hosts' array of names and ip addresses. For instance:
        hosts=[
            { 'name': 'domain.tld', 'ip': '192.168.4.2' }
        ]
    """
    def __init__( self, name, **kwargs ):
        self.temp_dir = mkdtemp( prefix='%s-' % name )
        copytree( '/etc', '%s/etc' % self.temp_dir )
        kwargs['privateDirs'] = [
            ( '/etc', '%s/etc' % self.temp_dir )
        ]
        Server.__init__( self, name, **kwargs )
        with open( '%s/etc/hosts' % self.temp_dir, 'w' ) as f:
            f.write( '127.0.0.1 localhost\n' )
            for host in kwargs['hosts']:
                f.write( '%s %s\n' % (host['ip'], host['name']) )

    def config( self, **params ):
        super( Server, self).config( **params )
        self.service( '/usr/sbin/dnsmasq', 53 )


class Services( Server ):
    """
    A host running sshd and nginx
    """
    def __init__( self, name, **kwargs ):
        Server.__init__( self, name, **kwargs )

    def config( self, **params ):
        super( Server, self).config( **params )
        self.service( '/usr/sbin/sshd', 22 )
        self.service( '/usr/sbin/nginx', 80 )


class NetworkTopo( Topo ):
    """
    Networks setup

    """

    def build( self ):
        private = self.addSwitch( 's1' )
        public  = self.addSwitch( 's2' )
        
        proxy = self.addHost( 'p1' )
        self.addLink( proxy, private, intfName2='p0-eth1', params2={ 'ip' : '192.168.35.1/24' } )
        self.addLink( proxy, public,  intfName2='p0-eth2', params2={ 'ip' : '10.0.0.1/24' } )

        dns = self.addHost( 'dns1', cls=DnsServer, ip='192.168.35.2/24',
            hosts=[
                { 'name': 'srv1.local', 'ip': '192.168.35.10' },
                { 'name': 'ws1.local',  'ip': '192.168.35.11' },
                { 'name': 'ws2.local',  'ip': '192.168.35.12' }
            ]
        )
        localServer1 = self.addHost( 'srv1', cls=Services, ip='192.168.35.10/24' )
        workstation1 = self.addHost( 'ws1', ip='192.168.35.11/24' )
        workstation2 = self.addHost( 'ws2', ip='192.168.35.12/24' )

        self.addLink( dns         , private )
        self.addLink( workstation1, private )
        self.addLink( workstation2, private )
        self.addLink( localServer1, private )

        publicServer1 = self.addHost( 'srv2', cls=Services, ip='10.0.0.10/24' )

        self.addLink( publicServer1, public )



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
