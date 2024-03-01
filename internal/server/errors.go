package server

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/infura"
	"github.com/goverland-labs/goverland-helpers-ens-resolver/internal/server/forms"
)

func ResolveError(err error) error {
	if err == nil {
		return nil
	}

	switch {
	case errors.Is(err, forms.ErrMissedAddresses), errors.Is(err, forms.ErrMissedDomains), errors.Is(err, forms.ErrWrongDomain), errors.Is(err, forms.ErrWrongAddress):
		return status.Error(codes.InvalidArgument, err.Error())
	case errors.Is(err, infura.ErrUnknownAddress):
		return status.Error(codes.NotFound, err.Error())
	default:
		return status.Error(codes.Internal, "something went wrong")
	}
}
