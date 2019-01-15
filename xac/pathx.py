from pathlib import Path as _Path_, PosixPath as _PosixPath_, WindowsPath  as _WindowsPath_
import os

"""
Subclassing pathlib for the sake of syntactic sugar.
Currently not working te way I want.
"""


def path_join(self, other):
    """
    / functionality for joining paths
    :param other: path-like
    :return:
    """
    return self.joinpath(other)


def path_rel(self, other):
    """
    % implements relative_to, giving the common base
    :param other:
    :return:
    """
    return self.relative_to(other)


def path_lop(self, other):
    """
    Path "lop" - remove common nodes from the left side until you get to a difference or anchor.
    Currently forces path diff to match.
    Also can do "path subtraction" by going up an integer number of directories
    :param self:
    :param other:
    :return:
    """
    self = type(self)(self) # janky copy, not sure if necessary


    if isinstance(other, int):
        if other < 0:
            raise ValueError('Cannot understand path manipulation, negative paths')
        while other > 0:
            self = self.parent()
            other -= 1
        return self

    other = type(self)(other)
    while self.name == other.name:
        self = self.parent
        other = other.parent
    # Ensure we can do a full "subtraction"
    if other != _Path_('.'):
        raise ValueError('Unable to fully lop-off (subtract) paths')
    return self


def path_name(self):
    return self.name


bolt_on_methods = {
    '__div__': path_join,
    '__mod__': path_rel,
    '__sub__': path_lop,
    '__invert__': path_name
}

Path = _Path_
PosixPath = _PosixPath_
WindowsPath = _WindowsPath_

for name, method in bolt_on_methods.items():
    setattr(Path, name, method)

class NewPath(_Path_):
    def __new__(cls, *args, **kvps):
        return super(NewPath).__new__(NewWindowsPath if os.name == 'nt' else NewPosixPath, *args, **kvps)

    def __add__(self, other):
        return self.joinpath(other)

    def __div__(self, other):
        """
        / functionality for joining paths
        :param other: path-like
        :return:
        """
        return self.joinpath(other)

    def __mod__(self, other):
        """
        % implements relative_to, giving the common base
        :param other:
        :return:
        """
        return self.relative_to(other)


class NewWindowsPath(_WindowsPath_, NewPath):
    pass

class NewPosixPath(_PosixPath_, NewPath):
    pass