package did

import (
	"context"
	"fmt"

	"github.com/TBD54566975/ssi-sdk/crypto"
	"github.com/TBD54566975/ssi-sdk/did"
	"github.com/TBD54566975/ssi-sdk/util"
	"github.com/mr-tron/base58"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/tbd54566975/ssi-service/pkg/service/keystore"
)

func NewWebHandler(s *Storage, ks *keystore.Service) (MethodHandler, error) {
	if s == nil {
		return nil, errors.New("storage cannot be empty")
	}
	if ks == nil {
		return nil, errors.New("keystore cannot be empty")
	}
	return &webHandler{method: did.WebMethod, storage: s, keyStore: ks}, nil
}

type webHandler struct {
	method   did.Method
	storage  *Storage
	keyStore *keystore.Service
}

type CreateWebDIDOptions struct {
	// e.g. did:web:example.com
	DIDWebID string `json:"didWebId" validate:"required"`
}

func (c CreateWebDIDOptions) Method() did.Method {
	return did.WebMethod
}

func (h *webHandler) GetMethod() did.Method {
	return h.method
}

func (h *webHandler) CreateDID(ctx context.Context, request CreateDIDRequest) (*CreateDIDResponse, error) {
	logrus.Debugf("creating DID: %+v", request)

	// process options
	if request.Options == nil {
		return nil, errors.New("options cannot be empty")
	}
	opts, ok := request.Options.(CreateWebDIDOptions)
	if !ok || request.Options.Method() != did.WebMethod {
		return nil, fmt.Errorf("invalid options for method, expected %s, got %s", did.WebMethod, request.Options.Method())
	}
	if err := util.IsValidStruct(opts); err != nil {
		return nil, errors.Wrap(err, "processing options")
	}

	didWeb := did.DIDWeb(opts.DIDWebID)

	if !didWeb.IsValid() {
		return nil, fmt.Errorf("could not resolve did:web DID: %s", didWeb)
	}

	pubKey, privKey, err := crypto.GenerateKeyByKeyType(request.KeyType)
	if err != nil {
		return nil, errors.Wrap(err, "could not generate key for did:web")
	}

	pubKeyBytes, err := crypto.PubKeyToBytes(pubKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not convert public key to byte")
	}

	doc, err := didWeb.CreateDoc(request.KeyType, pubKeyBytes)
	if err != nil {
		return nil, errors.Wrap(err, "could not create did:web doc")
	}

	// store metadata in DID storage
	id := didWeb.String()
	storedDID := DefaultStoredDID{
		ID:          id,
		DID:         *doc,
		SoftDeleted: false,
	}
	if err = h.storage.StoreDID(ctx, storedDID); err != nil {
		return nil, errors.Wrap(err, "could not store did:web value")
	}

	// convert to a serialized format for return to the client
	privKeyBytes, err := crypto.PrivKeyToBytes(privKey)
	if err != nil {
		return nil, errors.Wrap(err, "could not encode private key as base58")
	}
	privKeyBase58 := base58.Encode(privKeyBytes)

	// store private key in key storage
	keyStoreRequest := keystore.StoreKeyRequest{
		ID:               doc.VerificationMethod[0].ID,
		Type:             request.KeyType,
		Controller:       id,
		PrivateKeyBase58: privKeyBase58,
	}

	if err = h.keyStore.StoreKey(ctx, keyStoreRequest); err != nil {
		return nil, errors.Wrap(err, "could not store did:web private key")
	}

	return &CreateDIDResponse{
		DID:              storedDID.DID,
		PrivateKeyBase58: privKeyBase58,
		KeyType:          request.KeyType,
	}, nil
}

func (h *webHandler) GetDID(ctx context.Context, request GetDIDRequest) (*GetDIDResponse, error) {
	logrus.Debugf("getting DID: %+v", request)

	id := request.ID
	gotDID, err := h.storage.GetDIDDefault(ctx, id)
	if err != nil {
		return nil, errors.Wrapf(err, "error getting DID: %s", id)
	}
	if gotDID == nil {
		return nil, fmt.Errorf("did with id<%s> could not be found", id)
	}
	return &GetDIDResponse{DID: gotDID.GetDocument()}, nil
}

func (h *webHandler) GetDIDs(ctx context.Context) (*GetDIDsResponse, error) {
	logrus.Debug("getting did:web DID")

	gotDIDs, err := h.storage.GetDIDsDefault(ctx, did.WebMethod.String())
	if err != nil {
		return nil, errors.Wrap(err, "getting did:web DIDs")
	}
	dids := make([]did.Document, 0, len(gotDIDs))
	for _, gotDID := range gotDIDs {
		if !gotDID.IsSoftDeleted() {
			dids = append(dids, gotDID.GetDocument())
		}
	}
	return &GetDIDsResponse{DIDs: dids}, nil
}

func (h *webHandler) SoftDeleteDID(ctx context.Context, request DeleteDIDRequest) error {
	logrus.Debugf("soft deleting DID: %+v", request)

	id := request.ID
	gotStoredDID, err := h.storage.GetDIDDefault(ctx, id)
	if err != nil {
		return errors.Wrapf(err, "getting DID: %s", id)
	}
	if gotStoredDID == nil {
		return fmt.Errorf("did with id<%s> could not be found", id)
	}

	gotStoredDID.SoftDeleted = true

	return h.storage.StoreDID(ctx, *gotStoredDID)
}
