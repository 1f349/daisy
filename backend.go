package daisy

import (
	"context"
	"github.com/emersion/go-vcard"
	"github.com/emersion/go-webdav/carddav"
)

type Backend struct {
}

func (b *Backend) AddressbookHomeSetPath(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) AddressBook(ctx context.Context) (*carddav.AddressBook, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) GetAddressObject(ctx context.Context, path string, req *carddav.AddressDataRequest) (*carddav.AddressObject, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) ListAddressObjects(ctx context.Context, req *carddav.AddressDataRequest) ([]carddav.AddressObject, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) QueryAddressObjects(ctx context.Context, query *carddav.AddressBookQuery) ([]carddav.AddressObject, error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) PutAddressObject(ctx context.Context, path string, card vcard.Card, opts *carddav.PutAddressObjectOptions) (loc string, err error) {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) DeleteAddressObject(ctx context.Context, path string) error {
	//TODO implement me
	panic("implement me")
}

func (b *Backend) CurrentUserPrincipal(ctx context.Context) (string, error) {
	//TODO implement me
	panic("implement me")
}
