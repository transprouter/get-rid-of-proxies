#!/usr/bin/python

from mininet.clean import cleanup
from mininet.cli import CLI
from mininet.node import Node, Host, OVSController
from mininet.topo import Topo
from mininet.net import Mininet
from mininet.util import waitListening
from mininet.log import setLogLevel
from os import path, mkdir, removedirs
from shutil import copytree, rmtree
from tempfile import mkdtemp

class Router( Node ):
    "A Node with IP forwarding enabled."

    def config( self, **params ):
        super( Router, self).config( **params )
        # Enable forwarding on the router
        self.cmd( 'sysctl net.ipv4.ip_forward=1' )
        # Silently drop packets
        self.cmd( 'iptables -t filter -P DROP' )

    def terminate( self ):
        self.cmd( 'sysctl net.ipv4.ip_forward=0' )
        super( Router, self ).terminate()

class BaseHost ( Host ):
    """
    A base host definition with common setup for tests
    """
    def __init__( self, name, privateDirs=[], dns=None, **params ):
        self.root_dir = mkdtemp( prefix='%s-' % name )
        copytree( '%s/etc'  % path.dirname(path.abspath(__file__)), '%s/etc'  % self.root_dir )
        copytree( '%s/root' % path.dirname(path.abspath(__file__)), '%s/root' % self.root_dir )
        privateDirs = [
            ( '/etc', '%s/etc' % self.root_dir ),
            ( '/root', '%s/root' % self.root_dir )
        ]
        privateDirs += privateDirs
        self.__configure_hosts( name )
        if dns is not None:
            self.__configure_dns( dns )
        Host.__init__( self, name, privateDirs=privateDirs, **params )

    def __configure_hosts( self, name ):
        with open( '%s/etc/hostname' % self.root_dir, 'w' ) as f:
            f.write( '%s\n' % name )
        with open( '%s/etc/hosts' % self.root_dir, 'w' ) as f:
            f.write( '127.0.0.1 localhost\n' )
            f.write( '127.0.1.1 %s\n' % name )

    def __configure_dns( self, dns ):
        with open( '%s/etc/resolv.conf' % self.root_dir, 'w' ) as f:
            f.write( 'nameserver %s\n' % dns )

    def terminate( self ):
        """
        Remove temporary directories
        """
        for dir_path in self.privateDirs:
            rmtree( dir_path )
            removedirs( dir_path )
        super( Server, self ).terminate()

class Server( BaseHost ):
    """
    A node running some services
    """
    def __init__( self, name, **params ):
        BaseHost.__init__( self, name, **params )
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
        Server.__init__( self, name, **kwargs )
        with open( '%s/etc/hosts' % self.root_dir, 'a' ) as f:
            for host in kwargs['hosts']:
                f.write( '%s %s\n' % (host['ip'], host['name']) )

    def config( self, **params ):
        super( Server, self).config( **params )
        self.service( '/usr/sbin/dnsmasq -u root', 53 )


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

    def authorize( self, workstation ):
        self.cmd( 'echo %s >> /root/.ssh/authorized_keys' % workstation.get_public_key() )

class NetworkTopo( Topo ):
    """
    Networks setup

    """

    def build( self ):
        private = self.addSwitch( 's1' )
        public  = self.addSwitch( 's2' )
        
        router = self.addHost( 'r1', cls=Router, ip='192.168.35.1/24' )
        self.addLink( router, private, intfName2='r1-eth0', params2={ 'ip' : '192.168.35.1/24' } )
        self.addLink( router, public,  intfName2='r1-eth1', params2={ 'ip' : '10.0.0.1/24' } )

        dns = self.addHost( 'dns1', cls=DnsServer, ip='192.168.35.9/24',
            hosts=[
                { 'name': 'srv1.local',      'ip': '192.168.35.10' },
                { 'name': 'ws1.local',       'ip': '192.168.35.11' },
                { 'name': 'ws2.local',       'ip': '192.168.35.12' },
                { 'name': 'srv2.public.net', 'ip': '10.0.0.10' }
            ]
        )
        proxy = self.addHost( 'p1', ip='192.168.35.2' )
        localServer1 = self.addHost( 'srv1', cls=Services, ip='192.168.35.10/24',
                dns='192.168.35.9' )
        workstation1 = self.addHost( 'ws1', cls=BaseHost, ip='192.168.35.11/24',
                defaultRoute='via 192.168.35.1', dns='192.168.35.9' )
        workstation2 = self.addHost( 'ws2', cls=BaseHost, ip='192.168.35.12/24',
                defaultRoute='via 192.168.35.1', dns='192.168.35.9' )

        self.addLink( proxy       , private )
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
    try:
        net = start()
        try:
            CLI( net )
        finally:
            net.stop()
    finally:
        cleanup()
