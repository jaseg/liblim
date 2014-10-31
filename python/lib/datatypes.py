
import time
import operator

class CRDT:
	def __init__(self, oid):
		self._oid = oid
		self._version = 0
		# debugging-only convenience thing
		self._mtime = time.time()

	@classmethod
	def merge(kls, base, local, remote, fetch_cb):
		""" Override this. Merge two objects relative to their most recent common base and return the resulting object.

			remote contains the plain deserialized transport representation in python primitives, that is e.g. in case
			of JSON transport encoding the output of json.loads .
		
			This method should return a tuple(new_object, {dict of oid: item_crdt})
		"""
		raise NotImplementedError()
	
	@classmethod
	def serialize(kls, base, obj):
		""" Override this. Generate this object's transport representation relative to the given local base. This should
			be something json.dumps comprehends. """
		raise NotImplementedError()

	def _modified(self):
		self._mtime = time.time()

	@property
	def mtime(self):
		""" Modification timestamp of this object. This property is only local and not syncronized. """
		return self._mtime

	@property
	def oid(self):
		return self._oid

	@property
	def version(self):
		return self._version



def GrowSet(item_crdt):
	class GrowSet(set, CRDT):
		@classmethod
		def merge(kls, base, local, remote):
			# remote should contain the deserialized transport representation, which is a list of object ids.
			assert type(remote) is list
			return kls(local | set(remote)), { oid: item_crdt for oid in (set(remote) - local.keys()) }

		... # FIXME iterators


class Register:
	@property
	def value(self):
		return self._value

	@value.setter
	def set_value(self, newvalue):
		self._value = newvalue
		self._timestamp = Register._make_timestamp()
		self._modified() # debug thing
	
	@property
	def timestamp(self):
		return self._timestamp
	
	@classmethod
	def _make_timestamp(kls):
		return time.time()

	@classmethod
	def merge(kls, base, local, remote):
		return kls(oid, local.value if local.timestamp > remote.timestamp else remote.value)


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
	
