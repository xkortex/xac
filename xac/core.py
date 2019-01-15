import os
from .datastructures import Borg
from .error_handling import XacPathNotFound


class XacBorg(object):
    """
        Every instance of Borg() shares the same state. Data for consumption is stored in _data.
        You will be assimilated.
        """
    _data = dict()
    _state = {'_data': dict()}

    def __new__(cls, *p, **k):
        self = object.__new__(cls, *p, **k)
        self.pop = self._state.pop
        self.__dict__ = cls._state
        return self

    def update(self, adict):
        self._data.update(adict)

    def pop(self, key, d=object()):
        return self._data.pop(key)

    def get(self, key):
        return self._data.get(key)

    def __repr__(self):
        return 'Borg({})'.format(self._state)

    @property
    def xacpath(self):
        try:
            return os.environ['XACPATH']
        except KeyError:
            raise XacPathNotFound()
