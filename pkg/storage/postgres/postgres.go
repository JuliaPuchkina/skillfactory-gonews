package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

// хранилище данных
type Store struct {
	db *pgxpool.Pool
}

// конструктор объекта хранилища
func New(constr string) (*Store, error) {
	db, err := pgxpool.Connect(context.Background(), constr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

// Posts выводит все существующие публикации
func (s *Store) Posts() ([]storage.Post, error) {
	rows, err := s.db.Query(context.Background(), `
	SELECT
		id,
		author_id,
		title,
		content, 
		created_at
	FROM posts
	ORDER BY id;
`,
	)
	if err != nil {
		return nil, err
	}

	var posts []storage.Post
	// итерирование по результату выполнения запроса
	// и сканирование каждой строки в переменную
	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.AuthorID,
			&p.Title,
			&p.Content,
			&p.CreatedAt,
		)
		if err != nil {
			return nil, err
		}
		// добавление переменной в массив результатов
		posts = append(posts, p)

	}
	// ВАЖНО не забыть проверить rows.Err()
	return posts, rows.Err()
}

// AddPost создает новую публикацию
func (s *Store) AddPost(p storage.Post) error {
	err := s.db.QueryRow(context.Background(), `
		INSERT INTO tasks (title, content)
		VALUES ($1, $2);
		`,
		p.Title,
		p.Content,
	).Scan()
	return err
}

// UpdatePost обновляет публикацию
func (s *Store) UpdatePost(p storage.Post) error {
	err := s.db.QueryRow(context.Background(), `
	UPDATE posts
	SET content = '$2'
	WHERE id = $1
	`,
		p.ID,
		p.Content,
	).Scan()

	return err
}

// DeletePost удаляет публикацию
func (s *Store) DeletePost(p storage.Post) error {
	err := s.db.QueryRow(context.Background(), `
	DELETE FROM tasks
	WHERE id = $1;
	`,
		p.ID,
	).Scan()
	return err
}
