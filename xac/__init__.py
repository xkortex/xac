from .core import XacBorg
from .pathx import Path, PosixPath, WindowsPath

xacborg = XacBorg()

def xacpath():
    """
    Returns the xacpath
    :return:
    """
    return xacborg.xacpath

