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
	qBook = &quoteBook{
		idKeys:        make(map[int64]*QuoteWithID),
		structKeys:    make(map[Quote]struct{}),
		indexRandList: make(map[int64]int),
		byAuthor:      make(map[string][]int64),
	}
	ErrQuoteNotFound = errors.New("quote not found")
	ErrQuoteExists   = errors.New("quote of such author already exists")
)

func Create(quote *Quote) (int64, error) {
	err := findForCreate(quote)
	if err != nil {
		return 0, err
	}
	qBook.m.Lock()
	defer qBook.m.Unlock()
	qBook.IdCount += 1
	qBook.idKeys[qBook.IdCount] = &QuoteWithID{
		ID:    qBook.IdCount,
		Quote: *quote,
	}
	qBook.structKeys[*quote] = struct{}{}
	qBook.byAuthor[quote.Author] = append(qBook.byAuthor[quote.Author], qBook.IdCount)
	qBook.randomList = append(qBook.randomList, qBook.IdCount)
	qBook.indexRandList[qBook.IdCount] = len(qBook.randomList) - 1

	return qBook.IdCount, nil
}

func findForCreate(quote *Quote) error {
	qBook.m.RLock()
	defer qBook.m.RUnlock()
	_, ok := qBook.structKeys[*quote]
	if ok {
		return ErrQuoteExists
	}
	return nil
}

func GetRandom() (*QuoteWithID, error) {
	var err error
	var randomIndex *big.Int

	qBook.m.RLock()
	defer qBook.m.RUnlock()
	listLength := len(qBook.randomList)
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

	id := qBook.randomList[randomIndex.Int64()]
	return qBook.idKeys[id], nil
}

func List(author *string) []*QuoteWithID {
	qBook.m.RLock()
	defer qBook.m.RUnlock()
	quotes := []*QuoteWithID{}
	if author == nil {
		for _, quote := range qBook.idKeys {
			quotes = append(quotes, quote)
		}
	} else {
		ids := qBook.byAuthor[*author]
		for _, id := range ids {
			quotes = append(quotes, qBook.idKeys[id])
		}
	}

	return quotes
}

func Delete(id int64) error {
	structKey, err := findForDelete(id)
	if err != nil {
		return err
	}
	qBook.m.Lock()
	defer qBook.m.Unlock()
	delete(qBook.idKeys, id)
	delete(qBook.structKeys, *structKey)
	delete(qBook.byAuthor, structKey.Author)
	deletionIndex := qBook.indexRandList[id]
	substituteId := qBook.randomList[len(qBook.randomList)-1]
	qBook.randomList[deletionIndex] = substituteId
	qBook.indexRandList[substituteId] = deletionIndex
	delete(qBook.indexRandList, id)
	qBook.randomList = qBook.randomList[:len(qBook.randomList)-1]
	return nil
}

func findForDelete(id int64) (*Quote, error) {
	qBook.m.RLock()
	defer qBook.m.RUnlock()
	structKeyWithId, ok := qBook.idKeys[id]
	if !ok {
		return nil, ErrQuoteNotFound
	}
	return &structKeyWithId.Quote, nil
}
