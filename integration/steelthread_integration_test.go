package integration

import (
	"testing"

	credsdk "github.com/TBD54566975/ssi-sdk/credential"
	"github.com/stretchr/testify/assert"

	"github.com/tbd54566975/ssi-service/pkg/service/operation/storage"
)

var steelThreadContext = NewTestContext("SteelThread")

func TestCreateIssuerDIDKeyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	didKeyOutput, err := CreateDIDKey()
	assert.NoError(t, err)

	issuerDID, err := getJSONElement(didKeyOutput, "$.did.id")
	assert.NoError(t, err)
	assert.Contains(t, issuerDID, "did:key")
	SetValue(steelThreadContext, "issuerDID", issuerDID)

	issuerKID, err := getJSONElement(didKeyOutput, "$.did.verificationMethod[0].id")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerKID)
	SetValue(steelThreadContext, "issuerKID", issuerKID)
}

func TestResolveIONDIDIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	resolveOutput, err := ResolveDID("did:ion:EiD3DIbDgBCajj2zCkE48x74FKTV9_Dcu1u_imzZddDKfg")
	assert.NoError(t, err)
	assert.NotEmpty(t, resolveOutput)

	ionDID, err := getJSONElement(resolveOutput, "$.didDocument.id")
	assert.NoError(t, err)
	assert.Equal(t, "did:ion:EiD3DIbDgBCajj2zCkE48x74FKTV9_Dcu1u_imzZddDKfg", ionDID)
}

func TestCreateAliceDIDKeyIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	didKeyOutput, err := CreateDIDKey()
	assert.NoError(t, err)
	assert.NotEmpty(t, didKeyOutput)

	aliceDID, err := getJSONElement(didKeyOutput, "$.did.id")
	assert.NoError(t, err)
	assert.Contains(t, aliceDID, "did:key")
	SetValue(steelThreadContext, "aliceDID", aliceDID)

	aliceKID, err := getJSONElement(didKeyOutput, "$.did.verificationMethod[0].id")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceKID)
	SetValue(steelThreadContext, "aliceKID", aliceKID)

	aliceDIDPrivateKey, err := getJSONElement(didKeyOutput, "$.privateKeyBase58")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceDID)
	SetValue(steelThreadContext, "aliceDIDPrivateKey", aliceDIDPrivateKey)
}

func TestCreateSchemaIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	output, err := CreateKYCSchema()
	assert.NoError(t, err)

	schemaID, err := getJSONElement(output, "$.id")
	assert.NoError(t, err)
	assert.NotEmpty(t, schemaID)
	SetValue(steelThreadContext, "schemaID", schemaID)
}

func TestCreateVerifiableCredentialIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	issuerDID, err := GetValue(steelThreadContext, "issuerDID")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerDID)

	issuerKID, err := GetValue(steelThreadContext, "issuerKID")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerKID)

	schemaID, err := GetValue(steelThreadContext, "schemaID")
	assert.NoError(t, err)
	assert.NotEmpty(t, schemaID)

	vcOutput, err := CreateVerifiableCredential(credInputParams{
		IssuerID:  issuerDID.(string),
		IssuerKID: issuerKID.(string),
		SchemaID:  schemaID.(string),
		SubjectID: issuerDID.(string),
	}, false)
	assert.NoError(t, err)
	assert.NotEmpty(t, vcOutput)

	credentialJWT, err := getJSONElement(vcOutput, "$.credentialJwt")
	assert.NoError(t, err)
	assert.NotEmpty(t, credentialJWT)
	SetValue(steelThreadContext, "credentialJWT", credentialJWT)
}

func TestCreateCredentialManifestIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	issuerDID, err := GetValue(steelThreadContext, "issuerDID")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerDID)

	issuerKID, err := GetValue(steelThreadContext, "issuerKID")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerKID)

	schemaID, err := GetValue(steelThreadContext, "schemaID")
	assert.NoError(t, err)
	assert.NotEmpty(t, schemaID)

	cmOutput, err := CreateCredentialManifest(credManifestParams{
		IssuerID:  issuerDID.(string),
		IssuerKID: issuerKID.(string),
		SchemaID:  schemaID.(string),
	})
	assert.NoError(t, err)

	presentationDefinitionID, err := getJSONElement(cmOutput, "$.credential_manifest.presentation_definition.id")
	assert.NoError(t, err)
	assert.NotEmpty(t, presentationDefinitionID)
	SetValue(steelThreadContext, "presentationDefinitionID", presentationDefinitionID)

	manifestID, err := getJSONElement(cmOutput, "$.credential_manifest.id")
	assert.NoError(t, err)
	assert.NotEmpty(t, manifestID)
	SetValue(steelThreadContext, "manifestID", manifestID)
}

func TestCreateIssuanceTemplateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	issuerDID, err := GetValue(steelThreadContext, "issuerDID")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerDID)

	issuerKID, err := GetValue(steelThreadContext, "issuerKID")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuerKID)

	schemaID, err := GetValue(steelThreadContext, "schemaID")
	assert.NoError(t, err)
	assert.NotEmpty(t, schemaID)

	cmOutput, err := CreateCredentialManifest(credManifestParams{
		IssuerID:  issuerDID.(string),
		IssuerKID: issuerKID.(string),
		SchemaID:  schemaID.(string),
	})
	assert.NoError(t, err)

	manifestID, err := getJSONElement(cmOutput, "$.credential_manifest.id")
	assert.NoError(t, err)
	assert.NotEmpty(t, manifestID)
	SetValue(steelThreadContext, "manifestWithIssuanceTemplateID", manifestID)

	presentationDefinitionID, err := getJSONElement(cmOutput, "$.credential_manifest.presentation_definition.id")
	assert.NoError(t, err)
	assert.NotEmpty(t, presentationDefinitionID)
	SetValue(steelThreadContext, "presentationDefinitionWithIssuanceTemplateID", presentationDefinitionID)

	itOutput, err := CreateIssuanceTemplate(issuanceTemplateParams{
		SchemaID:   schemaID.(string),
		ManifestID: manifestID,
		IssuerID:   issuerDID.(string),
		IssuerKID:  issuerKID.(string),
	})
	assert.NoError(t, err)

	issuanceTemplateID, err := getJSONElement(itOutput, "$.id")
	assert.NoError(t, err)
	assert.NotEmpty(t, issuanceTemplateID)
	SetValue(steelThreadContext, "issuanceTemplateID", issuanceTemplateID)
}

func TestSubmitApplicationWithIssuanceTemplateIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	credentialJWT, err := GetValue(steelThreadContext, "credentialJWT")
	assert.NoError(t, err)
	assert.NotEmpty(t, credentialJWT)

	presentationDefinitionID, err := GetValue(steelThreadContext, "presentationDefinitionWithIssuanceTemplateID")
	assert.NoError(t, err)
	assert.NotEmpty(t, presentationDefinitionID)

	manifestID, err := GetValue(steelThreadContext, "manifestWithIssuanceTemplateID")
	assert.NoError(t, err)
	assert.NotEmpty(t, manifestID)

	aliceDID, err := GetValue(steelThreadContext, "aliceDID")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceDID)

	aliceKID, err := GetValue(steelThreadContext, "aliceKID")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceKID)

	aliceDIDPrivateKey, err := GetValue(steelThreadContext, "aliceDIDPrivateKey")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceDIDPrivateKey)

	credAppJWT, err := CreateCredentialApplicationJWT(credApplicationParams{
		DefinitionID: presentationDefinitionID.(string),
		ManifestID:   manifestID.(string),
	}, credentialJWT.(string), aliceDID.(string), aliceKID.(string), aliceDIDPrivateKey.(string))
	assert.NoError(t, err)
	assert.NotEmpty(t, credAppJWT)

	submitApplicationOutput, err := SubmitApplication(applicationParams{
		ApplicationJWT: credAppJWT,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, submitApplicationOutput)

	isDone, err := getJSONElement(submitApplicationOutput, "$.done")
	assert.NoError(t, err)
	assert.Equal(t, "true", isDone)

	credentialResponseID, err := getJSONElement(submitApplicationOutput, "$.result.response.credential_response.id")
	assert.NoError(t, err)

	opCredentialResponse, err := getJSONElement(submitApplicationOutput, "$.result.response")
	assert.NoError(t, err)

	responsesOutput, err := get(endpoint + version + "manifests/responses/" + credentialResponseID)
	assert.NoError(t, err)

	assert.JSONEq(t, responsesOutput, opCredentialResponse)
}
func TestSubmitAndReviewApplicationIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	credentialJWT, err := GetValue(steelThreadContext, "credentialJWT")
	assert.NoError(t, err)
	assert.NotEmpty(t, credentialJWT)

	presentationDefinitionID, err := GetValue(steelThreadContext, "presentationDefinitionID")
	assert.NoError(t, err)
	assert.NotEmpty(t, presentationDefinitionID)

	manifestID, err := GetValue(steelThreadContext, "manifestID")
	assert.NoError(t, err)
	assert.NotEmpty(t, manifestID)

	aliceDID, err := GetValue(steelThreadContext, "aliceDID")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceDID)

	aliceKID, err := GetValue(steelThreadContext, "aliceKID")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceKID)

	aliceDIDPrivateKey, err := GetValue(steelThreadContext, "aliceDIDPrivateKey")
	assert.NoError(t, err)
	assert.NotEmpty(t, aliceDIDPrivateKey)

	credAppJWT, err := CreateCredentialApplicationJWT(credApplicationParams{
		DefinitionID: presentationDefinitionID.(string),
		ManifestID:   manifestID.(string),
	}, credentialJWT.(string), aliceDID.(string), aliceKID.(string), aliceDIDPrivateKey.(string))
	assert.NoError(t, err)
	assert.NotEmpty(t, credAppJWT)

	submitApplicationOutput, err := SubmitApplication(applicationParams{
		ApplicationJWT: credAppJWT,
	})
	assert.NoError(t, err)
	assert.NotEmpty(t, submitApplicationOutput)

	isDone, err := getJSONElement(submitApplicationOutput, "$.done")
	assert.NoError(t, err)
	assert.Equal(t, "false", isDone)
	opID, err := getJSONElement(submitApplicationOutput, "$.id")
	assert.NoError(t, err)

	reviewApplicationOutput, err := ReviewApplication(reviewApplicationParams{
		ID:       storage.StatusObjectID(opID),
		Approved: true,
		Reason:   "oh yeah im testing",
	})
	assert.NoError(t, err)

	crManifestID, err := getJSONElement(reviewApplicationOutput, "$.credential_response.manifest_id")
	assert.NoError(t, err)
	assert.Equal(t, manifestID, crManifestID)

	vc, err := getJSONElement(reviewApplicationOutput, "$.verifiableCredentials[0]")
	assert.NoError(t, err)
	assert.NotEmpty(t, vc)
	_, _, typedVC, err := credsdk.ToCredential(vc)
	assert.NoError(t, err)
	assert.Equal(t, "Mister", typedVC.CredentialSubject["givenName"])
	assert.Equal(t, "Tee", typedVC.CredentialSubject["familyName"])

	operationOutput, err := get(endpoint + version + "operations/" + opID)
	assert.NoError(t, err)
	isDone, err = getJSONElement(operationOutput, "$.done")
	assert.NoError(t, err)
	assert.Equal(t, "true", isDone)

	opCredentialResponse, err := getJSONElement(operationOutput, "$.result.response")
	assert.NoError(t, err)
	assert.JSONEq(t, reviewApplicationOutput, opCredentialResponse)
}
