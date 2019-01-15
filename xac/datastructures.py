class Borg(object):
    """
    Every instance of Borg() shares the same state. Data for consumption is stored in _data.
    You will be assimilated.
    """
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
        return self._data[key]

    def __repr__(self):
        return 'Borg({})'.format(self._state)