package dropbox // import "tantalic.com/dropbox"

import (
	"io"
	"time"
)

func deliver(res listFolderRes, ch chan MetaData) {
	for _, m := range res.Entries {
		ch <- m
	}

	if !res.HasMore {
		close(ch)
	}
}

func (c *Client) List(o ListOptions, ch chan MetaData) (string, error) {
	var res listFolderRes
	err := c.rpc("files/list_folder", o, &res)

	if err != nil {
		return "", err
	}

	go deliver(res, ch)

	if res.HasMore {
		return c.Continue(res.Cursor, ch)
	}

	return res.Cursor, nil
}

type ListOptions struct {
	Path                            string `json:"path"`
	Recursive                       bool   `json:"recursive,omitempty"`
	IncludeMediaInfo                bool   `json:"include_media_info,omitempty"`
	IncludeDeleted                  bool   `json:"include_deleted,omitempty"`
	IncludeHasExplicitSharedMembers bool   `json:"include_has_explicit_shared_members,omitempty"`
}

type listFolderRes struct {
	Entries []MetaData `json:"entries"`
	Cursor  string     `json:"cursor"`
	HasMore bool       `json:"has_more"`
}

func (c *Client) Continue(cursor string, ch chan MetaData) (string, error) {
	var res listFolderRes

	err := c.rpc("files/list_folder/continue", &continueReq{
		Cursor: cursor,
	}, &res)

	if err != nil {
		return cursor, err
	}

	go deliver(res, ch)

	if res.HasMore {
		return c.Continue(res.Cursor, ch)
	}

	return res.Cursor, nil
}

type continueReq struct {
	Cursor string `json:"cursor"`
}

func (c *Client) Watch(interval time.Duration, o ListOptions, ch chan MetaData) error {

	listChan := make(chan MetaData)
	go func() {
		for item := range listChan {
			ch <- item
		}
	}()

	cursor, err := c.List(o, listChan)
	if err != nil {
		return err
	}

	ticker := time.NewTicker(interval)
	for _ = range ticker.C {
		contChan := make(chan MetaData)
		go func() {
			for item := range contChan {
				ch <- item
			}
		}()

		cursor, _ = c.Continue(cursor, contChan)
	}

	return nil
}

func (c *Client) Download(path string) (io.ReadCloser, error) {
	return c.content("files/download", downloadReq{
		Path: path,
	})
}

type downloadReq struct {
	Path string `json:"path"`
}

type MetaData struct {
	Type        string `json:".tag"`
	Name        string `json:"name"`
	Size        uint64 `json:"size"`
	ID          string `json:"id"`
	Revision    string `json:"rev"`
	Path        string `json:"path_lower"`
	DisplayPath string `json:"path_display"`
}

func (m *MetaData) IsFolder() bool {
	return m.Type == "folder"
}

func (m *MetaData) IsFile() bool {
	return m.Type == "file"
}
