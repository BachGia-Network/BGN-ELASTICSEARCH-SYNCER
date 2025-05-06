import click
import sqlalchemy as sa
from sqlalchemy.orm import DeclarativeBase, Mapped, mapped_column, relationship

from pgsync.base import pg_engine
from pgsync.utils import config_loader, get_config


class Base(DeclarativeBase):
    pass


class Book(Base):
    __tablename__ = "books"
    __table_args__ = (sa.PrimaryKeyConstraint("isbn"),)
    
    isbn: Mapped[str] = mapped_column(sa.String)
    title: Mapped[str] = mapped_column(sa.String)
    description: Mapped[str] = mapped_column(sa.String, nullable=True)
    
    authors = relationship("Author", secondary="book_authors")


class Author(Base):
    __tablename__ = "authors"
    __table_args__ = (sa.PrimaryKeyConstraint("id"),)
    
    id: Mapped[str] = mapped_column(sa.String)
    name: Mapped[str] = mapped_column(sa.String)
    
    books = relationship("Book", secondary="book_authors")


class BookAuthor(Base):
    __tablename__ = "book_authors"
    __table_args__ = (sa.PrimaryKeyConstraint("isbn", "author_id"),)
    
    isbn: Mapped[str] = mapped_column(sa.String, sa.ForeignKey("books.isbn"))
    author_id: Mapped[str] = mapped_column(sa.String, sa.ForeignKey("authors.id"))


def setup(config: str) -> None:
    for doc in config_loader(config):
        database: str = doc.get("database", doc["index"])
        with pg_engine(database) as engine:
            Base.metadata.drop_all(engine)
            Base.metadata.create_all(engine)


@click.command()
@click.option(
    "--config",
    "-c",
    help="Schema config",
    type=click.Path(exists=True),
)
def main(config: str) -> None:
    config: str = get_config(config)
    setup(config)


if __name__ == "__main__":
    main()
