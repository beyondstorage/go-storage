package s3

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	signerv4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	s3types "github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"

	"go.beyondstorage.io/credential"
	"go.beyondstorage.io/endpoint"
	ps "go.beyondstorage.io/v5/pairs"
	"go.beyondstorage.io/v5/pkg/httpclient"
	"go.beyondstorage.io/v5/services"
	typ "go.beyondstorage.io/v5/types"
)

// Service is the s3 service config.
type Service struct {
	cfg          aws.Config
	options      []func(*s3.Options)
	service      *s3.Client
	defaultPairs DefaultServicePairs
	features     ServiceFeatures

	typ.UnimplementedServicer
}

// String implements Servicer.String
func (s *Service) String() string {
	return fmt.Sprintf("Servicer s3")
}

// Storage is the s3 object storage service.
type Storage struct {
	service *s3.Client

	name    string
	workDir string

	defaultPairs DefaultStoragePairs
	features     StorageFeatures

	typ.UnimplementedStorager
	typ.UnimplementedDirer
	typ.UnimplementedMultiparter
	typ.UnimplementedLinker
	typ.UnimplementedStorageHTTPSigner
	typ.UnimplementedMultipartHTTPSigner
}

// String implements Storager.String
func (s *Storage) String() string {
	return fmt.Sprintf(
		"Storager s3 {Name: %s, WorkDir: %s}",
		s.name, s.workDir,
	)
}

// New will create both Servicer and Storager.
func New(pairs ...typ.Pair) (typ.Servicer, typ.Storager, error) {
	return newServicerAndStorager(pairs...)
}

// NewServicer will create Servicer only.
func NewServicer(pairs ...typ.Pair) (typ.Servicer, error) {
	return newServicer(pairs...)
}

// NewStorager will create Storager only.
func NewStorager(pairs ...typ.Pair) (typ.Storager, error) {
	_, store, err := newServicerAndStorager(pairs...)
	return store, err
}

func newServicer(pairs ...typ.Pair) (srv *Service, err error) {
	defer func() {
		if err != nil {
			err = services.InitError{Op: "new_servicer", Type: Type, Err: formatError(err), Pairs: pairs}
		}
	}()

	opt, err := parsePairServiceNew(pairs)
	if err != nil {
		return nil, err
	}

	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithHTTPClient(httpclient.New(opt.HTTPClientOptions)))
	if err != nil {
		return nil, err
	}

	// Set s3 config's http client
	if opt.HasHTTPClientOptions {
		cfg.HTTPClient = httpclient.New(opt.HTTPClientOptions)
	}

	var opts []func(*s3.Options)

	// Handle credential
	cp, err := credential.Parse(opt.Credential)
	if err != nil {
		return nil, err
	}
	switch cp.Protocol() {
	case credential.ProtocolHmac:
		ak, sk := cp.Hmac()
		cfg.Credentials = aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(ak, sk, ""))
	default:
		return nil, services.PairUnsupportedError{Pair: ps.WithCredential(opt.Credential)}
	}

	// Parse endpoint.
	if opt.HasEndpoint {
		ep, err := endpoint.Parse(opt.Endpoint)
		if err != nil {
			return nil, err
		}

		var url string
		switch ep.Protocol() {
		case endpoint.ProtocolHTTP:
			url, _, _ = ep.HTTP()
		case endpoint.ProtocolHTTPS:
			url, _, _ = ep.HTTPS()
		default:
			return nil, services.PairUnsupportedError{Pair: ps.WithEndpoint(opt.Endpoint)}
		}
		opts = append(opts, s3.WithEndpointResolver(s3.EndpointResolverFromURL(url)))
	}

	// Handle s3 API options.
	//
	// S3 will calculate payload's content-sha256 by default, we change this behavior for following reasons:
	// - To support uploading content without seek support: stdin, bytes.Reader
	// - To allow user decide when to calculate the hash, especially for big files
	//
	// We will ignore all errors returned by middleware.Stack handler, as we don't know whether this middleware exists.
	apiOptions := s3.WithAPIOptions(func(stack *middleware.Stack) error {
		// With removing PayloadSHA256 and adding UnsignedPayload, signer will set "X-Amz-Content-Sha256" to "UNSIGNED-PAYLOAD"
		_ = signerv4.RemoveComputePayloadSHA256Middleware(stack)
		_ = signerv4.AddUnsignedPayloadMiddleware(stack)
		_ = signerv4.RemoveContentSHA256HeaderMiddleware(stack)
		_ = signerv4.AddContentSHA256HeaderMiddleware(stack)
		return nil
	})
	opts = append(opts, apiOptions)

	if opt.ForcePathStyle {
		opts = append(opts, func(o *s3.Options) {
			o.UsePathStyle = true
		})
	}
	if opt.UseAccelerate {
		opts = append(opts, func(o *s3.Options) {
			o.UseAccelerate = true
		})
	}
	if opt.UseArnRegion {
		opts = append(opts, func(o *s3.Options) {
			o.UseARNRegion = true
		})
	}

	service := s3.NewFromConfig(cfg, opts...)

	srv = &Service{
		cfg:     cfg,
		options: opts,
		service: service,
	}

	if opt.HasDefaultServicePairs {
		srv.defaultPairs = opt.DefaultServicePairs
	}
	if opt.HasServiceFeatures {
		srv.features = opt.ServiceFeatures
	}
	return
}

// New will create a new s3 service.
func newServicerAndStorager(pairs ...typ.Pair) (srv *Service, store *Storage, err error) {
	srv, err = newServicer(pairs...)
	if err != nil {
		return
	}

	store, err = srv.newStorage(pairs...)
	if err != nil {
		err = services.InitError{Op: "new_storager", Type: Type, Err: formatError(err), Pairs: pairs}
		return
	}
	return
}

// All available storage classes are listed here.
const (
	StorageClassStandard           = string(s3types.ObjectStorageClassStandard)
	StorageClassReducedRedundancy  = string(s3types.ObjectStorageClassReducedRedundancy)
	StorageClassGlacier            = string(s3types.ObjectStorageClassGlacier)
	StorageClassStandardIa         = string(s3types.ObjectStorageClassStandardIa)
	StorageClassOnezoneIa          = string(s3types.ObjectStorageClassOnezoneIa)
	StorageClassIntelligentTiering = string(s3types.ObjectStorageClassIntelligentTiering)
	StorageClassDeepArchive        = string(s3types.ObjectStorageClassDeepArchive)
)

func formatError(err error) error {
	if _, ok := err.(services.InternalError); ok {
		return err
	}

	e := &smithy.GenericAPIError{}
	if errors.As(err, &e) {
		switch e.Code {
		// AWS SDK will use status code to generate awserr.Error, so "NotFound" should also be supported.
		case "NoSuchKey", "NotFound":
			return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
		case "AccessDenied":
			return fmt.Errorf("%w: %v", services.ErrPermissionDenied, err)
		}
	}

	noSuchKey := &s3types.NoSuchKey{}
	notFound := &s3types.NotFound{}
	if errors.As(err, &noSuchKey) || errors.As(err, &notFound) {
		return fmt.Errorf("%w: %v", services.ErrObjectNotExist, err)
	}

	return fmt.Errorf("%w: %v", services.ErrUnexpected, err)
}

// newStorage will create a new client.
func (s *Service) newStorage(pairs ...typ.Pair) (st *Storage, err error) {
	opt, err := parsePairStorageNew(pairs)
	if err != nil {
		return nil, err
	}

	// Copy config to prevent change the service config.
	cfg := s.cfg
	if opt.HasLocation {
		cfg.Region = opt.Location
	}

	st = &Storage{
		service: s3.NewFromConfig(cfg, s.options...),
		name:    opt.Name,
		workDir: "/",
	}

	if opt.HasDefaultStoragePairs {
		st.defaultPairs = opt.DefaultStoragePairs
	}
	if opt.HasStorageFeatures {
		st.features = opt.StorageFeatures
	}
	if opt.HasWorkDir {
		st.workDir = opt.WorkDir
	}
	return st, nil
}

func (s *Service) formatError(op string, err error, name string) error {
	if err == nil {
		return nil
	}

	return services.ServiceError{
		Op:       op,
		Err:      formatError(err),
		Servicer: s,
		Name:     name,
	}
}

// getAbsPath will calculate object storage's abs path
func (s *Storage) getAbsPath(path string) string {
	if strings.HasPrefix(path, "/") {
		return strings.TrimPrefix(path, "/")
	}
	prefix := strings.TrimPrefix(s.workDir, "/")
	return prefix + path
}

// getRelPath will get object storage's rel path.
func (s *Storage) getRelPath(path string) string {
	prefix := strings.TrimPrefix(s.workDir, "/")
	return strings.TrimPrefix(path, prefix)
}

func (s *Storage) formatError(op string, err error, path ...string) error {
	if err == nil {
		return nil
	}

	return services.StorageError{
		Op:       op,
		Err:      formatError(err),
		Storager: s,
		Path:     path,
	}
}

func (s *Storage) formatFileObject(v s3types.Object) (o *typ.Object, err error) {
	o = s.newObject(false)
	o.ID = *v.Key
	o.Path = s.getRelPath(*v.Key)
	// If you have enabled virtual link, you will not get the accurate object type.
	// If you want to get the exact object mode, please use `stat`
	o.Mode |= typ.ModeRead

	o.SetContentLength(v.Size)
	o.SetLastModified(aws.ToTime(v.LastModified))

	if v.ETag != nil {
		o.SetEtag(*v.ETag)
	}

	var sm ObjectSystemMetadata
	//v.StorageClass's type is s3types.ObjectStorageClass, which is equivalent to string
	sm.StorageClass = string(v.StorageClass)
	o.SetSystemMetadata(sm)

	return
}

func (s *Storage) newObject(done bool) *typ.Object {
	return typ.NewObject(s, done)
}

// All available server side algorithm are listed here.
const (
	ServerSideEncryptionAes256 = string(s3types.ServerSideEncryptionAes256)
	ServerSideEncryptionAwsKms = string(s3types.ServerSideEncryptionAwsKms)
)

func calculateEncryptionHeaders(algo string, key []byte) (algorithm, keyBase64, keyMD5Base64 *string, err error) {
	if len(key) != 32 {
		err = ErrServerSideEncryptionCustomerKeyInvalid
		return
	}
	kB64 := base64.StdEncoding.EncodeToString(key)
	kMD5 := md5.Sum(key)
	kMD5B64 := base64.StdEncoding.EncodeToString(kMD5[:])
	return &algo, &kB64, &kMD5B64, nil
}

// multipartXXX are multipart upload restriction in S3, see more details at:
// https://docs.aws.amazon.com/AmazonS3/latest/userguide/qfacts.html
const (
	// multipartNumberMaximum is the max part count supported.
	multipartNumberMaximum = 10000
	// multipartSizeMaximum is the maximum size for each part, 5GB.
	multipartSizeMaximum = 5 * 1024 * 1024 * 1024
	// multipartSizeMinimum is the minimum size for each part, 5MB.
	multipartSizeMinimum = 5 * 1024 * 1024
)

const (
	// writeSizeMaximum is the maximum size for each object with a single PUT operation, 5GB.
	// ref: https://docs.aws.amazon.com/AmazonS3/latest/userguide/upload-objects.html
	writeSizeMaximum = 5 * 1024 * 1024 * 1024
)

func (s *Storage) formatGetObjectInput(path string, opt pairStorageRead) (input *s3.GetObjectInput, err error) {
	rp := s.getAbsPath(path)

	input = &s3.GetObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}

	if opt.HasOffset && opt.HasSize {
		input.Range = aws.String(fmt.Sprintf("bytes=%d-%d", opt.Offset, opt.Offset+opt.Size-1))
	} else if opt.HasOffset && !opt.HasSize {
		input.Range = aws.String(fmt.Sprintf("bytes=%d-", opt.Offset))
	} else if !opt.HasOffset && opt.HasSize {
		input.Range = aws.String(fmt.Sprintf("bytes=0-%d", opt.Size-1))
	}

	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		input.SSECustomerAlgorithm, input.SSECustomerKey, input.SSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return nil, err
		}
	}

	return
}

func (s *Storage) formatPutObjectInput(path string, size int64, opt pairStorageWrite) (input *s3.PutObjectInput, err error) {
	rp := s.getAbsPath(path)

	input = &s3.PutObjectInput{
		Bucket:        aws.String(s.name),
		Key:           aws.String(rp),
		ContentLength: size,
	}

	if opt.HasContentMd5 {
		input.ContentMD5 = &opt.ContentMd5
	}
	if opt.HasStorageClass {
		input.StorageClass = s3types.StorageClass(opt.StorageClass)
	}
	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}
	if opt.HasServerSideEncryptionBucketKeyEnabled {
		input.BucketKeyEnabled = opt.ServerSideEncryptionBucketKeyEnabled
	}
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		input.SSECustomerAlgorithm, input.SSECustomerKey, input.SSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return nil, err
		}
	}
	if opt.HasServerSideEncryptionAwsKmsKeyID {
		input.SSEKMSKeyId = &opt.ServerSideEncryptionAwsKmsKeyID
	}
	if opt.HasServerSideEncryptionContext {
		encodedKMSEncryptionContext := base64.StdEncoding.EncodeToString([]byte(opt.ServerSideEncryptionContext))
		input.SSEKMSEncryptionContext = &encodedKMSEncryptionContext
	}
	if opt.HasServerSideEncryption {
		input.ServerSideEncryption = s3types.ServerSideEncryption(opt.ServerSideEncryption)
	}

	return
}

func (s *Storage) formatAbortMultipartUploadInput(path string, opt pairStorageDelete) (input *s3.AbortMultipartUploadInput) {
	rp := s.getAbsPath(path)

	input = &s3.AbortMultipartUploadInput{
		Bucket:   aws.String(s.name),
		Key:      aws.String(rp),
		UploadId: aws.String(opt.MultipartID),
	}

	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}

	return
}

func (s *Storage) formatDeleteObjectInput(path string, opt pairStorageDelete) (input *s3.DeleteObjectInput, err error) {
	rp := s.getAbsPath(path)

	if opt.HasObjectMode && opt.ObjectMode.IsDir() {
		if !s.features.VirtualDir {
			err = services.PairUnsupportedError{Pair: ps.WithObjectMode(opt.ObjectMode)}
			return nil, err
		}

		rp += "/"
	}

	input = &s3.DeleteObjectInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}

	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}

	return
}

func (s *Storage) formatCreateMultipartUploadInput(path string, opt pairStorageCreateMultipart) (input *s3.CreateMultipartUploadInput, err error) {
	rp := s.getAbsPath(path)

	input = &s3.CreateMultipartUploadInput{
		Bucket: aws.String(s.name),
		Key:    aws.String(rp),
	}

	if opt.HasServerSideEncryptionBucketKeyEnabled {
		input.BucketKeyEnabled = opt.ServerSideEncryptionBucketKeyEnabled
	}
	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		input.SSECustomerAlgorithm, input.SSECustomerKey, input.SSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return nil, err
		}
	}
	if opt.HasServerSideEncryptionAwsKmsKeyID {
		input.SSEKMSKeyId = &opt.ServerSideEncryptionAwsKmsKeyID
	}
	if opt.HasServerSideEncryptionContext {
		encodedKMSEncryptionContext := base64.StdEncoding.EncodeToString([]byte(opt.ServerSideEncryptionContext))
		input.SSEKMSEncryptionContext = &encodedKMSEncryptionContext
	}
	if opt.HasServerSideEncryption {
		input.ServerSideEncryption = s3types.ServerSideEncryption(opt.ServerSideEncryption)
	}

	return
}

func (s *Storage) formatCompleteMultipartUploadInput(o *typ.Object, parts []*typ.Part, opt pairStorageCompleteMultipart) (input *s3.CompleteMultipartUploadInput) {
	upload := &s3types.CompletedMultipartUpload{}
	for _, v := range parts {
		upload.Parts = append(upload.Parts, s3types.CompletedPart{
			ETag: aws.String(v.ETag),
			// For users the `PartNumber` is zero-based. But for S3, the effective `PartNumber` is [1, 10000].
			// Set PartNumber=v.Index+1 here to ensure pass in an effective `PartNumber` for `CompletedPart`.
			PartNumber: int32(v.Index + 1),
		})
	}

	input = &s3.CompleteMultipartUploadInput{
		Bucket:          aws.String(s.name),
		Key:             aws.String(o.ID),
		MultipartUpload: upload,
		UploadId:        aws.String(o.MustGetMultipartID()),
	}

	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}

	return
}

func (s *Storage) formatUploadPartInput(o *typ.Object, size int64, index int, opt pairStorageWriteMultipart) (input *s3.UploadPartInput, err error) {
	input = &s3.UploadPartInput{
		Bucket: &s.name,
		// For S3, the `PartNumber` is [1, 10000]. But for users, the `PartNumber` is zero-based.
		// Set PartNumber=index+1 here to ensure pass in an effective `PartNumber` for `UploadPart`.
		// ref: https://docs.aws.amazon.com/AmazonS3/latest/userguide/mpuoverview.html
		PartNumber:    int32(index + 1),
		Key:           aws.String(o.ID),
		UploadId:      aws.String(o.MustGetMultipartID()),
		ContentLength: size,
	}
	if opt.HasExceptedBucketOwner {
		input.ExpectedBucketOwner = &opt.ExceptedBucketOwner
	}
	if opt.HasServerSideEncryptionCustomerAlgorithm {
		input.SSECustomerAlgorithm, input.SSECustomerKey, input.SSECustomerKeyMD5, err = calculateEncryptionHeaders(opt.ServerSideEncryptionCustomerAlgorithm, opt.ServerSideEncryptionCustomerKey)
		if err != nil {
			return nil, err
		}
	}

	return
}
