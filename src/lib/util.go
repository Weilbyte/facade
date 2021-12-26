package lib

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
	peparser "github.com/saferwall/pe"
)

func GetAndValidate(targetDll string) (*peparser.File, error) {
	pe, err := peparser.New(targetDll, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while opening file (%s): %v", targetDll, err))
	}

	err = pe.Parse()
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Error while parsing file (%s): %v (is this a valid PE file?)", targetDll, err))
	}

	if !pe.IsDLL() {
		return nil, errors.New(fmt.Sprintf("File (%s) is not a DLL", targetDll))
	}
	return pe, nil
}

func getUUID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}
