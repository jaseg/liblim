
import socketserver
import json
import contextlib

class ProtocolError(Exception):
	pass

class ProtocolHandler: # because socketserver sucks.
	def __init__(self, db):
		self.db = db

	def send_objects(self, fobj, remote_roots):
		objs = { db.root_objects[remote_root]
				for remote_root, _remote_version in remote_roots
				if remote_root in db.root_objects }
		if objs:
			fobj.write('UCANHAZ ' + '\n'.join(obj.oid +':'+ json.dumps(obj.as_dict()) for obj,_desc in objs) + '\n\n')
	
	def send_ihaz(self, fobj):
		fobj.write('IHAZ ' + ','.join(oid +'/'+ obj.version for oid,(obj,_desc) in db.root_objects.items()) + '\n\n')

	# FIXME serious magick below. Clean up or document *really* well.
	def register_handlers(self, handlers):
		def insert_cb(key, handler):
			def remove_cb(oid, data):
				handler(oid, data)
				del self._queue[oid]
			self._queue[oid] = remove_cb
		for key,handler in handlers.items():
			self._queue.get(key, insert_cb)(key, handler)
	
	def enqueue_object(self, oid, data):
		def insert_cb(oid, data):
			def apply_cb(oid, handler):
				handler(oid, data)
				del self._queue[oid]
			self._queue[oid] = apply_cb
		self._queue.get(oid, insert_cb)(oid, data)



	def __call__(self, sock, client_address, server):

		baseversion = ...

		with sock.makefile() as f:

			# Say hello
			self.send_ihaz(fobj)

			line = f.readline()

			while not f.closed:
				SPACE = ' '
				if line.startswith('IHAZ'+SPACE):
					# strip command, space and newline
					args = line[5:-1].split(',')
					roots = { (bytes.fromhex(oid), int(version))
							for oid, _sep, version in ( root.partition('/') for root in args ) }
					# TODO protocol reset
					self.send_objects(f, roots)

				elif line.startswith('UCANHAZ'+SPACE):
					obj = line[8:-1]
					objects = {}
					with self.db.make_transaction() as (register, insert):
						while obj != '\n':
							oid, _sep, data = obj.partition(':')

							try:
								base = self.db[oid, baseversion]
								local = self.db[oid]
								remote = json.loads(data)
							except KeyError as e:
								raise KeyError('Object ID not found in database:', oid, e)

							merged, incoming = type(local).merge(base, local, remote)

							insert(merged)
							register(incoming)

							obj = f.readline()
				else:
					# pass
					assert line == '\n' # Test code!!


class DatabaseWrapper:
	def __init__(self, db):
		self.db = db
		with self.db as db:
			db.execute("CREATE TABLE IF NOT EXISTS roots (oid TEXT UNIQUE, description TEXT)")
			db.execute("CREATE TABLE IF NOT EXISTS objects (oid TEXT, version INTEGER, data TEXT)")
			db.execute("CREATE TABLE IF NOT EXISTS versions (version INTEGER UNIQUE, peers TEXT DEFAULT '')")
			self._head_version = db.execute("SELECT max(version) FROM versions")
			results = db.execute("SELECT oid, description FROM roots").fetchall()
			self._root_objects = { oid: (self[oid], desc) for oid, desc in results }

	@property
	def head_version(self):
		return self._head_version

	@property
	def root_objects(self):
		return self._root_objects

	def add_root(self, root, description):
		""" Add the given object to the list of root objects. Returns the object in case of success, None if the object
		was already present """
		if root.oid not in self._root_objects
			self._root_objects[root.oid] = (root, description)
			self.db.execute('INSERT INTO roots VALUES (?,?)', (root.oid, description))
			return root
		return None

	def add_version(self, version):
		self.db.execute("INSERT INTO versions (version) VALUES (?)", version)

	def __contains__(self, oid):
		count, = self.db.execute('SELECT COUNT(*) FROM objects WHERE oid = ?', oid).fetchone()
		assert count in {0,1}
		return count > 0

	def __getitem__(self, args):
		""" Get the object with the given OID, either at the given version or, if no version is given at the latest
			available version:
		
			foo[oid] or foo[oid, version]
		"""
		oid, *maybe_version = *args
		# Make sure to throw an error when no object is found
		if maybe_version:
			version, = maybe_version
			return self.db.execute('SELECT data FROM objects WHERE oid = ? AND version <= ? ORDER BY version LIMIT 1',
					oid, version).fetchone()[0]
		else:
			return self.db.execute('SELECT data FROM objects WHERE oid = ? ORDER BY version LIMIT 1', oid).fetchone()[0]

	def insert(self, obj):
		self.db.execute('INSERT INTO objects VALUES (?,?)', obj.oid, json.dumps(obj.as_dict()))

	@contextlib.contextmanager
	def make_transaction(self):
		incoming, received = {}, {}

		def register(typeinfo):
			# typeinfo is a set of tuple(oid, type)
			incoming |= typeinfo
		def insert(oid, data):
			try:
				incoming[oid](oid, data)
			except KeyError:
				raise ProtocolError('Got unexpected object with oid {}'.format(oid), data)

		try:
			yield register, insert
		
		# Sanity check
		missing = incoming.keys() - received.keys()
		if missing:
			raise ProtocolError('Received objects not matching required objects', missing)

		# Execute
		for obj in received.values():
			self.insert(obj)


if __name__ == '__main__':

	import argparse
	import sqlite3

	parser = argparse.ArgumentParser()
	parser.addArgument('db', help='The SQLite database file to operate on')
	parser.addArgument('op', choices=['listen', 'connect'], help='The operation to perform')
	parser.addArgument('ip', default=None, help='The IP address to bind or connect to')
	parser.addArgument('port', default=4242, type=int, help='The port to bind or connect to')
	args = parser.parse_args()

	print('Using database in "{}"'.format(args.db))
	db_backend = sqlite3.connect(args.db)
	db = DatabaseWrapper(db_backend)

	IP, PORT = args.host, args.port

	if args.op == 'listen':
		if IP is None:
			IP = '0.0.0.0'
		print('Listening on {}:{}'.format(IP, PORT))
		server = socketserver.TCPServer((IP, PORT), ProtocolHandler(db))
		server.serve_forever()

	elif sys.argv[1] == 'connect':
		if IP is None:
			IP = '127.0.0.1'
		print('Connecting to {}:{}'.format(IP, PORT))


