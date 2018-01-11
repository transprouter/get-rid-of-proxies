from network import topology

def before_all(context):
    context.mn = topology.start()

def after_all(context):
    context.mn.stop()


