package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type storedRecording struct {
	c *Client
}

func (sr *storedRecording) List() (ret []ari.StoredRecordingHandle, err error) {

	req := proxy.Request{
		RecordingStoredList: &proxy.RecordingStoredList{},
	}
	var resp proxy.Response
	err = sr.c.nc.Request(proxy.GetSubject(sr.c.prefix, sr.c.appName, ""), &req, &resp, sr.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	for _, i := range resp.EntityList.List {
		ret = append(ret, sr.Get(i.ID))
	}
	return
}

func (sr *storedRecording) Get(name string) ari.StoredRecordingHandle {
	return &storedRecordingHandle{
		id: name,
		s:  sr,
	}
}

func (sr *storedRecording) Data(name string) (srd *ari.StoredRecordingData, err error) {
	req := proxy.Request{
		RecordingStoredData: &proxy.RecordingStoredData{
			ID: name,
		},
	}
	var resp proxy.DataResponse
	err = sr.c.nc.Request(proxy.GetSubject(sr.c.prefix, sr.c.appName, ""), &req, &resp, sr.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	srd = resp.StoredRecordingData
	return
}

func (sr *storedRecording) Copy(name string, dest string) (h ari.StoredRecordingHandle, err error) {
	req := proxy.Request{
		RecordingStoredCopy: &proxy.RecordingStoredCopy{
			ID: name,
		},
	}
	var resp proxy.Response
	err = sr.c.nc.Request(proxy.CommandSubject(sr.c.prefix, sr.c.appName, ""), &req, &resp, sr.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	h = sr.Get(resp.Entity.ID)
	return
}

func (sr *storedRecording) Delete(name string) (err error) {
	req := proxy.Request{
		RecordingStoredDelete: &proxy.RecordingStoredDelete{
			ID: name,
		},
	}
	var resp proxy.Response
	err = sr.c.nc.Request(proxy.CommandSubject(sr.c.prefix, sr.c.appName, ""), &req, &resp, sr.c.requestTimeout)
	if err != nil {
		return
	}
	if err = resp.Err(); err != nil {
		return
	}
	return
}

type storedRecordingHandle struct {
	id string
	s  *storedRecording
}

func (s *storedRecordingHandle) Copy(dest string) (ari.StoredRecordingHandle, error) {
	return s.s.Copy(s.id, dest)
}

func (s *storedRecordingHandle) Data() (*ari.StoredRecordingData, error) {
	return s.s.Data(s.id)
}

func (s *storedRecordingHandle) Delete() error {
	return s.s.Delete(s.id)
}

func (s *storedRecordingHandle) ID() string {
	return s.id
}
