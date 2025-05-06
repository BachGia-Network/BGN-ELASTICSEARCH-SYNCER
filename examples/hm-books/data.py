import click
from schema import Book, Author, BookAuthor
from sqlalchemy.orm import sessionmaker

from pgsync.base import pg_engine, subtransactions
from pgsync.helper import teardown
from pgsync.utils import config_loader, get_config


@click.command()
@click.option(
    "--config",
    "-c",
    help="Schema config",
    type=click.Path(exists=True),
)
def main(config: str) -> None:
    config: str = get_config(config)
    teardown(drop_db=False, config=config)
    doc: dict = next(config_loader(config))
    database: str = doc.get("database", doc["index"])
    with pg_engine(database) as engine:
        Session = sessionmaker(bind=engine, autoflush=True)
        session = Session()

        # Create books first
        books = [
            Book(
                isbn="9785811243570",
                title="Charlie and the chocolate factory",
                description="Willy Wonka's famous chocolate factory is opening at last!",
            ),
            Book(
                isbn="9788374950978",
                title="Kafka on the Shore",
                description="Kafka on the Shore is a 2002 novel by Japanese author Haruki Murakami",
            ),
            Book(
                isbn="9781471331435",
                title="1984",
                description="1984 was George Orwell's chilling prophecy about the dystopian future"
            )
        ]

        # Create authors
        authors = [
            Author(id="Roald Dahl", name="Roald Dahl"),
            Author(id="Haruki Murakami", name="Haruki Murakami"),
            Author(id="Philip Gabriel", name="Philip Gabriel"),
            Author(id="George Orwell", name="George Orwell"),
        ]

        # Create book-author relationships
        book_authors = [
            BookAuthor(isbn="9785811243570", author_id="Roald Dahl"),
            BookAuthor(isbn="9788374950978", author_id="Haruki Murakami"),
            BookAuthor(isbn="9788374950978", author_id="Philip Gabriel"),
            BookAuthor(isbn="9781471331435", author_id="George Orwell"),
        ]

        with subtransactions(session):
            # First add books and authors
            session.add_all(books)
            session.add_all(authors)
            session.flush()  # Ensure books and authors are created before adding relationships
            
            # Then add the relationships
            session.add_all(book_authors)


if __name__ == "__main__":
    main()
