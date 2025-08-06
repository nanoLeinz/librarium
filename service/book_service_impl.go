package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/nanoLeinz/librarium/model"
	"github.com/nanoLeinz/librarium/repository"
	"github.com/sirupsen/logrus"
)

type BookServiceImpl struct {
	repo repository.BookRepository
	log  *logrus.Logger
}

func NewBookServiceImpl(log *logrus.Logger, repo repository.BookRepository) BookService {
	return &BookServiceImpl{
		log:  log,
		repo: repo,
	}
}

func (s *BookServiceImpl) Create(ctx context.Context, data *model.Book) (*model.Book, error)
func (s *BookServiceImpl) Update(ctx context.Context, id uuid.UUID, data map[string]any) error
func (s *BookServiceImpl) DeleteByID(ctx context.Context, id uuid.UUID) error
func (s *BookServiceImpl) GetByID(ctx context.Context, id uuid.UUID) (*model.Book, error)
func (s *BookServiceImpl) GetByTitle(ctx context.Context, name string) (*[]model.Book, error)
func (s *BookServiceImpl) GetAll(ctx context.Context) (*[]model.Book, error)
