package onfido

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"time"
)

// LivePhoto represents a LivePhoto in Onfido API
type LivePhoto struct {
	ID           string     `json:"id,omitempty"`
	CreatedAt    *time.Time `json:"created_at,omitempty"`
	Href         string     `json:"href,omitempty"`
	DownloadHref string     `json:"download_href,omitempty"`
	FileName     string     `json:"file_name,omitempty"`
	FileType     string     `json:"file_type,omitempty"`
	FileSize     int32      `json:"file_size,omitempty"`
}

// LivePhotoIter represents a LivePhoto iterator
type LivePhotoIter struct {
	*iter
}

// LivePhoto returns the current item in the iterator as a LivePhoto.
func (i *LivePhotoIter) LivePhoto() *LivePhoto {
	return i.Current().(*LivePhoto)
}

func (c *Client) UploadLivePhoto(ctx context.Context, applicantID string, file io.ReadSeeker) (*LivePhoto, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := createFormFile(writer, "file", file)
	if err != nil {
		return nil, err
	}
	if _, err := io.Copy(part, file); err != nil {
		return nil, err
	}
	if err := writer.WriteField("applicant_id", applicantID); err != nil {
		return nil, err
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}

	req, err := c.newRequest("POST", "/live_photos", body)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	if err != nil {
		return nil, err
	}

	var resp LivePhoto
	_, err = c.do(ctx, req, &resp)

	return &resp, err
}

func (c *Client) GetLivePhoto(ctx context.Context, id string) (*LivePhoto, error) {
	req, err := c.newRequest("GET", "/live_photos/"+id, nil)
	if err != nil {
		return nil, err
	}

	var resp LivePhoto
	_, err = c.do(ctx, req, &resp)
	return &resp, err
}

// ListPhotos retrieves the list of photos for the provided applicant.
// see https://documentation.onfido.com/?shell#live-photos
func (c *Client) ListLivePhotos(applicantID string) *LivePhotoIter {
	return &LivePhotoIter{&iter{
		c:       c,
		nextURL: "/live_photos?applicant_id=" + applicantID,
		handler: func(body []byte) ([]interface{}, error) {
			var r struct {
				LivePhotos []*LivePhoto `json:"live_photos"`
			}

			if err := json.Unmarshal(body, &r); err != nil {
				return nil, err
			}

			values := make([]interface{}, len(r.LivePhotos))
			for i, v := range r.LivePhotos {
				values[i] = v
			}
			return values, nil
		},
	}}
}
