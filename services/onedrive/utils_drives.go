package onedrive

// the refactor of github.com/goh-chunlin/go-onedrive/onedrive upper layer functions since its internal implementation does not meet our needs.
// 1. refactor all URL to support path(instead of item id).
// 2. basic onedrive item struct, add we needed information like etag and so on.
// 3. refactor list to support pagination.
// 4. refactor write/read function with io.reader.
// 5. add new upload functions with `upload session`.
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/goh-chunlin/go-onedrive/onedrive"
)

const MB = 1024 * 1024
const OneTimeUploadedSizeLimit = int64(60 * MB)

type onedriveClient struct {
	client *onedrive.Client
}

// Item is the basic definition of onedrive item.
type Item struct {
	ID                   string                    `json:"id"`
	Name                 string                    `json:"name"`
	Etag                 string                    `json:"eTag"`
	LastModifiedDateTime time.Time                 `json:"lastModifiedDateTime"`
	DownloadURL          string                    `json:"@microsoft.graph.downloadUrl"`
	Size                 int64                     `json:"size"`
	CreatedBy            CreatedBy                 `json:"createdBy"`
	File                 *onedrive.DriveItemFile   `json:"file"`
	Folder               *onedrive.DriveItemFolder `json:"folder"`
}

type CreatedBy struct {
	User onedrive.User `json:"user"`
}

type ItemList struct {
	ODataContext string  `json:"@odata.context"`
	Count        int     `json:"@odata.count"`
	NextLink     string  `json:"@odata.nextLink"`
	Items        []*Item `json:"value"`
}

// GetItem gets item from onedrive.
func (o *onedriveClient) GetItem(ctx context.Context, absPath string) (*Item, error) {
	if !path.IsAbs(absPath) {
		return nil, errors.New(fmt.Sprintf("%s not an absolute path", absPath))
	}

	var apiURL string

	if isRoot(absPath) {
		apiURL = "me/drive/root"
	} else {
		apiURL = "me/drive/root:" + url.PathEscape(absPath)
	}

	req, err := o.client.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	var driveItem *Item
	err = o.client.Do(ctx, req, false, &driveItem)
	if err != nil {
		return nil, err
	}

	return driveItem, nil
}

// DownloadItem download onedrive item from download url.
func (o *onedriveClient) DownloadItem(ctx context.Context, absPath string) (rc io.ReadCloser, n int64, err error) {
	i, err := o.GetItem(ctx, absPath)
	if err != nil {
		return
	}

	var body io.Reader

	req, err := http.NewRequestWithContext(ctx, "GET", i.DownloadURL, body)
	if err != nil {
		return
	}

	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 && resp.StatusCode != 201 {
		errContent, _ := ioutil.ReadAll(resp.Body)
		err = fmt.Errorf("upload resource failed, cause %s:%s", resp.Status, string(errContent))
		return
	}

	return resp.Body, i.Size, nil
}

// DeleteItem delete an item from one drive.
func (o *onedriveClient) DeleteItem(ctx context.Context, absPath string) error {
	if !path.IsAbs(absPath) {
		return errors.New("not an absoulute path")
	}

	var apiURL string

	if isRoot(absPath) {
		apiURL = "me/drive/root"
	} else {
		apiURL = "me/drive/root:" + url.PathEscape(absPath)
	}

	req, err := o.client.NewRequest("DELETE", apiURL, nil)
	if err != nil {
		return err
	}

	err = o.client.Do(ctx, req, false, nil)
	if err != nil {
		return err
	}

	return err
}

type CreateUploadSessionRequest struct {
	ConflictBehavior string `json:"@microsoft.graph.conflictBehavior"`
	FileSize         int64  `json:"fileSize"`
	Description      string `json:"description"`
}

type CreateUploadSessionResponse struct {
	UploadURL string `json:"uploadUrl"`
}

// createUploadSession create an upload session
// ref: https://docs.microsoft.com/en-us/graph/api/driveitem-createuploadsession#create-an-upload-session
func (o *onedriveClient) createUploadSession(ctx context.Context, absPath string, size int64, description string) (uploadUrl string, err error) {
	reqBody := CreateUploadSessionRequest{
		ConflictBehavior: "replace",
		FileSize:         size,
		Description:      description,
	}

	apiURL := "me/drive/root:" + url.PathEscape(absPath) + ":/createUploadSession"

	req, err := o.client.NewRequest("POST", apiURL, reqBody)
	if err != nil {
		return
	}

	var response *CreateUploadSessionResponse
	err = o.client.Do(ctx, req, false, &response)
	if err != nil {
		return
	}

	uploadUrl = response.UploadURL

	return
}

// cancelUploadSession cancel an upload session.
// ref: https://docs.microsoft.com/en-us/graph/api/driveitem-createuploadsession#cancel-the-upload-session
func (o *onedriveClient) cancelUploadSession(ctx context.Context, uploadUrl string) error {
	req, err := http.NewRequest("DELETE", uploadUrl, nil)
	if err != nil {
		return err
	}

	err = o.client.Do(ctx, req, false, nil)

	return err
}

// Upload uploads file in multipart(each within 60MB).
// ref: https://docs.microsoft.com/en-us/graph/api/driveitem-createuploadsession#upload-bytes-to-the-upload-session
func (o *onedriveClient) Upload(ctx context.Context, absPath string, fileSize int64, r io.Reader, description string) (n int64, err error) {
	if !path.IsAbs(absPath) {
		return 0, errors.New("please provide the abspath")
	}

	// get onedrive upload url
	uploadUrl, err := o.createUploadSession(ctx, absPath, fileSize, description)
	if err != nil {
		return
	}

	defer func() {
		if err != nil {
			// cancel current upload session when error occurred
			_ = o.cancelUploadSession(ctx, uploadUrl)
		}
	}()

	// init upload size
	// max 60MB one time
	oneTimeUploadedSize := OneTimeUploadedSizeLimit
	if fileSize < oneTimeUploadedSize {
		oneTimeUploadedSize = fileSize
	}

	for n < fileSize {
		// upload segment
		// make minimum needed length, avoid memory waste
		segment := make([]byte, minInt64(oneTimeUploadedSize, fileSize-n))
		nowSize, err := r.Read(segment)
		if err != nil && err != io.EOF {
			return n, err
		}
		bn := bytes.NewReader(segment)

		req, err := http.NewRequestWithContext(ctx, "PUT", uploadUrl, bn)
		if err != nil {
			return n, err
		}

		req.Header.Add("Content-Length", strconv.Itoa(nowSize))
		req.Header.Add("Content-Range", fmt.Sprintf("bytes %d-%d/%d", n, n+int64(nowSize)-1, fileSize))

		// needn't authentication, so use plain http client
		// ref: https://docs.microsoft.com/en-us/graph/api/driveitem-createuploadsession#remarks
		client := &http.Client{}
		res, err := client.Do(req)
		if err != nil {
			return n, err
		}
		if res.StatusCode != 202 && res.StatusCode != 201 && res.StatusCode != 200 {
			errContent, _ := ioutil.ReadAll(res.Body)
			err = fmt.Errorf("upload resource failed, cause %s:%s", res.Status, string(errContent))
			return n, err
		}
		// close body
		res.Body.Close()

		n += int64(nowSize)
	}

	return
}

// List list items based of the absPath, with skipToken opts.
// ref: https://docs.microsoft.com/en-us/graph/api/driveitem-list-children
//      https://docs.microsoft.com/en-us/graph/query-parameters
func (o *onedriveClient) List(ctx context.Context, absPath string, skipToken string, limit uint32) (items []*Item, nextSkipToken string, err error) {
	if !path.IsAbs(absPath) {
		err = errors.New("please provide the abspath")
		return
	}

	var apiURL string

	if isRoot(absPath) {
		apiURL = fmt.Sprintf("me/drive/root/children?$top=%d&$skiptoken=%s", limit, skipToken)
	} else {
		apiURL = fmt.Sprintf("me/drive/root:%s:/children?$top=%d&$skiptoken=%s", url.PathEscape(absPath), limit, skipToken)
	}

	req, err := o.client.NewRequest("GET", apiURL, nil)
	if err != nil {
		return
	}

	var response ItemList
	err = o.client.Do(ctx, req, false, &response)
	if err != nil {
		return
	}

	if response.NextLink != "" {
		index := strings.Index(response.NextLink, "skiptoken")
		if index != -1 && index+9 < len(response.NextLink) {
			// parse skiptoken=xxx
			nextSkipToken = response.NextLink[index+9:]
		}
	}

	items = response.Items

	return
}

func isRoot(path string) bool {
	if filepath.Join("/", path) == "/" {
		return true
	}

	return false
}

// isFile judge file type(file)
func isFile(item *Item) bool {
	if item.File != nil {
		return true
	}

	return false
}

func minInt64(a, b int64) int64 {
	if a < b {
		return a
	}

	return b
}

// isNotFoundError check onedrive error whether is itemNotFound
func isNotFoundError(err error) bool {
	if err == nil {
		return false
	}

	return strings.Contains(err.Error(), "itemNotFound")
}
