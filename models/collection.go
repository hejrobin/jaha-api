package models

import (
	// Native packages
	"math"
)

type Collection struct {
	Pointer   int         `json:"page"`
	PageCount int         `json:"pageCount"`
	Limit     int         `json:"recordsPerPage"`
	Count     int         `json:"recordCount"`
	Records   interface{} `json:"records"`
}

/**
 *	Return collection pointer.
 *
 *	@return int
 */
func (collection *Collection) GetPointer() int {
	return collection.Pointer
}

/**
 *	Sets collection pointer.
 *
 *	@param newPointer int
 *
 *	@return void
 */
func (collection *Collection) SetPointer(newPointer int) {
	collection.Pointer = int(math.Abs(float64(newPointer)))
}

/**
 *	Returns collection record count.
 *
 *	@return int
 */
func (collection *Collection) GetCount() int {
	return collection.Count
}

/**
 *	Sets collection record count.
 *
 *	@param newCount int
 *
 *	@return void
 */
func (collection *Collection) SetCount(newCount int) {
	collection.Count = int(math.Abs(float64(newCount)))
}

/**
 *	Returns collection record limit.
 *
 *	@return int
 */
func (collection *Collection) GetLimit() int {
	return collection.Limit
}

/**
 *	Sets collection record limit.
 *
 *	@param newLimit int
 *
 *	@return void
 */
func (collection *Collection) SetLimit(newLimit int) {
	collection.Limit = int(math.Abs(float64(newLimit)))
}

/**
 *	Returns collection page count.
 *
 *	@return int
 */
func (collection *Collection) GetPageCount() int {
	return collection.PageCount
}

/**
 *	Sets collection page count.
 *
 *	@param newPageCount int
 *
 *	@return void
 */
func (collection *Collection) SetPageCount(newPageCount int) {
	collection.PageCount = int(math.Abs(float64(newPageCount)))
}

/**
 *	Returns collection records.
 *
 *	@return interface{}
 */
func (collection *Collection) GetRecords() interface{} {
	return collection.Records
}

/**
 *	Sets collection records.
 *
 *	@param collectionRecords interface{}
 *
 *	@return void
 */
func (collection *Collection) SetRecords(collectionRecords interface{}) {
	collection.Records = collectionRecords
}

/**
 *	Return collection offset.
 *
 *	@return int
 */
func (collection *Collection) GetOffset() int {
	var pointer int
	var offset int

	pointer = collection.Pointer - 1

	if pointer < 0 {
		pointer = 0
	}

	offset = collection.Limit * pointer

	return int(math.Abs(float64(offset)))
}

/**
 *	Validates whether or not pointer is out of collection bounts.
 *
 *	@return bool
 */
func (collection *Collection) IsOutOfBounds() bool {
	numPages := collection.GetPageCount()
	if collection.Pointer > numPages || collection.Pointer < numPages {
		return true
	}

	return false
}

/**
 *	Grabs new records and sets appropriate collection properties.
 *
 *	@param collectionRecords interface{}
 *	@param newPointer int
 *	@param newCount int
 *
 *	@return void
 */
func (collection *Collection) Grab(collectionRecords interface{}, newPointer int, newCount int) {
	collection.SetRecords(collectionRecords)
	collection.SetPointer(newPointer)
	collection.SetCount(newCount)

	newPageCount := int(math.Ceil(float64(collection.Count) / float64(collection.Limit)))
	collection.SetPageCount(newPageCount)
}
