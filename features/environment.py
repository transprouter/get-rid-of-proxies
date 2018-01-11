from network import topology

def before_all(context):
    print( 'Starting network' )
    context.mn = topology.start()
    print( 'Network started' )

def after_all(context):
    print( 'Stopping network')
    context.mn.stop()


