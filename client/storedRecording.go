package client

import (
	"github.com/CyCoreSystems/ari"
	"github.com/CyCoreSystems/ari-proxy/proxy"
)

type storedRecording struct {
	c *Client
}

func (sr *storedRecording) List() (ret []ari.StoredRecordingHandle, err error) {
	el, err := sr.c.listRequest(&proxy.Request{
		RecordingStoredList: &proxy.RecordingStoredList{},
	})
	if err != nil {
		return
	}
	for _, i := range el.List {
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
	data, err := sr.c.dataRequest(&proxy.Request{
		RecordingStoredData: &proxy.RecordingStoredData{
			ID: name,
		},
	})
	if err != nil {
		return
	}
	srd = data.StoredRecording
	return
}

func (sr *storedRecording) Copy(name string, dest string) (h ari.StoredRecordingHandle, err error) {
	err = sr.c.commandRequest(&proxy.Request{
		RecordingStoredCopy: &proxy.RecordingStoredCopy{
			ID: name,
		},
	})
	return
}

func (sr *storedRecording) Delete(name string) (err error) {
	err = sr.c.commandRequest(&proxy.Request{
		RecordingStoredDelete: &proxy.RecordingStoredDelete{
			ID: name,
		},
	})
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
