
import time
import operator

class CRDT:
	def __init__(self):
		self._version = 0

	def merge(self, other):
		""" Merge this object with a remote object. The argument contains the remote object's data deserialized into a
		dict. """
		raise NotImplementedError()

	@property
	def oid(self):
		...

	@property
	def _timestamp(self):
		...

	@property
	def version(self):
		...

	@property
	def hash(self):
		...


class Set:
	def insert(self, obj):
		...

	def _add_object(self, data):
		raise NotImplementedError()
	
	def merge(self, other, syncprov):
		return { key: self._add_object for key in (set(other) - self.keys()) }
	
	def __getitem__(self, key):
		...
	
	def __contains__(self, key):
		...

	def keys(self):
		""" return a set of the OIDs of all objects included in this set """
		...

	def values(self):
		""" return a set of all objects in this set """
		...
	
	def hash(self):
		return reduce(operator.xor, (v.hash() for v in self.values()))


class Register:
	@property
	def value(self):
		return self._value

	@value.setter
	def set_value(self, newvalue):
		self._value = newvalue
		self._timestamp = Register._make_timestamp()
	
	@property
	def timestamp(self):
		return self._timestamp
	
	@classmethod
	def _make_timestamp(kls):
		return time.time()


class Immutable:
	
	def __getitem__(self, key):
		return self._items[key]

	def __len__(self):
		return len(self._items)
	
	def __iter__(self):
		return iter(self._items)
	
	def items():
		return self._items.items()
	
	def merge(self, other):
		# only check for equality.
		if set(self._items.items()) != set(other.items()):
			raise ValueError("Invalid data in remote immutable!")


class Post(Immutable):

	def __init__(self):
		self.content = None
		self.author = None
		self.upvotes = None
		self.downvotes = None
		...
	
	def 

