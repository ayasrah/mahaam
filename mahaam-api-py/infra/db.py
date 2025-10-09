
from contextlib import contextmanager
from sqlalchemy import create_engine
from infra import configs
from typing import TypeVar, Type, List
from sqlalchemy import text
from dataclasses import dataclass
import datetime

T = TypeVar('T')
U = TypeVar('U')


class DB:
	# Class-level variable to store the engine, initially None
	_engine = None

	@staticmethod
	def get_engine():
		"""Lazily initialize and return the SQLAlchemy engine with connection pooling."""
		if DB._engine is None:
			try:
				DB._engine = create_engine(
					configs.data.dbUrl, 
					pool_pre_ping=True,
					pool_recycle=3600,
					echo=False
				)
			except AttributeError as e:
				raise Exception(
					"Database configuration is missing. Ensure configs.data.dbUrl is set.") from e
		return DB._engine

	@staticmethod
	def select_one(model: Type[T], query: str, params: dict, conn=None) -> T | None:
		with DB.get_engine().connect() as conn:
			result = conn.execute(text(query), params).mappings().first()
		if not result:
			return None
		return model(**result)


	@staticmethod
	def select_many(model: Type[T], query: str, params: dict, conn=None) -> list[T]:
		with DB.get_engine().connect() as conn:
			result = conn.execute(text(query), params or {})
			rows = result.mappings().all()
			return [model(**row) for row in rows]


	@staticmethod
	def insert(sql, params=None, conn=None):
		return DB.execute(sql, params, conn)

	@staticmethod
	def update(sql, params=None, conn=None):
		return DB.execute(sql, params, conn)

	@staticmethod
	def delete(sql, params=None, conn=None):
		return DB.execute(sql, params, conn)

	@staticmethod
	def execute(sql, params=None, conn=None):
		if conn is None:
			conn = DB.get_engine().connect()
			result = conn.execute(text(sql), params or {})
			conn.commit()
			return result.rowcount
		else:
			result = conn.execute(text(sql), params or {})
			return result.rowcount



	@staticmethod
	@contextmanager
	def transaction_scope():
		conn = DB.get_engine().connect()
		try:
			conn.begin()
			yield conn
			conn.commit()
		except Exception:
			conn.rollback()
			raise