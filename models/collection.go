package models

import (
	// Native packages
	"math"
)

/**
 *	Collection struct, used to help with resource collection pagination.
 */
type Collection struct {
	Pointer   int `json:"page"`
	PageCount int `json:"pageCount"`
	Count     int `json:"records"`
	Limit     int `json:"recordsPerPage"`
}

/**
 *	Returns current collection pointer.
 *
 *	@return int
 */
func (c Collection) GetPointer() int {
	return c.Pointer
}

/**
 *	Sets current collection pointer, sets pointer to zero if newPointer is a negative value.
 *
 *	@param newPointer int - New collection pointer.
 *
 *	@return void
 */
func (c *Collection) SetPointer(newPointer int) {
	c.Pointer = newPointer

	if c.Pointer > 0 {
		c.Pointer = c.Pointer
	} else {
		c.Pointer = 0
	}
}

/**
 *	Sets collection count (number of resouces), also updates pages count.
 *
 *	@param count int - Current collection count.
 *
 *	@return void
 */
func (c *Collection) SetCount(count int) {
	c.Count = count
}

/**
 *	Returns the number of pages in current collection.
 *
 *	@return int
 */
func (c Collection) GetPageCount() int {
	numPages := int(math.Ceil(float64(c.Count) / float64(c.Limit)))
	return numPages
}

/**
 *	Sets collection page count.
 *
 *	@param pageCount int - Collection page count.
 *
 *	@return void
 */
func (c *Collection) SetPageCount(pageCount int) {
	c.PageCount = pageCount
}

/**
 *	Returns collection offset based on current pointer.
 *
 *	@return int
 */
func (c Collection) GetOffset() int {
	pointer := c.Pointer - 1

	if pointer < 0 {
		pointer = 0
	}

	offset := c.Limit * pointer

	return offset
}

/**
 *	Checks whether or not current collection pointer is out of bounds of collection range.
 *
 *	@return bool
 */
func (c *Collection) IsOutOfBounds() bool {
	numPages := c.GetPageCount()
	if c.Pointer > numPages || c.Pointer < numPages {
		return true
	}

	return false
}
