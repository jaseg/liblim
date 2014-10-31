
import socketserver
import json

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

	def perform_sync(self, oid, data):
		if oid in self.db:
			self.db.register_handlers(self.db[oid].merge(data))
		else:
			self.db.enqueue_object(oid, data)

	def __call__(self, sock, client_address, server):
		with sock.makefile() as f:

			# Say hello
			self.send_ihaz(fobj)

			line = f.readline()

			while not f.closed:
				if line.startswith('IHAZ '):
					# strip command, space and newline
					args = line[5:-1].split(',')
					roots = { (bytes.fromhex(oid), int(version))
							for oid, _sep, version in ( root.partition('/') for root in args ) }
					# TODO protocol reset
					self.send_objects(f, roots)

				elif line.startswith('UCANHAZ '):
					obj = line[8:-1]
					objects = {}
					while obj != '\n':
						oid, _sep, data = obj.partition(':')
						self.perform_sync(bytes.fromhex(oid), json.loads(data))
						obj = f.readline()
					# TODO do something awesome here

				else:
					# pass
					assert line == '\n' # Test code!!


class DatabaseWrapper:
	def __init__(self, db):
		self.db = db
		self._queue = set()
		with self.db as db:
			db.execute('CREATE TABLE IF NOT EXISTS roots (oid text UNIQUE, description text)')
			db.execute('CREATE TABLE IF NOT EXISTS objects (oid text UNIQUE, data text)')
			results = db.execute('SELECT oid, description FROM roots').fetchall()
			self.root_objects = { oid: (self[oid], desc) for oid, desc in results }

	def add_root(self, root, description):
		""" Add the given object to the list of root objects. Returns the object in case of success, None if the object
		was already present """
		if root.oid not in self._root_objects
			self._root_objects[root.oid] = (root, description)
			self.db.execute('INSERT INTO roots VALUES (?,?)', (root.oid, description))
			return root
		return None

	def __contains__(self, oid):
		count, = self.db.execute('SELECT COUNT(*) FROM objects WHERE oid = ?', oid).fetchone()
		assert count in {0,1}
		return count > 0

	def __getitem__(self, oid):
		# Make sure to throw an error when no object is found
		return self.db.execute('SELECT data FROM objects WHERE oid = ?', oid).fetchone()[0]

	def add_object(self, obj):
		self.db.execute('INSERT INTO objects VALUES (?,?)', obj.oid, json.dumps(obj.as_dict()))

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


