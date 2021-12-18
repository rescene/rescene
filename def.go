package rescene

import "errors"

// ErrCRC crc doesn't match
var ErrCRC = errors.New("rescene : crc error")

// ErrBadBlock block not formatted properly
var ErrBadBlock = errors.New("rescene : block not properly formatted")

// ErrBadFile file not properly formatted
var ErrBadFile = errors.New("rescene : file not properly formatted")

// ErrBadData invalid data
var ErrBadData = errors.New("rescene : incorrect data")

// ErrNoData data missing
var ErrNoData = errors.New("rescene : no data")

// ErrDuplSFV file referenced twice in sfv
var ErrDuplSFV = errors.New("rescene : duplicate file in sfv")
