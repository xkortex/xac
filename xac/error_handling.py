class XacException(Exception):
    """
    Generic exception
    """
    def __init__(self, original_exc, msg=None):
        msg = 'An exception occurred in the Xac service' if msg is None else msg
        super(XacException, self).__init__(
            msg=msg + ': {}'.format(original_exc))
        self.original_exc = original_exc


class XacPathNotFound(Exception):
    """
    Could not obtain the XACPATH from the environment.
    """
    def __init__(self, msg=None):
        msg = 'XACPATH environment variable not set' if msg is None else msg
        super(XacPathNotFound, self).__init__()
