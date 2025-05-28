package storage

import (
	"crypto/rand"
	"errors"
	"math/big"
	"sync"
)

type Quote struct {
	Author string `json:"author"`
	Quote  string `json:"quote"`
}

type QuoteWithID struct {
	ID int64 `json:"id"`
	Quote
}

type quoteBook struct {
	idKeys        map[int64]*QuoteWithID
	structKeys    map[Quote]struct{}
	indexRandList map[int64]int
	randomList    []int64
	byAuthor      map[string][]int64
	IdCount       int64
	m             sync.RWMutex
}

var (
	ErrQuoteNotFound = errors.New("quote not found")
	ErrQuoteExists   = errors.New("quote of such author already exists")
)

func NewQBook() *quoteBook {
	return &quoteBook{
		idKeys:        make(map[int64]*QuoteWithID),
		structKeys:    make(map[Quote]struct{}),
		indexRandList: make(map[int64]int),
		byAuthor:      make(map[string][]int64),
	}
}

func (qb *quoteBook) Create(quote *Quote) (int64, error) {
	err := qb.findForCreate(quote)
	if err != nil {
		return 0, err
	}
	qb.m.Lock()
	defer qb.m.Unlock()
	qb.IdCount += 1
	qb.idKeys[qb.IdCount] = &QuoteWithID{
		ID:    qb.IdCount,
		Quote: *quote,
	}
	qb.structKeys[*quote] = struct{}{}
	qb.byAuthor[quote.Author] = append(qb.byAuthor[quote.Author], qb.IdCount)
	qb.randomList = append(qb.randomList, qb.IdCount)
	qb.indexRandList[qb.IdCount] = len(qb.randomList) - 1

	return qb.IdCount, nil
}

func (qb *quoteBook) findForCreate(quote *Quote) error {
	qb.m.RLock()
	defer qb.m.RUnlock()
	_, ok := qb.structKeys[*quote]
	if ok {
		return ErrQuoteExists
	}
	return nil
}

func (qb *quoteBook) GetRandom() (*QuoteWithID, error) {
	var err error
	var randomIndex *big.Int

	qb.m.RLock()
	defer qb.m.RUnlock()
	listLength := len(qb.randomList)
	if listLength == 0 {
		return nil, ErrQuoteNotFound
	}

	keyAmount := big.NewInt(int64(listLength))
	for i := 0; i < 3; i++ {
		randomIndex, err = rand.Int(rand.Reader, keyAmount)
		if err == nil {
			break
		}
	}
	if err != nil {
		return nil, err
	}

	id := qb.randomList[randomIndex.Int64()]
	return qb.idKeys[id], nil
}

func (qb *quoteBook) List(author *string) []*QuoteWithID {
	qb.m.RLock()
	defer qb.m.RUnlock()
	quotes := []*QuoteWithID{}
	if author == nil {
		for _, quote := range qb.idKeys {
			quotes = append(quotes, quote)
		}
	} else {
		ids := qb.byAuthor[*author]
		for _, id := range ids {
			quotes = append(quotes, qb.idKeys[id])
		}
	}

	return quotes
}

func (qb *quoteBook) Delete(id int64) error {
	structKey, err := qb.findForDelete(id)
	if err != nil {
		return err
	}
	qb.m.Lock()
	defer qb.m.Unlock()
	delete(qb.idKeys, id)
	delete(qb.structKeys, *structKey)
	delete(qb.byAuthor, structKey.Author)
	deletionIndex := qb.indexRandList[id]
	substituteId := qb.randomList[len(qb.randomList)-1]
	qb.randomList[deletionIndex] = substituteId
	qb.indexRandList[substituteId] = deletionIndex
	delete(qb.indexRandList, id)
	qb.randomList = qb.randomList[:len(qb.randomList)-1]
	return nil
}

func (qb *quoteBook) findForDelete(id int64) (*Quote, error) {
	qb.m.RLock()
	defer qb.m.RUnlock()
	structKeyWithId, ok := qb.idKeys[id]
	if !ok {
		return nil, ErrQuoteNotFound
	}
	return &structKeyWithId.Quote, nil
}
