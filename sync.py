
import socketserver
import json

class ProtocolHandler: # because socketserver sucks.
	def __init__(self, db):
		self.db = db

	def __call__(self, sock, client_address, server):
		with sock.makefile() as f:
			line = f.readline()

			while not f.closed:
				if line.startswith('IHAZ '):
					# strip command, space and newline
					args = line[5:-1].split(',')
					roots = ( (bytes.fromhex(oid), int(version))
							for oid, _sep, version in ( root.partition('/') for root in args ) )
					# TODO protocol reset
					# TODO do something awesome here

				if line.startswith('UCANHAZ '):
					obj = line[8:-1]
					objects = {}
					while obj != '\n':
						oid, _sep, data = obj.partition(':')
						objects[bytes.fromhex(oid)] = json.loads(data)
						obj = f.readline()
					# TODO do something intelligent here


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
	db = sqlite3.connect(args.db)

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


