package libhosty

import "errors"

//ErrNotAnAddressLine used when operating on a non-address line for operation
// related to address lines, such as comment/uncomment
var ErrNotAnAddressLine = errors.New("This line is not of type ADDRESS")

//ErrUncommentableLine used when try to comment a line that cannot be commented
var ErrUncommentableLine = errors.New("This line cannot be commented")

//ErrAlredyCommentedLine used when try to comment an alredy commented line
var ErrAlredyCommentedLine = errors.New("This line is alredy commented")

//ErrAlredyUncommentedLine used when try to uncomment an alredy uncommented line
var ErrAlredyUncommentedLine = errors.New("This line is alredy uncommented")

//ErrAddressNotFound used when provided address is not found
var ErrAddressNotFound = errors.New("Cannot find a line with given address")

//ErrHostnameNotFound used when provided hostname is not found
var ErrHostnameNotFound = errors.New("Cannot find a line with given hostname")

//ErrUnknown used when we don't know what's happened
var ErrUnknown = errors.New("Unknown error")
